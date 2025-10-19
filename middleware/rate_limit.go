package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

// ==================== 预定义的限流策略 ====================

// IPRateLimit IP限流
func IPRateLimit(prefix string, seconds int, limit int) gin.HandlerFunc {
	return LimitConfig{
		GenerationKey: func(c *gin.Context) string {
			return prefix + ":" + c.ClientIP()
		},
		CheckOrMark: DefaultCheckOrMark,
		Expire:      seconds,
		Limit:       limit,
	}.LimitWithTime()
}

// UserRateLimit 用户限流（需要先登录）
func UserRateLimit(prefix string, seconds int, limit int) gin.HandlerFunc {
	return LimitConfig{
		GenerationKey: func(c *gin.Context) string {
			userID := c.GetString("userID") // 从JWT中获取
			if userID == "" {
				userID = c.ClientIP() // 未登录则用IP
			}
			return prefix + ":" + userID
		},
		CheckOrMark: DefaultCheckOrMark,
		Expire:      seconds,
		Limit:       limit,
	}.LimitWithTime()
}

// PhoneRateLimit 手机号限流
func PhoneRateLimit(prefix string, seconds int, limit int) gin.HandlerFunc {
	return LimitConfig{
		GenerationKey: func(c *gin.Context) string {
			// 先读取原始 Body 数据
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				return prefix + ":unknown"
			}

			// 重新设置 Body，让后续 handler 可以继续使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 解析 JSON 获取手机号
			var req struct {
				Phone string `json:"phone"`
			}
			if err := json.Unmarshal(bodyBytes, &req); err != nil {
				return prefix + ":unknown"
			}
			if req.Phone == "" {
				return prefix + ":unknown"
			}
			return prefix + ":" + req.Phone
		},
		CheckOrMark: DefaultCheckOrMark,
		Expire:      seconds,
		Limit:       limit,
	}.LimitWithTime()
}

// ==================== 组合限流 ====================

// MultiRateLimit 多重限流（所有条件都要满足）
func MultiRateLimit(limits ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, limit := range limits {
			limit(c)
			if c.IsAborted() {
				return
			}
		}
		c.Next()
	}
}

// ==================== 常用场景封装 ====================

// RegisterLimit 注册接口限流（多重保护）
func RegisterLimit() gin.HandlerFunc {
	return MultiRateLimit(
		// IP小时级限制
		IPRateLimit("register:ip:hour", 3600, 3),
		// IP天级限制
		IPRateLimit("register:ip:day", 86400, 10),
	)
}

// LoginLimit 登录接口限流
func LoginLimit() gin.HandlerFunc {
	return IPRateLimit("login:ip", 3600, 20)
}

// SmsLimit 短信验证码限流
func SmsLimit() gin.HandlerFunc {
	return MultiRateLimit(
		// IP限制：每小时5次
		IPRateLimit("sms:ip", 3600, 5),
		// 手机号限制：每5分钟1次
		PhoneRateLimit("sms:phone", 300, 1),
	)
}

// CaptchaLimit 验证码限流
func CaptchaLimit() gin.HandlerFunc {
	return IPRateLimit("captcha:ip", 60, 10)
}

// ==================== 动态限流（根据配置）====================

type RateLimitConfig struct {
	Register struct {
		IPHour int `yaml:"ip_hour"`
		IPDay  int `yaml:"ip_day"`
	} `yaml:"register"`

	Login struct {
		IPHour int `yaml:"ip_hour"`
	} `yaml:"login"`

	SMS struct {
		IPHour      int `yaml:"ip_hour"`
		PhoneMinute int `yaml:"phone_minute"`
	} `yaml:"sms"`
}

var RateLimitConf = RateLimitConfig{
	Register: struct {
		IPHour int `yaml:"ip_hour"`
		IPDay  int `yaml:"ip_day"`
	}{
		IPHour: 3,
		IPDay:  10,
	},
	Login: struct {
		IPHour int `yaml:"ip_hour"`
	}{
		IPHour: 20,
	},
	SMS: struct {
		IPHour      int `yaml:"ip_hour"`
		PhoneMinute int `yaml:"phone_minute"`
	}{
		IPHour:      5,
		PhoneMinute: 1,
	},
}
