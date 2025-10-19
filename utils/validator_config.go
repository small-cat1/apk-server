package utils

import (
	"fmt"
	"strings"
)

// ValidateConfig 自定义配置验证函数
func ValidateConfig(providerCode string, config map[string]interface{}) error {
	// 根据不同的支付服务商验证配置
	switch providerCode {
	case "wechat":
		return validateWechatConfig(config)
	case "alipay":
		return validateAlipayConfig(config)
	case "unionpay":
		return validateUnionpayConfig(config)
	default:
		return validateGenericConfig(config)
	}
}

// validateWechatConfig 微信支付配置验证
func validateWechatConfig(config map[string]interface{}) error {
	requiredFields := []string{"app_id", "mch_id", "key"}
	for _, field := range requiredFields {
		if value, exists := config[field]; !exists {
			return fmt.Errorf("缺少必需字段: %s", field)
		} else if str, ok := value.(string); !ok || strings.TrimSpace(str) == "" {
			return fmt.Errorf("字段 %s 不能为空", field)
		}
	}

	// 验证app_id格式
	if appId, ok := config["app_id"].(string); ok {
		if len(appId) != 18 {
			return fmt.Errorf("微信app_id格式错误，应为18位字符")
		}
	}

	// 验证mch_id格式
	if mchId, ok := config["mch_id"].(string); ok {
		if len(mchId) < 8 || len(mchId) > 10 {
			return fmt.Errorf("微信mch_id格式错误，应为8-10位数字")
		}
	}

	return nil
}

// validateAlipayConfig 支付宝配置验证
func validateAlipayConfig(config map[string]interface{}) error {
	requiredFields := []string{"app_id", "private_key", "public_key"}

	for _, field := range requiredFields {
		if value, exists := config[field]; !exists {
			return fmt.Errorf("缺少必需字段: %s", field)
		} else if str, ok := value.(string); !ok || strings.TrimSpace(str) == "" {
			return fmt.Errorf("字段 %s 不能为空", field)
		}
	}

	// 验证app_id格式
	if appId, ok := config["app_id"].(string); ok {
		if len(appId) != 16 && len(appId) != 20 {
			return fmt.Errorf("支付宝app_id格式错误")
		}
	}

	return nil
}

// validateUnionpayConfig 银联支付配置验证
func validateUnionpayConfig(config map[string]interface{}) error {
	requiredFields := []string{"mer_id", "private_key", "public_key"}

	for _, field := range requiredFields {
		if value, exists := config[field]; !exists {
			return fmt.Errorf("缺少必需字段: %s", field)
		} else if str, ok := value.(string); !ok || strings.TrimSpace(str) == "" {
			return fmt.Errorf("字段 %s 不能为空", field)
		}
	}

	return nil
}

// validateGenericConfig 通用配置验证
func validateGenericConfig(config map[string]interface{}) error {
	// 通用必需字段
	requiredFields := []string{"app_id", "app_secret"}

	for _, field := range requiredFields {
		if value, exists := config[field]; !exists {
			return fmt.Errorf("缺少必需字段: %s", field)
		} else if str, ok := value.(string); !ok || strings.TrimSpace(str) == "" {
			return fmt.Errorf("字段 %s 不能为空", field)
		}
	}

	// 验证callback_url格式（如果存在）
	if callbackUrl, exists := config["callback_url"]; exists {
		if str, ok := callbackUrl.(string); ok && str != "" {
			if !isValidURL(str) {
				return fmt.Errorf("callback_url 格式不正确")
			}
		}
	}

	return nil
}

// isValidURL 验证URL格式
func isValidURL(urlStr string) bool {
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}

// ValidateJSONSchema JSON Schema验证（需要引入第三方库）
func ValidateJSONSchema(schema string, data map[string]interface{}) error {
	// 这里可以使用github.com/xeipuuv/gojsonschema库进行JSON Schema验证
	// 示例实现
	if schema == "" {
		return nil
	}

	// TODO: 实现JSON Schema验证
	return nil
}
