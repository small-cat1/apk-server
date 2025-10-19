package project

import (
	"ApkAdmin/constants"
	"time"
)

// AppPackage 应用安装包表（修改后）
type AppPackage struct {
	ID            uint64                  `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	AppID         string                  `json:"app_id" gorm:"type:varchar(100);not null;comment:应用唯一标识符"`
	AppName       string                  `json:"app_name" gorm:"type:varchar(50);not null;comment:应用名称"`
	VersionName   string                  `json:"version_name" gorm:"type:varchar(50);not null;comment:版本名称（如1.0.0）"`
	CountryCode   string                  `json:"country_code" gorm:"size:3;comment:国家代码(NULL表示通用)"`
	VersionCode   *int                    `json:"version_code" gorm:"not null;index:idx_version_code;comment:版本号（用于版本比较）"`
	Platform      constants.Platform      `json:"platform" gorm:"type:enum('android','ios','harmony','windows');not null;index:idx_platform;comment:平台类型"`
	FileURL       *string                 `json:"file_url" gorm:"type:varchar(500);comment:文件下载链接"`
	Status        constants.PackageStatus `json:"status" gorm:"type:enum('draft','testing','review_pending','approved','published','rejected','suspended','archived');default:draft;index:idx_status;comment:包状态"`
	DownloadCount int                     `json:"download_count" gorm:"default:0;comment:下载次数"`
	RatingAverage *float64                `json:"rating_average" gorm:"type:decimal(3,2);comment:平均评分"`
	RatingCount   int                     `json:"rating_count" gorm:"not null;default:0;comment:评分次数"`
	UploadedAt    *time.Time              `json:"uploaded_at" gorm:"comment:上传时间"`
	PublishedAt   *time.Time              `json:"published_at" gorm:"index:idx_published_at;comment:发布时间"`
	CreatedAt     time.Time               `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt     time.Time               `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间"`
	CreatedBy     uint64                  `json:"created_by" gorm:"not null;comment:创建人ID"`
	UpdatedBy     *uint64                 `json:"updated_by" gorm:"comment:更新人ID"`
	PublishedBy   *uint64                 `json:"published_by" gorm:"comment:发布人ID"`

	// 关联关系
	Application     Application           `json:"application,omitempty" gorm:"foreignKey:AppID;references:AppID"`
	PlanRelations   []PackagePlanRelation `json:"plan_relations,omitempty" gorm:"foreignKey:AppPackageID"`
	MembershipPlans []MembershipPlan      `json:"membership_plans,omitempty" gorm:"many2many:package_plan_relations;"`
}

// TableName 指定表名
func (AppPackage) TableName() string {
	return "app_packages"
}

// PackagePlanRelation 安装包套餐关联表（新增的中间表）
type PackagePlanRelation struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	AppPackageID     uint64    `json:"app_package_id" gorm:"not null;index:idx_package_plan;comment:安装包ID"`
	MembershipPlanID uint      `json:"membership_plan_id" gorm:"not null;index:idx_package_plan;comment:套餐ID"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// 关联关系
	Package    MembershipPlan `json:"plan,omitempty" gorm:"foreignKey:MembershipPlanID"`
	AppPackage AppPackage     `json:"package,omitempty" gorm:"foreignKey:AppPackageID"`
}

// TableName 指定表名
func (PackagePlanRelation) TableName() string {
	return "package_plan_relations"
}
