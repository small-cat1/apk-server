package project

import "time"

type UserCommissionAccount struct {
	ID              int64     `gorm:"primarykey;comment:账户ID" json:"id"`
	UserID          uint      `gorm:"not null;uniqueIndex;comment:用户ID" json:"userId"`
	AvailableAmount float64   `gorm:"type:decimal(10,2);default:0.00;comment:可提现金额" json:"availableAmount"`
	FrozenAmount    float64   `gorm:"type:decimal(10,2);default:0.00;comment:冻结金额" json:"frozenAmount"`
	TotalEarnings   float64   `gorm:"type:decimal(10,2);default:0.00;comment:累计收益" json:"totalEarnings"`
	WithdrawnAmount float64   `gorm:"type:decimal(10,2);default:0.00;comment:已提现金额" json:"withdrawnAmount"`
	CreatedAt       time.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"comment:更新时间" json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserCommissionAccount) TableName() string {
	return "user_commission_account"
}

// UserCommissionAccountSimple 创建一个只包含需要字段的结构体
type UserCommissionAccountSimple struct {
	ID              int64   `json:"id"`
	UserID          uint    `json:"user_id"`
	AvailableAmount float64 `json:"available_amount"`
}

func (UserCommissionAccountSimple) TableName() string {
	return "user_commission_account"
}
