package request

import (
	"ApkAdmin/utils"
)

type BaseLoginRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captchaKey" binding:"required"`
}

func (r BaseLoginRequest) Validate() error {
	// 验证手机号
	if err := utils.ValidatePhone(r.Phone); err != nil {
		return err
	}
	// 验证密码
	if err := utils.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

type BaseRegisterRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captchaKey" binding:"required"`
	InviteCode string `json:"inviteCode"`
}

func (r BaseRegisterRequest) Validate() error {
	// 验证手机号
	if err := utils.ValidatePhone(r.Phone); err != nil {
		return err
	}
	// 验证邮箱
	if err := utils.ValidateEmail(r.Email); err != nil {
		return err
	}
	// 验证密码
	if err := utils.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}
