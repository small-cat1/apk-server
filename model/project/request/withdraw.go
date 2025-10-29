package request

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// UserWithdrawRequest 用户提现请求
type UserWithdrawRequest struct {
	Amount        float64 `json:"amount"`        // 提现金额
	WithdrawType  string  `json:"withdrawType"`  // 提现方式：alipay/wechat/bank
	AlipayAccount string  `json:"alipayAccount"` // 支付宝账号
	AlipayName    string  `json:"alipayName"`    // 支付宝姓名
	WechatAccount string  `json:"wechatAccount"` // 微信账号
	WechatName    string  `json:"wechatName"`    // 微信姓名
	BankName      string  `json:"bankName"`      // 开户银行
	BankAccount   string  `json:"bankAccount"`   // 银行卡号
	BankHolder    string  `json:"bankHolder"`    // 持卡人姓名
}

// Validate 验证提现请求参数
func (r UserWithdrawRequest) Validate() error {
	// 1. 验证提现金额
	if r.Amount <= 0 {
		return errors.New("提现金额必须大于0")
	}
	// 金额最多保留2位小数
	if !isValidAmount(r.Amount) {
		return errors.New("提现金额格式不正确")
	}

	// 2. 验证提现方式
	validTypes := map[string]bool{
		"alipay": true,
		"wechat": true,
		"bank":   true,
	}
	if !validTypes[r.WithdrawType] {
		return errors.New("无效的提现方式")
	}

	// 3. 根据提现方式验证对应的账户信息
	switch r.WithdrawType {
	case "alipay":
		return r.validateAlipay()
	case "wechat":
		return r.validateWechat()
	case "bank":
		return r.validateBank()
	}

	return nil
}

// validateAlipay 验证支付宝账户信息
func (r UserWithdrawRequest) validateAlipay() error {
	// 验证支付宝账号
	if strings.TrimSpace(r.AlipayAccount) == "" {
		return errors.New("请输入支付宝账号")
	}

	// 支付宝账号：手机号或邮箱
	if !isValidPhone(r.AlipayAccount) && !isValidEmail(r.AlipayAccount) {
		return errors.New("支付宝账号格式不正确")
	}

	// 验证姓名
	if strings.TrimSpace(r.AlipayName) == "" {
		return errors.New("请输入支付宝实名姓名")
	}
	if !isValidName(r.AlipayName) {
		return errors.New("姓名格式不正确")
	}

	return nil
}

// validateWechat 验证微信账户信息
func (r UserWithdrawRequest) validateWechat() error {
	// 验证微信账号
	if strings.TrimSpace(r.WechatAccount) == "" {
		return errors.New("请输入微信账号")
	}
	// 微信账号：微信号或手机号
	if len(r.WechatAccount) < 6 || len(r.WechatAccount) > 20 {
		return errors.New("微信账号格式不正确")
	}

	// 验证姓名
	if strings.TrimSpace(r.WechatName) == "" {
		return errors.New("请输入微信实名姓名")
	}
	if !isValidName(r.WechatName) {
		return errors.New("姓名格式不正确")
	}

	return nil
}

// validateBank 验证银行卡信息
func (r UserWithdrawRequest) validateBank() error {
	// 验证开户银行
	if strings.TrimSpace(r.BankName) == "" {
		return errors.New("请输入开户银行")
	}
	if len(r.BankName) < 2 || len(r.BankName) > 50 {
		return errors.New("开户银行名称长度不正确")
	}

	// 验证银行卡号
	if strings.TrimSpace(r.BankAccount) == "" {
		return errors.New("请输入银行卡号")
	}
	if !isValidBankCard(r.BankAccount) {
		return errors.New("银行卡号格式不正确")
	}

	// 验证持卡人姓名
	if strings.TrimSpace(r.BankHolder) == "" {
		return errors.New("请输入持卡人姓名")
	}
	if !isValidName(r.BankHolder) {
		return errors.New("持卡人姓名格式不正确")
	}

	return nil
}

// GetAccountInfo 获取账户信息（用于保存到数据库）
func (r UserWithdrawRequest) GetAccountInfo() (accountName, accountNo string) {
	switch r.WithdrawType {
	case "alipay":
		return r.AlipayName, r.AlipayAccount
	case "wechat":
		return r.WechatName, r.WechatAccount
	case "bank":
		return r.BankHolder, r.BankAccount
	}
	return "", ""
}

