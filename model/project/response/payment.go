package response

type PaymentProviderResp struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
}
