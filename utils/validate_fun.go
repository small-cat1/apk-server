package utils

import (
	"errors"
	"regexp"
)

// ValidatePhone 验证手机号（支持中国大陆手机号）
func ValidatePhone(phone string) error {
	if phone == "" {
		return errors.New("手机号不能为空")
	}

	// 中国大陆手机号正则：1开头，第二位是3-9，后面9位数字
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("手机号格式不正确")
	}

	return nil
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}

	// 基本的邮箱格式验证
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// ValidatePassword 验证密码
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	// 密码长度至少8位
	if len(password) < 8 {
		return errors.New("密码长度至少为8位")
	}

	// 密码最大长度限制
	if len(password) > 20 {
		return errors.New("密码长度不能超过20位")
	}

	// 检查是否包含至少一个字母
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	// 检查是否包含至少一个数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含字母和数字")
	}

	return nil
}