// GetFullAccountInfo 获取完整账户信息（用于展示）
func (r UserWithdrawRequest) GetFullAccountInfo() string {
	switch r.WithdrawType {
	case "alipay":
		return "支付宝 " + r.AlipayAccount + " (" + r.AlipayName + ")"
	case "wechat":
		return "微信 " + r.WechatAccount + " (" + r.WechatName + ")"
	case "bank":
		return r.BankName + " " + r.BankAccount + " (" + r.BankHolder + ")"
	}
	return ""
}

// ==================== 辅助验证函数 ====================

// isValidAmount 验证金额格式（最多2位小数）
func isValidAmount(amount float64) bool {
	// 转换为字符串检查小数位数
	str := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.10f", amount), "0"), ".")
	parts := strings.Split(str, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return false
	}
	return true
}

// isValidName 验证姓名（2-20个字符，支持中英文）
func isValidName(name string) bool {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 60 {
		return false
	}
	// 允许中文、英文字母、空格、点号
	matched, _ := regexp.MatchString(`^[\p{Han}a-zA-Z\s.·]+$`, name)
	return matched
}

// isValidBankCard 验证银行卡号（Luhn算法）
func isValidBankCard(cardNo string) bool {
	// 移除空格和横杠
	cardNo = strings.ReplaceAll(cardNo, " ", "")
	cardNo = strings.ReplaceAll(cardNo, "-", "")

	// 长度检查：银行卡号通常是16-19位
	if len(cardNo) < 16 || len(cardNo) > 19 {
		return false
	}

	// 必须全是数字
	matched, _ := regexp.MatchString(`^\d+$`, cardNo)
	if !matched {
		return false
	}

	// Luhn算法校验
	return luhnCheck(cardNo)
}

// luhnCheck Luhn算法（模10算法）
func luhnCheck(cardNo string) bool {
	sum := 0
	alternate := false

	// 从右到左遍历
	for i := len(cardNo) - 1; i >= 0; i-- {
		n := int(cardNo[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n = n%10 + n/10
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

type WithdrawRecordRequest struct {
	PageInfo
	Status string `json:"status" form:"status"` // 状态筛选：all/pending/success/rejected
	Time   string `json:"time" form:"time"`     // 时间筛选：all/week/month/quarter
}

// Validate 验证请求参数
func (r WithdrawRecordRequest) Validate() error {
	// 1. 验证分页参数
	if r.Page < 1 {
		return errors.New("页码必须大于0")
	}
	if r.PageSize < 1 {
		return errors.New("每页数量必须大于0")
	}
	if r.PageSize > 100 {
		return errors.New("每页数量不能超过100")
	}

	// 2. 验证状态参数（可选）
	if r.Status != "" {
		validStatus := map[string]bool{
			"all":      true,
			"pending":  true,
			"success":  true,
			"rejected": true,
		}
		if !validStatus[r.Status] {
			return errors.New("无效的状态参数")
		}
	}

	// 3. 验证时间参数（可选）
	if r.Time != "" {
		validTime := map[string]bool{
			"all":     true,
			"week":    true,
			"month":   true,
			"quarter": true,
		}
		if !validTime[r.Time] {
			return errors.New("无效的时间参数")
		}
	}

	return nil
}

// GetStatusCondition 获取状态查询条件
func (r WithdrawRecordRequest) GetStatusCondition() []string {
	switch r.Status {
	case "pending":
		return []string{"pending"}
	case "success":
		// 成功包含：已通过和已完成
		return []string{"approved", "completed"}
	case "rejected":
		return []string{"rejected"}
	default:
		// "all" 或空值，返回 nil 表示不筛选
		return nil
	}
}

// GetTimeRange 获取时间范围
func (r WithdrawRecordRequest) GetTimeRange() (startTime, endTime *time.Time) {
	now := time.Now()

	switch r.Time {
	case "week":
		// 本周：从本周一 00:00:00 到当前时间
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 周日为7
		}
		start := now.AddDate(0, 0, -(weekday - 1)).
			Truncate(24 * time.Hour)
		startTime = &start

	case "month":
		// 本月：从本月1号 00:00:00 到当前时间
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startTime = &start

	case "quarter":
		// 近三个月：从3个月前的今天 00:00:00 到当前时间
		start := now.AddDate(0, -3, 0).Truncate(24 * time.Hour)
		startTime = &start

	default:
		// "all" 或空值，不限制时间
		return nil, nil
	}

	return startTime, nil
}
