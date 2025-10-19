package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// SetSystemConfigRequest 设置系统配置请求
type SetSystemConfigRequest struct {
	Scope  string                 `json:"scope" binding:"required"`  // 作用域：website、service、pay、seo
	Config map[string]interface{} `json:"config" binding:"required"` // 配置项
}

// Validate 验证请求参数
func (r *SetSystemConfigRequest) Validate() error {
	// 验证scope
	validScopes := []string{"website", "service", "pay", "seo"}
	if !contains(validScopes, r.Scope) {
		return fmt.Errorf("invalid scope: %s, must be one of: %s", r.Scope, strings.Join(validScopes, ", "))
	}
	// 验证config不能为空
	if r.Config == nil || len(r.Config) == 0 {
		return errors.New("config cannot be empty")
	}
	// 根据不同的scope进行特定验证
	switch r.Scope {
	case "website":
		return r.validateWebsiteConfig()
	case "service":
		return r.validateServiceConfig()
	case "pay":
		return r.validatePayConfig()
	case "seo":
		return r.validateSEOConfig()
	default:
		return nil
	}
}

// validateWebsiteConfig 验证站点配置
func (r *SetSystemConfigRequest) validateWebsiteConfig() error {
	// 验证站点名称
	if name, ok := r.Config["website_name"].(string); ok && name != "" {
		if len(name) < 2 || len(name) > 50 {
			return errors.New("website_name length must be between 2 and 50")
		}
	}
	// 验证站点域名
	if domain, ok := r.Config["website_domain"].(string); ok && domain != "" {
		if !isValidURL(domain) {
			return errors.New("website_domain must be a valid URL")
		}
	}
	// 验证Logo URL
	if logo, ok := r.Config["website_logo"].(string); ok && logo != "" {
		if !isValidURL(logo) {
			return errors.New("website_logo must be a valid URL")
		}
	}

	// 验证ICP备案号格式（可选）
	if icp, ok := r.Config["website_icp"].(string); ok && icp != "" {
		if len(icp) > 100 {
			return errors.New("website_icp is too long")
		}
	}

	return nil
}

// validateServiceConfig 验证客服配置
func (r *SetSystemConfigRequest) validateServiceConfig() error {
	// 验证QQ号
	if qq, ok := r.Config["qq"].(string); ok && qq != "" {
		if !isValidQQ(qq) {
			return errors.New("invalid qq number format")
		}
	}

	// 验证邮箱
	if email, ok := r.Config["email"].(string); ok && email != "" {
		if !isValidEmail(email) {
			return errors.New("invalid email format")
		}
	}

	// 验证电话
	if phone, ok := r.Config["phone"].(string); ok && phone != "" {
		if !isValidPhone(phone) {
			return errors.New("invalid phone number format")
		}
	}

	// 验证微信二维码URL
	if qrcode, ok := r.Config["wechat_qrcode"].(string); ok && qrcode != "" {
		if !isValidURL(qrcode) {
			return errors.New("wechat_qrcode must be a valid URL")
		}
	}

	// 验证IM链接
	if imLink, ok := r.Config["im_link"].(string); ok && imLink != "" {
		if !isValidURL(imLink) {
			return errors.New("im_link must be a valid URL")
		}
	}

	return nil
}

// validatePayConfig 验证支付配置
func (r *SetSystemConfigRequest) validatePayConfig() error {
	// 验证回调地址
	if callbackURL, ok := r.Config["pay_callback_url"].(string); ok && callbackURL != "" {
		if !isValidURL(callbackURL) {
			return errors.New("pay_callback_url must be a valid URL")
		}
	}

	return nil
}

// validateSEOConfig 验证SEO配置
func (r *SetSystemConfigRequest) validateSEOConfig() error {
	// 验证Google Analytics ID格式
	if gaID, ok := r.Config["google_analytics"].(string); ok && gaID != "" {
		if !isValidGoogleAnalyticsID(gaID) {
			return errors.New("invalid google_analytics ID format")
		}
	}
	return nil
}

// GetSystemConfigRequest 获取系统配置请求
type GetSystemConfigRequest struct {
	Scope string `form:"scope" uri:"scope" binding:"required"` // 作用域
}

