package response

import "ApkAdmin/model/project"

// CommissionStats 佣金统计信息
type CommissionStats struct {
	TotalCommission float64 `json:"totalCommission"` // 累计佣金
	TotalOrders     int64   `json:"totalOrders"`     // 订单数
}

// CommissionDetailListResponse 分佣明细列表响应（包含分页数据和统计信息）
type CommissionDetailListResponse struct {
	List     []project.CommissionDetail `json:"list"`     // 明细列表
	Total    int64                      `json:"total"`    // 总记录数
	Page     int                        `json:"page"`     // 当前页码
	PageSize int                        `json:"pageSize"` // 每页大小
	Stats    CommissionStats            `json:"stats"`    // 统计信息
}
