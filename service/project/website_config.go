package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SystemConfigService struct{}

// GetConfig 获取指定scope的配置
func (s *SystemConfigService) GetConfig(scope string) (map[string]interface{}, error) {
	var configs []project.SystemConfig

	// 查询指定scope的所有配置
	db := global.GVA_DB.Where("scope = ?", scope)

	// 特殊处理：website scope 排除 ios_account
	if scope == "website" {
		db = db.Where("`key` != ?", "ios_account")
	}

	err := db.Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// 如果没有配置，返回空map
	if len(configs) == 0 {
		return make(map[string]interface{}), nil
	}

	// 将配置转换为map
	result := make(map[string]interface{})
	for _, config := range configs {
		// 尝试解析JSON值
		var value interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			// 如果解析失败，直接使用字符串值
			result[config.Key] = config.Value
		} else {
			result[config.Key] = value
		}
	}

	return result, nil
}

// SetConfig 设置指定scope的配置
func (s *SystemConfigService) SetConfig(scope string, config map[string]interface{}) error {
	// 开启事务
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		for key, value := range config {
			// 将值序列化为JSON
			valueJSON, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
			}

			// 查询是否存在该配置
			var existingConfig project.SystemConfig
			err = tx.Where("scope = ? AND `key` = ?", scope, key).First(&existingConfig).Error

			if err == nil {
				// 配置存在，更新
				existingConfig.Value = string(valueJSON)
				existingConfig.UpdatedAt = &now
				if err := tx.Save(&existingConfig).Error; err != nil {
					return fmt.Errorf("failed to update config for key %s: %w", key, err)
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				// 配置不存在，创建
				newConfig := project.SystemConfig{
					Scope:     scope,
					Name:      s.generateConfigName(scope, key),
					Key:       key,
					Value:     string(valueJSON),
					CreatedAt: &now,
					UpdatedAt: &now,
				}
				if err := tx.Create(&newConfig).Error; err != nil {
					return fmt.Errorf("failed to create config for key %s: %w", key, err)
				}
			} else {
				// 其他错误
				return fmt.Errorf("failed to query config for key %s: %w", key, err)
			}
		}

		return nil
	})
}

// GetConfigByKey 获取指定scope和key的配置值
func (s *SystemConfigService) GetConfigByKey(scope, key string) (interface{}, error) {
	var config project.SystemConfig
	if err := global.GVA_DB.Where("scope = ? AND `key` = ?", scope, key).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("config not found for scope: %s, key: %s", scope, key)
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// 尝试解析JSON值
	var value interface{}
	if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
		// 如果解析失败，直接返回字符串值
		return config.Value, nil
	}

	return value, nil
}

// SetConfigByKey 设置指定scope和key的配置值
func (s *SystemConfigService) SetConfigByKey(scope, key string, value interface{}) error {
	// 将值序列化为JSON
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	now := time.Now()

	// 查询是否存在该配置
	var existingConfig project.SystemConfig
	err = global.GVA_DB.Where("scope = ? AND `key` = ?", scope, key).First(&existingConfig).Error

	if err == nil {
		// 配置存在，更新
		existingConfig.Value = string(valueJSON)
		existingConfig.UpdatedAt = &now
		return global.GVA_DB.Save(&existingConfig).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 配置不存在，创建
		newConfig := project.SystemConfig{
			Scope:     scope,
			Name:      s.generateConfigName(scope, key),
			Key:       key,
			Value:     string(valueJSON),
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		return global.GVA_DB.Create(&newConfig).Error
	}

	return fmt.Errorf("failed to query config: %w", err)
}

// DeleteConfigByKey 删除指定scope和key的配置
func (s *SystemConfigService) DeleteConfigByKey(scope, key string) error {
	result := global.GVA_DB.Where("scope = ? AND `key` = ?", scope, key).Delete(&project.SystemConfig{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete config: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("config not found for scope: %s, key: %s", scope, key)
	}
	return nil
}

// DeleteConfigByScope 删除指定scope的所有配置
func (s *SystemConfigService) DeleteConfigByScope(scope string) error {
	result := global.GVA_DB.Where("scope = ?", scope).Delete(&project.SystemConfig{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete configs: %w", result.Error)
	}
	return nil
}

// GetAllScopes 获取所有scope列表
func (s *SystemConfigService) GetAllScopes() ([]string, error) {
	var scopes []string

	if err := global.GVA_DB.Model(&project.SystemConfig{}).
		Distinct("scope").
		Pluck("scope", &scopes).Error; err != nil {
		return nil, fmt.Errorf("failed to get scopes: %w", err)
	}

	return scopes, nil
}

// BatchSetConfig 批量设置多个scope的配置
func (s *SystemConfigService) BatchSetConfig(configs map[string]map[string]interface{}) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for scope, config := range configs {
			if err := s.SetConfig(scope, config); err != nil {
				return fmt.Errorf("failed to set config for scope %s: %w", scope, err)
			}
		}
		return nil
	})
}

