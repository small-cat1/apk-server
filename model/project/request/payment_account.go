package request

import "ApkAdmin/model/common/request"

// 请求结构体
type PaymentAccountListReq struct {
	PageInfo
	ProviderCode string `json:"provider_code" form:"provider_code"`
	Status       string `json:"status" form:"status"`
	Group        string `json:"group" form:"group"`
	Region       string `json:"region" form:"region"`
	Name         string `json:"name" form:"name"`
}

type CreatePaymentAccountReq struct {
	Name           string                 `json:"name" binding:"required"`
	ProviderCode   string                 `json:"provider_code" binding:"required"`
	AccountType    string                 `json:"account_type"`
	Config         map[string]interface{} `json:"config" binding:"required"`
	Status         string                 `json:"status"`
	Weight         int                    `json:"weight"`
	MaxDailyAmount float64                `json:"max_daily_amount"`
	Group          string                 `json:"group"`
	Tags           string                 `json:"tags"`
	Region         string                 `json:"region"`
	Remark         string                 `json:"remark"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// PaymentAccountUpdate 更新支付账号请求
type PaymentAccountUpdate struct {
	ID             uint                   `json:"id" binding:"required" comment:"账号ID"`
	Name           string                 `json:"name" binding:"required" comment:"账号名称"`
	Config         map[string]interface{} `json:"config" binding:"required" comment:"支付配置参数"`
	Status         string                 `json:"status" binding:"omitempty,oneof=active inactive maintenance deleted" comment:"状态"`
	Weight         int                    `json:"weight" binding:"omitempty,min=1,max=100" comment:"权重"`
	MaxDailyAmount float64                `json:"max_daily_amount" binding:"omitempty,min=0" comment:"日最大交易金额"`
	Group          string                 `json:"group" binding:"omitempty,max=50" comment:"分组"`
	Region         string                 `json:"region" binding:"omitempty,max=50" comment:"地区"`
	Tags           string                 `json:"tags" binding:"omitempty,max=255" comment:"标签"`
	Remark         string                 `json:"remark" binding:"omitempty,max=500" comment:"备注"`
}

// PaymentAccountSearch 支付账号搜索请求
type PaymentAccountSearch struct {
	request.PageInfo
	Name         string `form:"name" comment:"账号名称"`
	ProviderCode string `form:"provider_code" comment:"支付服务商代码"`
	AccountType  string `form:"account_type" comment:"账号类型"`
	Status       string `form:"status" comment:"状态"`
	Group        string `form:"group" comment:"分组"`
	Region       string `form:"region" comment:"地区"`
	StartDate    string `form:"start_date" comment:"创建开始日期"`
	EndDate      string `form:"end_date" comment:"创建结束日期"`
}

// PaymentAccountById 根据ID查询支付账号请求
type PaymentAccountById struct {
	ID uint `form:"id" json:"id" binding:"required" comment:"账号ID"`
}

// PaymentAccountDelete 删除支付账号请求
type PaymentAccountDelete struct {
	ID uint `json:"id" binding:"required" comment:"账号ID"`
}

// PaymentAccountBatchDelete 批量删除支付账号请求
type PaymentAccountBatchDelete struct {
	IDs []uint `json:"ids" binding:"required,min=1" comment:"账号ID列表"`
}

// PaymentAccountBatchStatus 批量更新支付账号状态请求
type PaymentAccountBatchStatus struct {
	IDs    []uint `json:"ids" binding:"required,min=1" comment:"账号ID列表"`
	Status string `json:"status" binding:"required,oneof=active inactive maintenance deleted" comment:"状态"`
}

// PaymentAccountWeight 更新账号权重请求
type PaymentAccountWeight struct {
	ID     uint `json:"id" binding:"required" comment:"账号ID"`
	Weight int  `json:"weight" binding:"required,min=1,max=100" comment:"权重"`
}

// PaymentAccountMaintenance 设置账号维护模式请求
type PaymentAccountMaintenance struct {
	ID          uint `json:"id" binding:"required" comment:"账号ID"`
	Maintenance bool `json:"maintenance" binding:"required" comment:"是否维护模式"`
}

// PaymentAccountConfigValidate 验证账号配置请求
type PaymentAccountConfigValidate struct {
	ProviderCode string                 `json:"provider_code" binding:"required" comment:"支付服务商代码"`
	Config       map[string]interface{} `json:"config" binding:"required" comment:"配置参数"`
}

// PaymentAccountStatistics 获取账号交易统计请求
type PaymentAccountStatistics struct {
	ID        uint   `form:"id" comment:"账号ID"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
	TimeRange string `form:"time_range" comment:"时间范围：today,week,month,year"`
}

