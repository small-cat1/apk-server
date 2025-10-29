package request

// CommissionDetailSearch 分页查询请求参数
type CommissionDetailSearch struct {
	ClientCommissionDetailSearch
	UserId *uint `json:"userId" form:"userId"` // 用户ID
}

// ClientCommissionDetailSearch  分页查询请求参数
type ClientCommissionDetailSearch struct {
	Status     string `json:"status" form:"status" binding:"omitempty,oneof=all pending settled frozen"`   // 状态筛选：all-全部, pending-待结算, settled-已结算, frozen-冻结
	TimeFilter string `json:"timeFilter" form:"timeFilter" binding:"omitempty,oneof=all today week month"` // 时间筛选：all-全部, today-今天, week-本周, month-本月
	PageInfo          // 分页参数
}
