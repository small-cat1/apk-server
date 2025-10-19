package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"gorm.io/gorm"
)

// UserMembershipService 用户套餐服务
type UserMembershipService struct {
}

func (u UserMembershipService) UpdateUserMembership(data map[string]interface{}, conditions ...func(*gorm.DB) *gorm.DB) error {
	query := global.GVA_DB.Model(&project.User{})
	// 应用所有条件
	for _, condition := range conditions {
		query = condition(query)
	}
	return query.Updates(data).Error
}
