package request

type MembershipPlanOrderRequest struct {
	PackageId     int    `json:"packageId" binding:"required"`     // 套餐ID
	PaymentMethod string `json:"paymentMethod" binding:"required"` //支付方式
}

type AccountOrderRequest struct {
	AppId         int    `json:"appId" binding:"required"`         //应用ID
	Quantity      int    `json:"quantity" binding:"required"`      //数量
	Amount        string `json:"amount" binding:"required"`        //金额
	PaymentMethod string `json:"paymentMethod" binding:"required"` //支付方式
}
