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
	AppID       string                  `json:"app_id" binding:"required"`
	PlanID      string                  `json:"plan_id"`
	VersionName string                  `json:"version_name" binding:"required"`
	VersionCode int                     `json:"version_code" binding:"required"`
	Platform    constants.Platform      `json:"platform" binding:"required"`
	FileURL     string                  `json:"file_url"`
	Status      constants.PackageStatus `json:"status"`
	PlanList    []int                   `json:"planList"` // 套餐ID列表
}

func (req *AppPackageCreateRequest) Validate() error {
	// 验证应用ID格式
	if strings.TrimSpace(req.AppID) == "" {
		return errors.New("应用ID不能为空")
	}

	// 验证包名格式 (一般为反向域名格式)
	//if matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+$`, r.PackageName); !matched {
	//	return errors.New("包名格式不正确，应为反向域名格式，如：com.example.app")
	//}

	// 验证版本号必须为正整数
	if req.VersionCode <= 0 {
		return errors.New("版本号必须大于0")
	}
	return nil
}

func (req *AppPackageCreateRequest) ToAppPackage() *project.AppPackage {
	now := time.Now()
	return &project.AppPackage{
		AppID:         req.AppID,
		VersionName:   req.VersionName,
		VersionCode:   &req.VersionCode,
		Platform:      req.Platform,
		FileURL:       &req.FileURL,
		Status:        req.Status,
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
	FileURL     string             `json:"file_url,omitempty" validate:"omitempty,url,max=500"`                               // 文件URL
	// 发布设置
	Status   constants.PackageStatus `json:"status,omitempty" validate:"omitempty,oneof=draft testing review_pending approved published rejected suspended archived"` // 包状态
	PlanList []int                   `json:"planList"`                                                                                                                // 套餐ID列表

}

func (r *AppPackageUpdateRequest) Validate() error {
	// 验证ID
	if r.ID == 0 {
		return errors.New("包ID不能为空")
	}

	// 验证版本号（如果提供）
	if r.VersionCode != nil && *r.VersionCode <= 0 {
		return errors.New("版本号必须大于0")
	}

	return nil
}

type ApkBatchUpdateStatusRequest struct {
	IDs    []uint                  `json:"ids" binding:"required"`
	Status constants.PackageStatus `json:"status" binding:"required"`
}
