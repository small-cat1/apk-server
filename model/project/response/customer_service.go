package response

type CustomerServiceContacts struct {
	Qq           string `json:"qq"`
	Wechat       string `json:"wechat"`
	WechatQrcode string `json:"wechatQrcode"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Im           bool   `json:"im"`
}

type ImConfig struct {
	Token string `json:"token"`
}

type CustomerServicePosition struct {
	Right  string `json:"right"`
	Bottom string `json:"bottom"`
	ZIndex int    `json:"zIndex"`
}

type CustomerServiceConfigResp struct {
	Enabled    bool                    `json:"enabled"`
	ShowText   bool                    `json:"showText"`
	ButtonText string                  `json:"buttonText"`
	Tooltip    string                  `json:"tooltip"`
	Position   CustomerServicePosition `json:"position"`
	Contacts   CustomerServiceContacts `json:"contacts"`
	WorkTime   interface{}             `json:"workTime"`
	Notice     interface{}             `json:"notice"`
	ImType     string                  `json:"imType"`
	ImConfig   ImConfig                `json:"imConfig"`
	Preload    bool                    `json:"preload"`
}