// Validate 验证请求参数
func (r *GetSystemConfigRequest) Validate() error {
	validScopes := []string{"website", "service", "pay", "seo"}
	if !contains(validScopes, r.Scope) {
		return fmt.Errorf("invalid scope: %s, must be one of: %s", r.Scope, strings.Join(validScopes, ", "))
	}
	return nil
}

type GetConfigByKeyRequest struct {
	Scope string `form:"scope" uri:"scope" binding:"required"` // 作用域
	Key   string `form:"key" uri:"key" binding:"required"`     // 作用域
}

// Validate 验证请求参数
func (r *GetConfigByKeyRequest) Validate() error {
	validScopes := []string{"website", "service", "pay", "seo"}
	if !contains(validScopes, r.Scope) {
		return fmt.Errorf("invalid scope: %s, must be one of: %s", r.Scope, strings.Join(validScopes, ", "))
	}
	return nil
}

// ============ 辅助验证函数 ============

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidURL 验证URL格式
func isValidURL(str string) bool {
	return true
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidPhone 验证手机号格式（支持中国大陆、香港、台湾等）
func isValidPhone(phone string) bool {
	if phone == "" {
		return false
	}

	// 保存原始号码
	original := phone

	// 移除空格和横线进行标准化
	normalized := strings.ReplaceAll(phone, " ", "")
	normalized = strings.ReplaceAll(normalized, "-", "")

	// 定义所有支持的电话号码格式
	patterns := []struct {
		name        string
		pattern     string
		useOriginal bool // 是否使用原始格式（带横线）
	}{
		{"中国大陆手机号", `^1[3-9]\d{9}$`, false},
		{"固定电话（区号+号码）", `^0\d{2,3}\d{7,8}$`, false},
		{"国际格式", `^\+\d{1,3}\d{7,14}$`, false},
		{"400电话", `^400\d{7,8}$`, false},
		{"800免费电话", `^800\d{7,8}$`, false},
		{"95客服热线", `^95\d{3,4}$`, false},
		{"运营商短号", `^1\d{4}$`, false},
		{"400/800带横线格式", `^(400|800)-?\d{3}-?\d{4,5}$`, true},
	}
	for _, p := range patterns {
		testPhone := normalized
		if p.useOriginal {
			testPhone = original
		}

		matched, err := regexp.MatchString(p.pattern, testPhone)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// isValidQQ 验证QQ号格式（5-11位数字）
func isValidQQ(qq string) bool {
	pattern := `^[1-9]\d{4,10}$`
	matched, _ := regexp.MatchString(pattern, qq)
	return matched
}

// isValidGoogleAnalyticsID 验证Google Analytics ID格式
func isValidGoogleAnalyticsID(id string) bool {
	// 支持 UA-XXXXXX-X 和 G-XXXXXXXXXX 格式
	pattern := `^(UA-\d{4,10}-\d{1,4}|G-[A-Z0-9]{10,})$`
	matched, _ := regexp.MatchString(pattern, id)
	return matched
}

// GetConfigValue 获取配置值（类型安全）
func (r *SetSystemConfigRequest) GetConfigValue(key string) (interface{}, bool) {
	val, ok := r.Config[key]
	return val, ok
}

// GetStringValue 获取字符串类型配置值
func (r *SetSystemConfigRequest) GetStringValue(key string) string {
	if val, ok := r.Config[key].(string); ok {
		return val
	}
	return ""
}

// GetBoolValue 获取布尔类型配置值
func (r *SetSystemConfigRequest) GetBoolValue(key string) bool {
	if val, ok := r.Config[key].(bool); ok {
		return val
	}
	return false
}

// GetIntValue 获取整数类型配置值
func (r *SetSystemConfigRequest) GetIntValue(key string) int {
	switch val := r.Config[key].(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		// 尝试解析字符串
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	default:
		return 0
	}
}

// ToJSON 转换为JSON字符串
func (r *SetSystemConfigRequest) ToJSON() (string, error) {
	data, err := json.Marshal(r.Config)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
