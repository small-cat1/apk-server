package request

import (
	"ApkAdmin/constants"
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"ApkAdmin/utils"
	"errors"
	"regexp"
	"strings"
	"time"
)

// AppPackageListRequest  国家地区列表请求
type AppPackageListRequest struct {
	request.PageInfo
	AppID       string `form:"app_id" json:"app_id"`             // 应用ID
	Platform    string `form:"platform" json:"platform"`         // 平台筛选
	Status      string `form:"status" json:"status"`             // 状态筛选
	CountryCode string `form:"country_code" json:"country_code"` // 国家代码筛选
}

func (r *AppPackageListRequest) Validate() error {
	// 验证平台
	if r.Platform != "" {
		validPlatforms := []string{"android", "ios", "harmony", "windows"}
		if !utils.Contains(validPlatforms, r.Platform) {
			return errors.New("无效的平台类型")
		}
	}

	// 验证状态
	if r.Status != "" {
		validStatuses := []string{"draft", "testing", "review_pending", "approved", "published", "rejected", "suspended", "archived"}
		if !utils.Contains(validStatuses, r.Status) {
			return errors.New("无效的状态")
		}
	}

	// 验证国家代码
	if r.CountryCode != "" {
		if len(r.CountryCode) < 2 || len(r.CountryCode) > 4 {
			return errors.New("国家代码长度必须为2-3位")
		}
		if matched, _ := regexp.MatchString(`^[A-Z]{2,4}$`, r.CountryCode); !matched {
			return errors.New("国家代码格式不正确")
		}
	}

	return nil
}

type AppPackageCreateRequest struct {
	AppID       string             `json:"app_id" binding:"required"`
	PlanID      string             `json:"plan_id"`
	VersionName string             `json:"version_name" binding:"required"`
	VersionCode int                `json:"version_code" binding:"required"`
	Platform    constants.Platform `json:"platform" binding:"required"`
	FileURL     string             `json:"file_url,omitempty" validate:"omitempty,max=500"`    // URL格式
	ObjectName  string             `json:"object_name,omitempty" validate:"omitempty,max=500"` // OSS路径（不是URL！）
	FileName    string             `json:"file_name,omitempty" validate:"omitempty,max=255"`   // 文件名（新增）
	PackageSize int64              `json:"package_size,omitempty" validate:"omitempty,min=0"`  // 非负数

	Status   string `json:"status"`
	PlanList []int  `json:"planList"` // 套餐ID列表
}

func (r *AppPackageCreateRequest) Validate() error {
	// 验证应用ID格式
	if strings.TrimSpace(r.AppID) == "" {
		return errors.New("应用ID不能为空")
	}

	// 验证版本号必须为正整数
	if r.VersionCode <= 0 {
		return errors.New("版本号必须大于0")
	}

	// 验证ObjectName（重要！）
	if r.ObjectName != "" {
		if len(r.ObjectName) > 500 {
			return errors.New("OSS对象名称长度不能超过500个字符")
		}
		if strings.Contains(r.ObjectName, "..") {
			return errors.New("OSS对象名称格式不正确")
		}
		if strings.HasPrefix(r.ObjectName, "/") {
			return errors.New("OSS对象名称不能以'/'开头")
		}
		// 必须以 private/ 或 public/ 开头
		validPrefixes := []string{"private/", "public/"}
		hasValidPrefix := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(r.ObjectName, prefix) {
				hasValidPrefix = true
				break
			}
		}
		if !hasValidPrefix {
			return errors.New("OSS对象名称必须以'private/'或'public/'开头")
		}
	}
	// 验证PackageSize
	if r.PackageSize < 0 {
		return errors.New("文件大小不能为负数")
	}
	maxSize := int64(500 * 1024 * 1024) // 500MB
	if r.PackageSize > maxSize {
		return errors.New("文件大小不能超过500MB")
	}
	return nil
}

func (r *AppPackageCreateRequest) ToAppPackage() *project.AppPackage {
	now := time.Now()
	return &project.AppPackage{
		AppID:         r.AppID,
		VersionName:   r.VersionName,
		VersionCode:   &r.VersionCode,
		Platform:      r.Platform,
		FileURL:       &r.FileURL,
		FileName:      &r.FileName,
		ObjectName:    &r.ObjectName,
		PackageSize:   r.PackageSize,
		Status:        r.Status,
		DownloadCount: 0,
		RatingCount:   0,
		UploadedAt:    &now,
		CreatedAt:     now,
		UpdatedAt:     now,
		// CreatedBy 需要从上下文中获取，这里暂时设为0
		CreatedBy: 0,
	}
}

// AppPackageUpdateRequest 更新应用安装包请求
type AppPackageUpdateRequest struct {
	ID uint `json:"id" binding:"required" validate:"required,min=1"` // 包ID (必填)
	// 可更新的基础信息
	VersionName string             `json:"version_name,omitempty" validate:"omitempty,min=1,max=50"`                          // 版本名称
	VersionCode *int               `json:"version_code,omitempty" validate:"omitempty,min=1"`                                 // 版本号
	Platform    constants.Platform `json:"platform" binding:"required" validate:"required,oneof=android ios harmony windows"` // 平台
	FileURL     string             `json:"file_url,omitempty" validate:"omitempty,max=500"`                                   // URL格式
	ObjectName  string             `json:"object_name,omitempty" validate:"omitempty,max=500"`                                // OSS路径（不是URL！）
	FileName    string             `json:"file_name,omitempty" validate:"omitempty,max=255"`                                  // 文件名（新增）
	PackageSize int64              `json:"package_size,omitempty" validate:"omitempty,min=0"`                                 // 非负数
	// 发布设置
	Status   constants.PackageStatus `json:"status,omitempty" validate:"omitempty,oneof=draft testing review_pending approved published rejected suspended archived"` // 包状态
	PlanList []int                   `json:"planList"`                                                                                                                // 套餐ID列表
}

func (r *AppPackageUpdateRequest) Validate() error {
	// 验证ID
	if r.ID <= 0 {
		return errors.New("包ID不能为空")
	}
	// 验证版本号（如果提供）
	if r.VersionCode != nil && *r.VersionCode <= 0 {
		return errors.New("版本号必须大于0")
	}
	// 验证ObjectName（重要！）
	if r.ObjectName != "" {
		if len(r.ObjectName) > 500 {
			return errors.New("OSS对象名称长度不能超过500个字符")
		}
		if strings.Contains(r.ObjectName, "..") {
			return errors.New("OSS对象名称格式不正确")
		}
		if strings.HasPrefix(r.ObjectName, "/") {
			return errors.New("OSS对象名称不能以'/'开头")
		}
		// 必须以 private/ 或 public/ 开头
		validPrefixes := []string{"private/", "public/"}
		hasValidPrefix := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(r.ObjectName, prefix) {
				hasValidPrefix = true
				break
			}
		}
		if !hasValidPrefix {
			return errors.New("OSS对象名称必须以'private/'或'public/'开头")
		}
	}
	// 验证PackageSize
	if r.PackageSize < 0 {
		return errors.New("文件大小不能为负数")
	}
	maxSize := int64(500 * 1024 * 1024) // 500MB
	if r.PackageSize > maxSize {
		return errors.New("文件大小不能超过500MB")
	}
	return nil
}

type ApkBatchUpdateStatusRequest struct {
	IDs    []uint                  `json:"ids" binding:"required"`
	Status constants.PackageStatus `json:"status" binding:"required"`
}
