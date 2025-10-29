package utils

import (
	"ApkAdmin/global"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// WithdrawConfig 提现配置结构
type WithdrawConfig struct {
	MinWithdraw         float64  `json:"minWithdraw"`         // 最低提现金额
	MaxWithdraw         float64  `json:"maxWithdraw"`         // 最高提现金额
	DailyWithdrawCount  int      `json:"dailyWithdrawCount"`  // 每日提现次数
	WithdrawMethods     []string `json:"withdrawMethods"`     // 提现方式
	WithdrawFee         float64  `json:"withdrawFee"`         // 手续费百分比
	SettlementCycle     string   `json:"settlementCycle"`     // 结算周期
	WithdrawProcessDays string   `json:"withdrawProcessDays"` // 处理时长
}

// 工具函数
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GenerateFlowNo(prefix string) string {
	now := time.Now()
	// 格式：COM20251022001234567890
	return fmt.Sprintf("%s%s%09d",
		prefix,
		now.Format("20060102"),
		now.UnixNano()%1000000000,
	)
}

// parseWithdrawConfig 解析提现配置
func ParseWithdrawConfig(config interface{}) (*WithdrawConfig, error) {
	// 将 interface{} 转换为 map
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, errors.New("配置格式错误")
	}

	withdrawConfig := &WithdrawConfig{
		MinWithdraw:        getFloat64(configMap, "minWithdraw", 10),
		MaxWithdraw:        getFloat64(configMap, "maxWithdraw", 5000),
		DailyWithdrawCount: getInt(configMap, "dailyWithdrawCount", 3),
		WithdrawFee:        getFloat64(configMap, "withdrawFee", 0),
	}

	// 解析提现方式
	if methods, ok := configMap["withdrawMethods"].([]interface{}); ok {
		withdrawConfig.WithdrawMethods = make([]string, 0, len(methods))
		for _, m := range methods {
			if method, ok := m.(string); ok {
				withdrawConfig.WithdrawMethods = append(withdrawConfig.WithdrawMethods, method)
			}
		}
	}

	return withdrawConfig, nil
}

// GenerateWithdrawNo generateWithdrawNo 生成提现单号
func GenerateWithdrawNo(userID uint) string {
	// 格式：WD + 日期 + 用户ID后4位 + 随机4位数
	now := time.Now()
	dateStr := now.Format("20060102")
	userSuffix := fmt.Sprintf("%04d", userID%10000)
	randomNum := now.UnixNano() % 10000
	return fmt.Sprintf("WD%s%s%04d", dateStr, userSuffix, randomNum)
}

// getFloat64 从 map 中获取 float64 值，提供默认值
func getFloat64(m map[string]interface{}, key string, defaultVal float64) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return defaultVal
}

// getInt 从 map 中获取 int 值，提供默认值
func getInt(m map[string]interface{}, key string, defaultVal int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		}
	}
	return defaultVal
}

// ParseAccounts 解析账号 - 支持多种数据类型
func ParseAccounts(data interface{}) ([]string, error) {
	var accounts []string

	switch v := data.(type) {
	case string:
		// JSON 字符串
		if err := json.Unmarshal([]byte(v), &accounts); err != nil {
			return nil, fmt.Errorf("解析JSON失败: %w", err)
		}

	case []interface{}:
		// 接口数组
		for _, item := range v {
			if str, ok := item.(string); ok {
				accounts = append(accounts, str)
			}
		}

	case []string:
		// 字符串数组
		accounts = v

	default:
		return nil, fmt.Errorf("不支持的数据类型: %T", v)
	}

	return accounts, nil
}

// 随机选择账号
func RandomSelectAccount(accounts []string) string {
	if len(accounts) == 0 {
		return ""
	}
	rand.NewSource(time.Now().UnixNano())
	return accounts[rand.Intn(len(accounts))]
}

// 过滤有效账号
func FilterValidAccounts(accounts []string) []string {
	var valid []string
	for _, account := range accounts {
		trimmed := strings.TrimSpace(account)
		if trimmed != "" {
			valid = append(valid, trimmed)
		}
	}
	return valid
}

// BuildPublicUrl  构建公开文件URL
func BuildPublicUrl(objectName string) string {
	return fmt.Sprintf("https://%s.%s/%s",
		global.GVA_CONFIG.AliyunOSS.BucketName,
		global.GVA_CONFIG.AliyunOSS.Endpoint,
		objectName,
	)
}
