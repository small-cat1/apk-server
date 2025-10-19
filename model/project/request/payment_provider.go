package request

import (
	"ApkAdmin/model/common/request"
	"errors"
	"mime/multipart"
	"strings"
)

// PaymentProviderListRequest 支付服务商列表请求
type PaymentProviderListRequest struct {
	request.PageInfo
	Code         string `json:"code" form:"code"`                     // 服务商代码
	Name         string `json:"name" form:"name"`                     // 服务商名称
	Status       string `json:"status" form:"status"`                 // 状态
	CreatedAtGte string `json:"created_at_gte" form:"created_at_gte"` // 创建时间起始
	CreatedAtLte string `json:"created_at_lte" form:"created_at_lte"` // 创建时间结束
}

// PaymentProviderCreateRequest 创建支付服务商请求
type PaymentProviderCreateRequest struct {
	Code        string `json:"code" binding:"required"` // 服务商代码
	Name        string `json:"name" binding:"required"` // 服务商名称
	Description string `json:"description"`             // 描述
	Status      string `json:"status"`                  // 状态
	Icon        string `json:"icon"`                    // 图标URL
	SortOrder   int    `json:"sort_order"`              // 排序
}

// Validate 验证创建请求
func (r *PaymentProviderCreateRequest) Validate() error {
	if strings.TrimSpace(r.Code) == "" {
		return errors.New("服务商代码不能为空")
	}
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("服务商名称不能为空")
	}
	if r.Status != "" && r.Status != "active" && r.Status != "inactive" {
		return errors.New("状态只能是 active 或 inactive")
	}
	if r.Status == "" {
		r.Status = "active" // 默认为启用
	}
	return nil
}

// PaymentProviderUpdateRequest 更新支付服务商请求
type PaymentProviderUpdateRequest struct {
	ID           uint   `json:"id" binding:"required"`   // ID
	Code         string `json:"code" binding:"required"` // 服务商代码
	Name         string `json:"name" binding:"required"` // 服务商名称
	Description  string `json:"description"`             // 描述
	Status       string `json:"status"`                  // 状态
	Icon         string `json:"icon"`                    // 图标URL
	ConfigSchema string `json:"config_schema"`           // 配置字段模板
	SortOrder    int    `json:"sort_order"`              // 排序
}

// Validate 验证更新请求
func (r *PaymentProviderUpdateRequest) Validate() error {
	if r.ID == 0 {
		return errors.New("ID不能为空")
	}
	if strings.TrimSpace(r.Code) == "" {
		return errors.New("服务商代码不能为空")
	}
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("服务商名称不能为空")
	}
	if r.Status != "" && r.Status != "active" && r.Status != "inactive" {
		return errors.New("状态只能是 active 或 inactive")
	}
	return nil
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // ID列表
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1"`                    // ID列表
	Status string `json:"status" binding:"required,oneof=active inactive"` // 状态
}

// CheckCodeRequest 检查代码可用性请求
type CheckCodeRequest struct {
	Code string `json:"code" form:"code" binding:"required"` // 服务商代码
	ID   uint   `json:"id" form:"id"`                        // 排除的ID（用于更新时检查）
}

// ValidateConfigRequest 验证配置请求
type ValidateConfigRequest struct {
	ProviderID uint        `json:"provider_id" binding:"required"` // 服务商ID
	Config     interface{} `json:"config" binding:"required"`      // 配置数据
}

// UpdateSortOrderRequest 更新排序请求
type UpdateSortOrderRequest struct {
	ID        uint `json:"id" binding:"required"`         // ID
	SortOrder *int `json:"sort_order" binding:"required"` // 排序值
}

// CloneProviderRequest 克隆服务商请求
type CloneProviderRequest struct {
	SourceID    uint   `json:"source_id" binding:"required"` // 源服务商ID
	NewCode     string `json:"new_code" binding:"required"`  // 新服务商代码
	NewName     string `json:"new_name" binding:"required"`  // 新服务商名称
	CopyConfig  bool   `json:"copy_config"`                  // 是否复制配置
	Description string `json:"description"`                  // 描述
}

// Validate 验证克隆请求
func (r *CloneProviderRequest) Validate() error {
	if r.SourceID == 0 {
		return errors.New("源服务商ID不能为空")
	}
	if strings.TrimSpace(r.NewCode) == "" {
		return errors.New("新服务商代码不能为空")
	}
	if strings.TrimSpace(r.NewName) == "" {
		return errors.New("新服务商名称不能为空")
	}
	return nil
}

// StatisticsRequest 统计信息请求
type StatisticsRequest struct {
	ProviderID uint   `json:"provider_id" form:"provider_id"` // 特定服务商ID
	StartDate  string `json:"start_date" form:"start_date"`   // 开始日期
	EndDate    string `json:"end_date" form:"end_date"`       // 结束日期
	MetricType string `json:"metric_type" form:"metric_type"` // 指标类型
	GroupBy    string `json:"group_by" form:"group_by"`       // 分组方式
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
	ID     uint        `json:"id" binding:"required"`     // 服务商ID
	Config interface{} `json:"config" binding:"required"` // 测试配置
}

// ProviderLogsRequest 服务商日志请求
type ProviderLogsRequest struct {
	request.PageInfo
	ProviderID uint   `json:"provider_id" form:"provider_id" binding:"required"` // 服务商ID
	LogLevel   string `json:"log_level" form:"log_level"`                        // 日志级别
	StartTime  string `json:"start_time" form:"start_time"`                      // 开始时间
	EndTime    string `json:"end_time" form:"end_time"`                          // 结束时间
	Keyword    string `json:"keyword" form:"keyword"`                            // 关键词搜索
}