// PaymentAccountTransactions 获取账号交易记录请求
type PaymentAccountTransactions struct {
	request.PageInfo
	AccountID uint   `form:"account_id" binding:"required" comment:"账号ID"`
	Status    string `form:"status" comment:"交易状态"`
	OrderNo   string `form:"order_no" comment:"订单号"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
	MinAmount string `form:"min_amount" comment:"最小金额"`
	MaxAmount string `form:"max_amount" comment:"最大金额"`
}

// PaymentAccountOperationLogs 获取账号操作日志请求
type PaymentAccountOperationLogs struct {
	request.PageInfo
	AccountID  uint   `form:"account_id" comment:"账号ID"`
	Operation  string `form:"operation" comment:"操作类型"`
	OperatorID uint   `form:"operator_id" comment:"操作人ID"`
	StartDate  string `form:"start_date" comment:"开始日期"`
	EndDate    string `form:"end_date" comment:"结束日期"`
}

// PaymentAccountGroupCreate 创建账号分组请求
type PaymentAccountGroupCreate struct {
	Name        string `json:"name" binding:"required,max=50" comment:"分组名称"`
	Description string `json:"description" binding:"omitempty,max=200" comment:"分组描述"`
	SortOrder   int    `json:"sort_order" binding:"omitempty,min=0" comment:"排序"`
}

// PaymentAccountGroupUpdate 更新账号分组请求
type PaymentAccountGroupUpdate struct {
	ID          uint   `json:"id" binding:"required" comment:"分组ID"`
	Name        string `json:"name" binding:"required,max=50" comment:"分组名称"`
	Description string `json:"description" binding:"omitempty,max=200" comment:"分组描述"`
	SortOrder   int    `json:"sort_order" binding:"omitempty,min=0" comment:"排序"`
}

// PaymentAccountGroupDelete 删除账号分组请求
type PaymentAccountGroupDelete struct {
	ID uint `json:"id" binding:"required" comment:"分组ID"`
}

// PaymentAccountBatchAssignGroup 批量分配账号到分组请求
type PaymentAccountBatchAssignGroup struct {
	IDs     []uint `json:"ids" binding:"required,min=1" comment:"账号ID列表"`
	GroupID uint   `json:"group_id" binding:"required" comment:"分组ID"`
}

// PaymentAccountBatchTags 批量标签操作请求
type PaymentAccountBatchTags struct {
	IDs  []uint   `json:"ids" binding:"required,min=1" comment:"账号ID列表"`
	Tags []string `json:"tags" binding:"required,min=1" comment:"标签列表"`
}

// PaymentAccountUsageStatistics 获取账号使用率统计请求
type PaymentAccountUsageStatistics struct {
	TimeRange    string `form:"time_range" comment:"时间范围：today,week,month,year"`
	StartDate    string `form:"start_date" comment:"开始日期"`
	EndDate      string `form:"end_date" comment:"结束日期"`
	ProviderCode string `form:"provider_code" comment:"支付服务商代码"`
	Group        string `form:"group" comment:"分组"`
}

// PaymentAccountPopularMethods 获取热门支付方式统计请求
type PaymentAccountPopularMethods struct {
	TimeRange string `form:"time_range" comment:"时间范围"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
	Limit     int    `form:"limit" comment:"返回数量限制"`
}

// PaymentAccountRegionDistribution 获取地区分布统计请求
type PaymentAccountRegionDistribution struct {
	TimeRange string `form:"time_range" comment:"时间范围"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
}

// PaymentAccountPerformanceMetrics 获取账号性能指标请求
type PaymentAccountPerformanceMetrics struct {
	AccountID uint   `form:"account_id" comment:"账号ID"`
	TimeRange string `form:"time_range" comment:"时间范围"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
}

// PaymentAccountAlertRules 设置账号告警规则请求
type PaymentAccountAlertRules struct {
	AccountID             uint    `json:"account_id" binding:"required" comment:"账号ID"`
	DailyAmountThreshold  float64 `json:"daily_amount_threshold" binding:"omitempty,min=0" comment:"日交易金额阈值"`
	DailyOrdersThreshold  int     `json:"daily_orders_threshold" binding:"omitempty,min=0" comment:"日订单数阈值"`
	FailureRateThreshold  float64 `json:"failure_rate_threshold" binding:"omitempty,min=0,max=1" comment:"失败率阈值"`
	ResponseTimeThreshold int     `json:"response_time_threshold" binding:"omitempty,min=0" comment:"响应时间阈值(毫秒)"`
	BalanceThreshold      float64 `json:"balance_threshold" binding:"omitempty,min=0" comment:"余额阈值"`
	NotificationEmails    string  `json:"notification_emails" binding:"omitempty" comment:"通知邮箱，多个用逗号分隔"`
	NotificationWebhooks  string  `json:"notification_webhooks" binding:"omitempty" comment:"通知Webhook，多个用逗号分隔"`
	EnableEmailAlert      bool    `json:"enable_email_alert" comment:"是否启用邮件告警"`
	EnableWebhookAlert    bool    `json:"enable_webhook_alert" comment:"是否启用Webhook告警"`
	EnableSMSAlert        bool    `json:"enable_sms_alert" comment:"是否启用短信告警"`
}