// generateConfigName 生成配置名称
func (s *SystemConfigService) generateConfigName(scope, key string) string {
	// 根据key生成友好的名称
	nameMap := map[string]map[string]string{
		"website": {
			"website_name":        "站点名称",
			"website_domain":      "站点域名",
			"website_logo":        "站点Logo",
			"website_description": "站点描述",
			"website_keywords":    "站点关键词",
			"website_switch":      "网站开关",
			"website_close_tip":   "关闭提示",
			"website_icp":         "ICP备案号",
			"copyright":           "版权信息",
			"ios_account":         "IOS账号配置",
		},
		"service": {
			"enabled":       "客服开关",
			"qq":            "QQ客服",
			"wechat":        "微信客服",
			"wechat_qrcode": "微信二维码",
			"phone":         "电话客服",
			"email":         "邮箱客服",
			"im_switch":     "第三方IM开关",
			"im_type":       "IM类型",
			"im_link":       "在线客服链接",
			"work_time":     "工作时间",
			"notice":        "客服通知",
		},
		"commission": {
			"minWithdraw":         "最低提现金额",
			"maxWithdraw":         "单次最高提现金额",
			"dailyWithdrawCount":  "每日提现次数",
			"settlementCycle":     "结算周期",
			"withdrawFee":         "提现手续费",
			"withdrawProcessDays": "提现到账时间",
		},
		"seo": {
			"seo_title":        "SEO标题",
			"seo_description":  "SEO描述",
			"seo_keywords":     "SEO关键词",
			"baidu_analytics":  "百度统计代码",
			"google_analytics": "Google Analytics ID",
		},
	}

	if scopeMap, ok := nameMap[scope]; ok {
		if name, ok := scopeMap[key]; ok {
			return name
		}
	}

	return key
}

// ValidateConfig 验证配置的完整性
func (s *SystemConfigService) ValidateConfig(scope string, config map[string]interface{}) error {
	// 根据scope验证必填项
	requiredFields := map[string][]string{
		"website":    {"website_name"},
		"service":    {},
		"commission": {"minWithdraw", "maxWithdraw", "dailyWithdrawCount", "settlementCycle", "withdrawProcessDays"},
		"seo":        {},
	}

	if fields, ok := requiredFields[scope]; ok {
		for _, field := range fields {
			if _, exists := config[field]; !exists {
				return fmt.Errorf("required field %s is missing", field)
			}
		}
	}

	return nil
}

// RefreshCache 刷新配置缓存（如果使用了缓存）
func (s *SystemConfigService) RefreshCache(scope string) error {
	// 如果使用了Redis等缓存，在这里实现缓存刷新逻辑
	// 例如：
	// cacheKey := fmt.Sprintf("system_config:%s", scope)
	// return global.Redis.Del(cacheKey).Err()
	return nil
}
