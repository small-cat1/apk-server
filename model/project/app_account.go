package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/system"
	"ApkAdmin/utils/crypto"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// AppAccount 应用账号表
type AppAccount struct {
	ID            uint                       `json:"id" gorm:"primarykey;comment:主键ID"`
	AppID         string                     `json:"app_id" gorm:"type:varchar(100);not null;index:idx_app_id;comment:应用唯一标识符" binding:"required"`
	AccountDetail string                     `json:"account_detail" gorm:"type:text;not null;comment:登录账号详情（加密存储）" binding:"required"`
	CategoryID    uint                       `json:"category_id" gorm:"not null;index:idx_category_id;comment:分类ID" binding:"required"`
	AccountNo     string                     `json:"account_no" gorm:"type:varchar(50);not null;uniqueIndex:uk_account_no;comment:账号编号（系统生成）"`
	ExtraInfo     string                     `json:"extra_info" gorm:"type:text;comment:额外信息"`
	AccountStatus constants.AppAccountStatus `json:"account_status" gorm:"type:tinyint;default:1;index:idx_account_status;comment:账号本身状态 1正常 2封禁 3过期 4风险"`
	CreatedAt     time.Time                  `json:"created_at" gorm:"index:idx_created_at;comment:创建时间"`
	UpdatedAt     time.Time                  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt     gorm.DeletedAt             `json:"deleted_at" gorm:"index:idx_deleted_at;comment:删除时间"`
	CreatedBy     uint                       `json:"created_by" gorm:"not null;comment:创建人ID"`
	UpdatedBy     *uint                      `json:"updated_by" gorm:"comment:更新人ID"`

	// 关联字段（不存储到数据库）
	Application *Application    `json:"application,omitempty" gorm:"foreignKey:AppID;references:AppID"`
	Category    *AppCategory    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Creator     *system.SysUser `json:"creator,omitempty" gorm:"foreignKey:CreatedBy;references:ID"`
}

// TableName 指定表名
func (AppAccount) TableName() string {
	return "app_accounts"
}

// BeforeCreate 创建前钩子
func (a *AppAccount) BeforeCreate(tx *gorm.DB) error {
	// 生成账号编号
	if a.AccountNo == "" {
		a.AccountNo = generateAccountNo()
	}
	if a.AccountDetail != "" {
		encrypted, err := crypto.EncryptAccountDetail(a.AccountDetail)
		if err != nil {
			return err
		}
		a.AccountDetail = encrypted
	}
	// 设置默认状态
	if a.AccountStatus == 0 {
		a.AccountStatus = constants.AppAccountStatusNormal
	}
	return nil
}

// BeforeUpdate 更新前加密
func (a *AppAccount) BeforeUpdate(tx *gorm.DB) error {
	// 检查 AccountDetail 是否被修改
	if tx.Statement.Changed("AccountDetail") && a.AccountDetail != "" {
		// 判断是否已经加密（简单判断，可以更严格）
		if _, err := crypto.DecryptAccountDetail(a.AccountDetail); err != nil {
			// 如果解密失败，说明是明文，需要加密
			encrypted, err := crypto.EncryptAccountDetail(a.AccountDetail)
			if err != nil {
				return err
			}
			a.AccountDetail = encrypted
		}
	}
	return nil
}

// AfterFind 查询后解密
func (a *AppAccount) AfterFind(tx *gorm.DB) error {
	if a.AccountDetail != "" {
		decrypted, err := crypto.DecryptAccountDetail(a.AccountDetail)
		if err != nil {
			// 解密失败，可能是旧数据或数据损坏
			global.GVA_LOG.Error("解密数据失败", zap.Error(err))
			return nil
		}
		a.AccountDetail = decrypted
	}
	return nil
}

// GetEncryptedDetail 获取加密的详情（用于特殊场景）
func (a *AppAccount) GetEncryptedDetail() (string, error) {
	return crypto.EncryptAccountDetail(a.AccountDetail)
}

// IsAvailable 账号是否可用
func (a *AppAccount) IsAvailable() bool {
	return a.AccountStatus == constants.AppAccountStatusNormal && a.DeletedAt.Time.IsZero()
}

// 生成账号编号
func generateAccountNo() string {
	// 格式：ACC + 年月日 + 6位随机数
	now := time.Now()
	dateStr := now.Format("20060102")
	randomNum := rand.Intn(1000000)
	return fmt.Sprintf("ACC%s%06d", dateStr, randomNum)
}