// PaymentAccountAlerts 获取账号告警记录请求
type PaymentAccountAlerts struct {
	request.PageInfo
	AccountID uint   `form:"account_id" comment:"账号ID"`
	AlertType string `form:"alert_type" comment:"告警类型"`
	Status    string `form:"status" comment:"告警状态：pending,acknowledged,resolved"`
	Severity  string `form:"severity" comment:"告警级别：low,medium,high,critical"`
	StartDate string `form:"start_date" comment:"开始日期"`
	EndDate   string `form:"end_date" comment:"结束日期"`
}

// PaymentAccountAlertAcknowledge 确认告警请求
type PaymentAccountAlertAcknowledge struct {
	AlertID         uint   `json:"alert_id" binding:"required" comment:"告警ID"`
	AcknowledgedBy  uint   `json:"acknowledged_by" binding:"required" comment:"确认人ID"`
	AcknowledgeNote string `json:"acknowledge_note" binding:"omitempty,max=500" comment:"确认备注"`
}

// PaymentAccountBackups 获取账号备份列表请求
type PaymentAccountBackups struct {
	request.PageInfo
	AccountID  uint   `form:"account_id" comment:"账号ID"`
	BackupType string `form:"backup_type" comment:"备份类型：manual,auto,scheduled"`
	StartDate  string `form:"start_date" comment:"开始日期"`
	EndDate    string `form:"end_date" comment:"结束日期"`
}

// PaymentAccountBackupCreate 创建账号备份请求
type PaymentAccountBackupCreate struct {
	AccountID   uint   `json:"account_id" binding:"required" comment:"账号ID"`
	BackupType  string `json:"backup_type" binding:"required,oneof=manual auto scheduled" comment:"备份类型"`
	Description string `json:"description" binding:"omitempty,max=200" comment:"备份描述"`
}

// PaymentAccountBackupRestore 恢复账号备份请求
type PaymentAccountBackupRestore struct {
	BackupID  uint `json:"backup_id" binding:"required" comment:"备份ID"`
	AccountID uint `json:"account_id" binding:"required" comment:"目标账号ID"`
}

// PaymentAccountBackupDelete 删除账号备份请求
type PaymentAccountBackupDelete struct {
	BackupID uint `json:"backup_id" binding:"required" comment:"备份ID"`
}

// PaymentAccountExport 导出账号数据请求
type PaymentAccountExport struct {
	Format       string   `form:"format" binding:"omitempty,oneof=excel csv json" comment:"导出格式"`
	AccountIDs   []uint   `form:"account_ids" comment:"指定账号ID列表"`
	ProviderCode string   `form:"provider_code" comment:"支付服务商代码"`
	Status       string   `form:"status" comment:"状态"`
	Group        string   `form:"group" comment:"分组"`
	Fields       []string `form:"fields" comment:"导出字段"`
	StartDate    string   `form:"start_date" comment:"开始日期"`
	EndDate      string   `form:"end_date" comment:"结束日期"`
}

// PaymentAccountImport 导入账号数据请求
type PaymentAccountImport struct {
	OverwriteExisting bool `form:"overwrite_existing" comment:"是否覆盖已存在的账号"`
	ValidateOnly      bool `form:"validate_only" comment:"仅验证不导入"`
}

// PaymentAccountHealthCheck 账号健康检查请求
type PaymentAccountHealthCheck struct {
	AccountID uint   `json:"account_id" binding:"required" comment:"账号ID"`
	CheckType string `json:"check_type" binding:"omitempty,oneof=connection balance config all" comment:"检查类型"`
}

// PaymentAccountSyncBalance 同步账号余额请求
type PaymentAccountSyncBalance struct {
	AccountID uint `json:"account_id" binding:"required" comment:"账号ID"`
}

// PaymentAccountTestConnection 测试账号连接请求
type PaymentAccountTestConnection struct {
	AccountID uint `json:"account_id" binding:"required" comment:"账号ID"`
}
