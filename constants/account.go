package constants

import (
	"database/sql/driver"
	"fmt"
)

// AccountStatus 账户状态类型
type AccountStatus int

// 账户状态常量
const (
	AccountStatusNormal    AccountStatus = 1 // 正常
	AccountStatusDisabled  AccountStatus = 2 // 已禁用（管理员操作）
	AccountStatusLocked    AccountStatus = 3 // 临时锁定（登录失败过多）
	AccountStatusPending   AccountStatus = 4 // 待审核
	AccountStatusSuspended AccountStatus = 5 // 暂停使用（违规等）
	AccountStatusDeleted   AccountStatus = 9 // 已删除（软删除）
)

// statusNames 状态名称映射
var statusNames = map[AccountStatus]string{
	AccountStatusNormal:    "正常",
	AccountStatusDisabled:  "已禁用",
	AccountStatusLocked:    "临时锁定",
	AccountStatusPending:   "待审核",
	AccountStatusSuspended: "暂停使用",
	AccountStatusDeleted:   "已删除",
}

// statusColors 状态颜色映射（用于前端展示）
var statusColors = map[AccountStatus]string{
	AccountStatusNormal:    "success",
	AccountStatusDisabled:  "danger",
	AccountStatusLocked:    "warning",
	AccountStatusPending:   "info",
	AccountStatusSuspended: "danger",
	AccountStatusDeleted:   "default",
}

// String 实现 Stringer 接口
func (s AccountStatus) String() string {
	if name, ok := statusNames[s]; ok {
		return name
	}
	return fmt.Sprintf("未知状态(%d)", s)
}

// IsValid 检查状态是否有效
func (s AccountStatus) IsValid() bool {
	_, ok := statusNames[s]
	return ok
}

// IsNormal 是否为正常状态
func (s AccountStatus) IsNormal() bool {
	return s == AccountStatusNormal
}

// IsActive 是否为活跃状态（可以登录）
func (s AccountStatus) IsActive() bool {
	return s == AccountStatusNormal
}

// IsBlocked 是否被阻止登录
func (s AccountStatus) IsBlocked() bool {
	return s == AccountStatusDisabled ||
		s == AccountStatusLocked ||
		s == AccountStatusSuspended ||
		s == AccountStatusDeleted
}

// CanLogin 是否可以登录
func (s AccountStatus) CanLogin() bool {
	return !s.IsBlocked()
}

// IsTemporary 是否为临时状态（有过期时间）
func (s AccountStatus) IsTemporary() bool {
	return s == AccountStatusLocked || s == AccountStatusSuspended
}

// Color 获取状态对应的颜色
func (s AccountStatus) Color() string {
	if color, ok := statusColors[s]; ok {
		return color
	}
	return "default"
}

// GetBlockReason 获取阻止原因（用于错误提示）
func (s AccountStatus) GetBlockReason() string {
	switch s {
	case AccountStatusDisabled:
		return "账户已被禁用，请联系管理员"
	case AccountStatusLocked:
		return "账户已被临时锁定"
	case AccountStatusPending:
		return "账户待审核，请等待管理员审核"
	case AccountStatusSuspended:
		return "账户已被暂停使用"
	case AccountStatusDeleted:
		return "账户已注销"
	default:
		return "账户状态异常"
	}
}

// Scan 实现 sql.Scanner 接口（从数据库读取）
func (s *AccountStatus) Scan(value interface{}) error {
	if value == nil {
		*s = AccountStatusNormal
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = AccountStatus(v)
	case int:
		*s = AccountStatus(v)
	case []byte:
		var i int
		if _, err := fmt.Sscan(string(v), &i); err != nil {
			return err
		}
		*s = AccountStatus(i)
	default:
		return fmt.Errorf("无法将 %T 转换为 AccountStatus", value)
	}

	return nil
}

// Value 实现 driver.Valuer 接口（写入数据库）
func (s AccountStatus) Value() (driver.Value, error) {
	return int64(s), nil
}

// ParseAccountStatus 从字符串或数字解析状态
func ParseAccountStatus(v interface{}) (AccountStatus, error) {
	switch val := v.(type) {
	case int:
		status := AccountStatus(val)
		if !status.IsValid() {
			return 0, fmt.Errorf("无效的状态值: %d", val)
		}
		return status, nil
	case string:
		// 支持通过名称查找
		for status, name := range statusNames {
			if name == val {
				return status, nil
			}
		}
		return 0, fmt.Errorf("无效的状态名称: %s", val)
	default:
		return 0, fmt.Errorf("不支持的类型: %T", v)
	}
}

// GetAllStatuses 获取所有可用状态
func GetAllStatuses() []AccountStatus {
	return []AccountStatus{
		AccountStatusNormal,
		AccountStatusDisabled,
		AccountStatusLocked,
		AccountStatusPending,
		AccountStatusSuspended,
		AccountStatusDeleted,
	}
}

// GetStatusOptions 获取状态选项（用于前端下拉框）
func GetStatusOptions() []map[string]interface{} {
	options := make([]map[string]interface{}, 0)
	for _, status := range GetAllStatuses() {
		options = append(options, map[string]interface{}{
			"value": int(status),
			"label": status.String(),
			"color": status.Color(),
		})
	}
	return options
}
