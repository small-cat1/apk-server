package response

import "time"

// WithdrawRecordResp 提现记录响应
type WithdrawRecordResp struct {
	ID           int64      `json:"id"`
	WithdrawNo   string     `json:"withdrawNo"`             // 提现单号
	Amount       float64    `json:"amount"`                 // 提现金额
	Fee          float64    `json:"fee"`                    // 手续费
	ActualAmount float64    `json:"actualAmount"`           // 实际到账金额
	WithdrawType string     `json:"withdrawType"`           // 提现方式
	AccountName  *string    `json:"accountName,omitempty"`  // 账户名
	AccountNo    *string    `json:"accountNo,omitempty"`    // 账户号（脱敏）
	Status       string     `json:"status"`                 // 状态
	RejectReason *string    `json:"rejectReason,omitempty"` // 拒绝原因
	AuditTime    *time.Time `json:"auditTime,omitempty"`    // 审核时间
	CompleteTime *time.Time `json:"completeTime,omitempty"` // 完成时间
	Remark       *string    `json:"remark,omitempty"`       // 备注
	CreateTime   time.Time  `json:"createTime"`             // 创建时间
	UpdateTime   time.Time  `json:"updateTime"`             // 更新时间
}

// WithdrawRecordListResp 提现记录列表响应
type WithdrawRecordListResp struct {
	List           []WithdrawRecordResp `json:"list"`
	Total          int64                `json:"total"`
	TotalWithdrawn float64              `json:"totalWithdrawn"` // 累计提现金额
	TotalCount     int64                `json:"totalCount"`     // 提现次数
}

// MaskAccountNo 账号脱敏处理
func (r *WithdrawRecordResp) MaskAccountNo() {
	if r.AccountNo != nil && *r.AccountNo != "" {
		accountNo := *r.AccountNo
		length := len(accountNo)

		if length <= 4 {
			// 如果长度小于等于4，只显示最后一位
			masked := "***" + accountNo[length-1:]
			r.AccountNo = &masked
		} else if length <= 8 {
			// 长度5-8位，显示前后各1位
			masked := accountNo[:1] + "***" + accountNo[length-1:]
			r.AccountNo = &masked
		} else {
			// 长度大于8位，显示前后各2位
			masked := accountNo[:2] + "******" + accountNo[length-2:]
			r.AccountNo = &masked
		}
	}
}
