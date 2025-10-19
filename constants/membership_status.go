package constants

// MembershipStatus 会员状态
type MembershipStatus int

const (
	MembershipStatusActive    MembershipStatus = iota + 1 // 生效中
	MembershipStatusExpired                               // 已过期
	MembershipStatusCancelled                             // 已取消
	MembershipStatusSuspended                             // 已暂停
	MembershipStatusReplaced                              // 已被替代（升级/降级)
)

// String 实现 Stringer 接口，方便日志输出和调试
func (s MembershipStatus) String() string {
	switch s {
	case MembershipStatusActive:
		return "active"
	case MembershipStatusExpired:
		return "expired"
	case MembershipStatusCancelled:
		return "cancelled"
	case MembershipStatusSuspended:
		return "suspended"
	case MembershipStatusReplaced:
		return "replaced"
	default:
		return "unknown"
	}
}

// IsValid 验证状态是否有效
func (s MembershipStatus) IsValid() bool {
	return s >= MembershipStatusActive && s <= MembershipStatusReplaced
}

// IsActive 判断是否为活跃状态
func (s MembershipStatus) IsActive() bool {
	return s == MembershipStatusActive
}
