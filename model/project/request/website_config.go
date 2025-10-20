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
	Scope      string                 `json:"scope" binding:"required"`  // 作用域：website、service、pay、seo
	Config     map[string]interface{} `json:"config" binding:"required"` // 配置项
	GoogleCode string                 `json:"googleCode"`
}

// Validate 验证请求参数
func (r *SetSystemConfigRequest) Validate() error {
	// 验证scope
	validScopes := []string{"website", "service", "commission", "seo"}
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
	case "commission":
		return r.validateCommissionConfig()
	case "seo":
		return r.validateSEOConfig()
	default:
		return nil
	}
}

// validateCommissionConfig 验证分佣配置
func (r *SetSystemConfigRequest) validateCommissionConfig() error {

	// 验证最低提现金额
	if minWithdraw := r.GetFloat64Value("minWithdraw"); minWithdraw > 0 {
		if minWithdraw < 1 || minWithdraw > 10000 {
			return errors.New("minWithdraw must be between 1 and 10000")
		}
	}

	// 验证最高提现金额
	if maxWithdraw := r.GetFloat64Value("maxWithdraw"); maxWithdraw > 0 {
		if maxWithdraw < 100 || maxWithdraw > 100000 {
			return errors.New("maxWithdraw must be between 100 and 100000")
		}
	}

	// 验证最低提现不能大于等于最高提现
	minWithdraw := r.GetFloat64Value("minWithdraw")
	maxWithdraw := r.GetFloat64Value("maxWithdraw")
	if minWithdraw > 0 && maxWithdraw > 0 && minWithdraw >= maxWithdraw {
		return errors.New("minWithdraw must be less than maxWithdraw")
	}

	// 验证每日提现次数
	if dailyCount := r.GetIntValue("dailyWithdrawCount"); dailyCount > 0 {
		if dailyCount < 1 || dailyCount > 10 {
			return errors.New("dailyWithdrawCount must be between 1 and 10")
		}
	}

	// 验证提现手续费
	if fee := r.GetFloat64Value("withdrawFee"); fee > 0 {
		if fee < 0 || fee > 10 {
			return errors.New("withdrawFee must be between 0 and 10")
		}
	}

	// 验证必填字段
	requiredFields := []string{"minWithdraw", "maxWithdraw", "dailyWithdrawCount", "settlementCycle", "withdrawProcessDays"}
	for _, field := range requiredFields {
		if _, ok := r.Config[field]; !ok {
			return fmt.Errorf("required field %s is missing", field)
		}
	}

	return nil
}

// GetFloat64Value 获取浮点数类型配置值
func (r *SetSystemConfigRequest) GetFloat64Value(key string) float64 {
	switch val := r.Config[key].(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
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
	validScopes := []string{"website", "service", "commission", "seo"}
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
	validScopes := []string{"website", "service", "commission", "seo"}
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
