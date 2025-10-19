package constants

import "fmt"

type AppAccountStatus int

const (
	AppAccountStatusNormal  AppAccountStatus = 1 // 正常
	AppAccountStatusBanned  AppAccountStatus = 2 // 封禁
	AppAccountStatusExpired AppAccountStatus = 3 // 过期
	AppAccountStatusRisk    AppAccountStatus = 4 // 风险
	AppAccountStatusSold    AppAccountStatus = 5 // 已卖出
)

//`status` enum('available','locked','assigned','expired','revoked') NOT NULL DEFAULT 'available' COMMENT '状态：available-可用, locked-被订单锁定, assigned-已分配, expired-过期, revoked-已回收',

// Int 实现 Stringer 接口
func (s AppAccountStatus) Int() int {
	return int(s)
}

func (s AppAccountStatus) GetAccountStatusText() string {
	statusMap := map[AppAccountStatus]string{
		AppAccountStatusNormal:  "正常",
		AppAccountStatusBanned:  "封禁",
		AppAccountStatusExpired: "过期",
		AppAccountStatusRisk:    "风险",
		AppAccountStatusSold:    "已卖出",
	}
	if text, ok := statusMap[s]; ok {
		return text
	}
	return "未知"
}

func (s AppAccountStatus) CanTransitionTo(target AppAccountStatus) bool {
	// 定义状态转换规则
	allowedTransitions := map[AppAccountStatus][]AppAccountStatus{
		AppAccountStatusNormal: {
			AppAccountStatusBanned,  // 正常 -> 封禁
			AppAccountStatusExpired, // 正常 -> 过期
			AppAccountStatusRisk,    // 正常 -> 风险
			AppAccountStatusSold,    // 正常 -> 已卖出
		},
		AppAccountStatusBanned: {
			AppAccountStatusNormal, // 封禁 -> 正常（解封）
			AppAccountStatusRisk,   // 封禁 -> 风险
		},
		AppAccountStatusExpired: {
			AppAccountStatusNormal, // 过期 -> 正常（续期）
			AppAccountStatusBanned, // 过期 -> 封禁
		},
		AppAccountStatusRisk: {
			AppAccountStatusNormal, // 风险 -> 正常（解除风险）
			AppAccountStatusBanned, // 风险 -> 封禁
		},
		AppAccountStatusSold: {
			// 已卖出是终态，一般不允许转换
			// 如果需要退款等场景，可以添加: AppAccountStatusNormal
		},
	}

	// 检查是否允许转换
	if allowed, ok := allowedTransitions[s]; ok {
		for _, status := range allowed {
			if status == target {
				return true
			}
		}
	}
	return false
}

// ValidateTransition 验证状态转换并返回错误信息
func (s AppAccountStatus) ValidateTransition(target AppAccountStatus) error {
	if s == target {
		return fmt.Errorf("当前状态已经是 %s", target.GetAccountStatusText())
	}

	if !s.CanTransitionTo(target) {
		return fmt.Errorf("不允许从 %s 转换到 %s",
			s.GetAccountStatusText(),
			target.GetAccountStatusText())
	}

	return nil
}

// IsTerminal 判断是否为终态（不可再转换）
func (s AppAccountStatus) IsTerminal() bool {
	return s == AppAccountStatusSold
}

// IsAvailable 判断账号是否可用（可以被查询、分配等）
func (s AppAccountStatus) IsAvailable() bool {
	return s == AppAccountStatusNormal
}
