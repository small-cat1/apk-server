package request

type UserWithdrawRequest struct {
	Amount        string `json:"amount"`
	WithdrawType  string `json:"withdrawType"`
	AlipayAccount string `json:"alipayAccount"`
	AlipayName    string `json:"alipayName"`
	WechatAccount string `json:"wechatAccount"`
	WechatName    string `json:"wechatName"`
	BankName      string `json:"bankName"`
	BankAccount   string `json:"bankAccount"`
	BankHolder    string `json:"bankHolder"`
}

func (r UserWithdrawRequest) Validate() error {

	return nil
}