// ConfigHistoryRequest 配置历史请求
type ConfigHistoryRequest struct {
	request.PageInfo
	ProviderID uint   `json:"provider_id" form:"provider_id" binding:"required"` // 服务商ID
	StartTime  string `json:"start_time" form:"start_time"`                      // 开始时间
	EndTime    string `json:"end_time" form:"end_time"`                          // 结束时间
	Action     string `json:"action" form:"action"`                              // 操作类型
}

// RestoreConfigRequest 恢复配置请求
type RestoreConfigRequest struct {
	ProviderID uint `json:"provider_id" binding:"required"` // 服务商ID
	VersionID  uint `json:"version_id" binding:"required"`  // 版本ID
}

// Validate 验证恢复配置请求
func (r *RestoreConfigRequest) Validate() error {
	if r.ProviderID == 0 {
		return errors.New("服务商ID不能为空")
	}
	if r.VersionID == 0 {
		return errors.New("版本ID不能为空")
	}
	return nil
}

// ImportConfigRequest 导入配置请求
type ImportConfigRequest struct {
	File         *multipart.FileHeader `json:"file" binding:"required"` // 导入文件
	Override     bool                  `json:"override"`                // 是否覆盖已存在的配置
	ValidateOnly bool                  `json:"validate_only"`           // 是否仅验证不导入
}

// SyncProviderRequest 同步服务商信息请求
type SyncProviderRequest struct {
	ID       uint   `json:"id" binding:"required"` // 服务商ID
	SyncType string `json:"sync_type"`             // 同步类型: basic, config, methods, all
}

// ProviderPaymentMethodsRequest 获取服务商支付方式请求
type ProviderPaymentMethodsRequest struct {
	ProviderID uint   `json:"provider_id" form:"provider_id" binding:"required"` // 服务商ID
	Currency   string `json:"currency" form:"currency"`                          // 货币类型
	Country    string `json:"country" form:"country"`                            // 国家代码
}

// ExportConfigRequest 导出配置请求
type ExportConfigRequest struct {
	ID         uint     `json:"id" binding:"required"` // 服务商ID
	ExportType string   `json:"export_type"`           // 导出类型: json, yaml, xml
	Fields     []string `json:"fields"`                // 导出字段
}

// ClearCacheRequest 清理缓存请求
type ClearCacheRequest struct {
	ID        uint   `json:"id" binding:"required"` // 服务商ID
	CacheType string `json:"cache_type"`            // 缓存类型: config, methods, statistics, all
}

// EnableProviderRequest 启用服务商请求
type EnableProviderRequest struct {
	ID uint `json:"id" binding:"required"` // 服务商ID
}

// DisableProviderRequest 禁用服务商请求
type DisableProviderRequest struct {
	ID     uint   `json:"id" binding:"required"` // 服务商ID
	Reason string `json:"reason"`                // 禁用原因
}

// SearchProviderRequest 搜索服务商请求
type SearchProviderRequest struct {
	request.PageInfo
	Keyword    string   `json:"keyword" form:"keyword"`       // 搜索关键词
	Status     string   `json:"status" form:"status"`         // 状态筛选
	Categories []string `json:"categories" form:"categories"` // 分类筛选
	SortBy     string   `json:"sort_by" form:"sort_by"`       // 排序字段
}

// ProviderHealthCheckRequest 服务商健康检查请求
type ProviderHealthCheckRequest struct {
	ID          uint `json:"id" binding:"required"` // 服务商ID
	CheckConfig bool `json:"check_config"`          // 是否检查配置
	CheckAPI    bool `json:"check_api"`             // 是否检查API连通性
}

// BulkOperationRequest 批量操作请求
type BulkOperationRequest struct {
	IDs       []uint                 `json:"ids" binding:"required,min=1"` // ID列表
	Operation string                 `json:"operation" binding:"required"` // 操作类型
	Data      map[string]interface{} `json:"data"`                         // 操作数据
}

// Validate 验证批量操作请求
func (r *BulkOperationRequest) Validate() error {
	if len(r.IDs) == 0 {
		return errors.New("请选择要操作的数据")
	}

	validOperations := []string{"enable", "disable", "delete", "update_sort", "sync", "clear_cache"}
	isValid := false
	for _, op := range validOperations {
		if r.Operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("无效的操作类型")
	}

	return nil
}

// GetProviderDetailsRequest 获取服务商详细信息请求
type GetProviderDetailsRequest struct {
	ID            uint `json:"id" form:"id" binding:"required"`      // 服务商ID
	IncludeConfig bool `json:"include_config" form:"include_config"` // 是否包含配置信息
	IncludeLogs   bool `json:"include_logs" form:"include_logs"`     // 是否包含日志信息
	IncludeStats  bool `json:"include_stats" form:"include_stats"`   // 是否包含统计信息
}

// UpdateProviderConfigRequest 更新服务商配置请求
type UpdateProviderConfigRequest struct {
	ID     uint        `json:"id" binding:"required"`     // 服务商ID
	Config interface{} `json:"config" binding:"required"` // 配置数据
	Reason string      `json:"reason"`                    // 更新原因
}

// Validate 验证更新配置请求
func (r *UpdateProviderConfigRequest) Validate() error {
	if r.ID == 0 {
		return errors.New("服务商ID不能为空")
	}
	if r.Config == nil {
		return errors.New("配置数据不能为空")
	}
	return nil
}
