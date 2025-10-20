package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type CommissionTierService struct{}

// GetCommissionTierList 获取阶梯等级列表
func (s *CommissionTierService) GetCommissionTierList(name string, status *int8, page, pageSize int) (list []project.CommissionTier, total int64, err error) {
	db := global.GVA_DB.Model(&project.CommissionTier{})

	// 搜索条件
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	// 统计总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询，按排序字段排序
	offset := (page - 1) * pageSize
	err = db.Order("sort desc").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// GetCommissionTierById 根据ID获取阶梯等级
func (s *CommissionTierService) GetCommissionTierById(id int) (tier project.CommissionTier, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&tier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tier, fmt.Errorf("等级不存在")
	}
	return tier, err
}

// GetAllEnabledTiers 获取所有启用的等级（按排序）
func (s *CommissionTierService) GetAllEnabledTiers() ([]project.CommissionTier, error) {
	var tiers []project.CommissionTier
	err := global.GVA_DB.Where("status = ?", 1).Order("sort ASC").Find(&tiers).Error
	return tiers, err
}

// GetTierBySubordinateCount 根据下级人数获取匹配的等级
func (s *CommissionTierService) GetTierBySubordinateCount(count int) (tier project.CommissionTier, err error) {
	// 查找最高的符合条件的等级
	err = global.GVA_DB.Where("status = ? AND min_subordinates <= ?", 1, count).
		Order("min_subordinates DESC").
		First(&tier).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到，返回最低等级（min_subordinates = 0）
		err = global.GVA_DB.Where("status = ?", 1).
			Order("min_subordinates ASC").
			First(&tier).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tier, fmt.Errorf("没有可用的等级配置")
		}
	}

	return tier, err
}

// CreateCommissionTier 创建阶梯等级
func (s *CommissionTierService) CreateCommissionTier(tier *project.CommissionTier) error {
	// 检查等级名称是否重复
	var count int64
	err := global.GVA_DB.Model(&project.CommissionTier{}).Where("name = ?", tier.Name).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("等级名称已存在")
	}

	// 检查最低下级人数是否重复
	err = global.GVA_DB.Model(&project.CommissionTier{}).Where("min_subordinates = ?", tier.MinSubordinates).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该下级人数条件已存在")
	}

	now := time.Now()
	tier.CreateTime = &now
	tier.UpdateTime = &now

	return global.GVA_DB.Create(tier).Error
}

// UpdateCommissionTier 更新阶梯等级
func (s *CommissionTierService) UpdateCommissionTier(tier *project.CommissionTier) error {
	// 检查等级是否存在
	var existingTier project.CommissionTier
	err := global.GVA_DB.Where("id = ?", tier.ID).First(&existingTier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("等级不存在")
	}

	// 检查等级名称是否与其他等级重复
	var count int64
	err = global.GVA_DB.Model(&project.CommissionTier{}).
		Where("name = ? AND id != ?", tier.Name, tier.ID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("等级名称已存在")
	}

	// 检查最低下级人数是否与其他等级重复
	err = global.GVA_DB.Model(&project.CommissionTier{}).
		Where("min_subordinates = ? AND id != ?", tier.MinSubordinates, tier.ID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该下级人数条件已存在")
	}

	now := time.Now()
	tier.UpdateTime = &now

	// 只更新允许修改的字段
	return global.GVA_DB.Model(&project.CommissionTier{}).Where("id = ?", tier.ID).Updates(map[string]interface{}{
		"name":             tier.Name,
		"min_subordinates": tier.MinSubordinates,
		"rate":             tier.Rate,
		"color":            tier.Color,
		"icon":             tier.Icon,
		"sort":             tier.Sort,
		"status":           tier.Status,
		"update_time":      tier.UpdateTime,
	}).Error
}

// DeleteCommissionTiers 删除阶梯等级（支持批量）
func (s *CommissionTierService) DeleteCommissionTiers(ids []int) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择要删除的等级")
	}
	// 检查是否有用户正在使用这些等级
	var count int64
	err := global.GVA_DB.Model(&project.TeamStatistics{}).
		Where("current_tier_id IN ?", ids).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("有用户正在使用这些等级，无法删除")
	}

	// 批量删除
	result := global.GVA_DB.Where("id IN ?", ids).Delete(&project.CommissionTier{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("等级不存在")
	}

	return nil
}

// UpdateCommissionTierStatus 更新等级状态
func (s *CommissionTierService) UpdateCommissionTierStatus(id int, status *int) error {
	// 检查等级是否存在
	var tier project.CommissionTier
	err := global.GVA_DB.Where("id = ?", id).First(&tier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("等级不存在")
	}

	// 如果是禁用操作，检查是否有用户正在使用该等级
	if *status == 0 {
		var count int64
		err = global.GVA_DB.Model(&project.TeamStatistics{}).
			Where("current_tier_id = ?", id).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("有用户正在使用该等级，无法禁用")
		}
	}

	now := time.Now()
	return global.GVA_DB.Model(&project.CommissionTier{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_time": &now,
	}).Error
}

// UpdateCommissionTierSort 批量更新等级排序
func (s *CommissionTierService) UpdateCommissionTierSort(sorts []map[string]interface{}) error {
	// 开启事务
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for _, item := range sorts {
			id, ok1 := item["id"].(float64)
			sort, ok2 := item["sort"].(float64)

			if !ok1 || !ok2 {
				return fmt.Errorf("参数格式错误")
			}

			err := tx.Model(&project.CommissionTier{}).Where("id = ?", int(id)).Updates(map[string]interface{}{
				"sort":        int(sort),
				"update_time": &now,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUserCurrentTier 获取用户当前应享受的等级
func (s *CommissionTierService) GetUserCurrentTier(userId int64) (tier project.CommissionTier, err error) {
	// 统计用户的直属下级人数
	var count int64
	err = global.GVA_DB.Model(&project.User{}).Where("inviter_id = ?", userId).Count(&count).Error
	if err != nil {
		return tier, err
	}

	// 根据下级人数获取匹配的等级
	return s.GetTierBySubordinateCount(int(count))
}

// ValidateTierConfig 验证等级配置的完整性
func (s *CommissionTierService) ValidateTierConfig() error {
	// 检查是否存在 min_subordinates = 0 的基础等级
	var count int64
	err := global.GVA_DB.Model(&project.CommissionTier{}).
		Where("status = ? AND min_subordinates = ?", 1, 0).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("缺少基础等级配置（min_subordinates = 0）")
	}

	return nil
}
