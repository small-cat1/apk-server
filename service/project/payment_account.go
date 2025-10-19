package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// 不同支付方式的配置结构体
type WechatPayConfig struct {
	AppID       string `json:"app_id" validate:"required"`
	MchID       string `json:"mch_id" validate:"required"`
	ApiKey      string `json:"api_key" validate:"required"`
	CertPath    string `json:"cert_path"`
	KeyPath     string `json:"key_path"`
	NotifyURL   string `json:"notify_url"`
	ReturnURL   string `json:"return_url"`
	Environment string `json:"environment"` // sandbox, production
}

type AlipayConfig struct {
	AppID       string `json:"app_id" validate:"required"`
	PrivateKey  string `json:"private_key" validate:"required"`
	PublicKey   string `json:"public_key" validate:"required"`
	Gateway     string `json:"gateway"`
	NotifyURL   string `json:"notify_url"`
	ReturnURL   string `json:"return_url"`
	Environment string `json:"environment"` // sandbox, production
}

type StripeConfig struct {
	PublishableKey string `json:"publishable_key" validate:"required"`
	SecretKey      string `json:"secret_key" validate:"required"`
	WebhookSecret  string `json:"webhook_secret"`
	Currency       string `json:"currency"`
	Environment    string `json:"environment"`
}

type PayPalConfig struct {
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	WebhookID    string `json:"webhook_id"`
	Environment  string `json:"environment"` // sandbox, live
}

// 配置接口 - 所有支付配置都需要实现
type PaymentConfigInterface interface {
	Validate() error
	GetNotifyURL() string
	GetReturnURL() string
	IsProduction() bool
}

// 实现配置接口
func (c *WechatPayConfig) Validate() error {
	if c.AppID == "" || c.MchID == "" || c.ApiKey == "" {
		return errors.New("微信支付必填参数不能为空")
	}
	return nil
}

func (c *WechatPayConfig) GetNotifyURL() string {
	return c.NotifyURL
}

func (c *WechatPayConfig) GetReturnURL() string {
	return c.ReturnURL
}

func (c *WechatPayConfig) IsProduction() bool {
	return c.Environment == "production"
}

func (c *AlipayConfig) Validate() error {
	if c.AppID == "" || c.PrivateKey == "" || c.PublicKey == "" {
		return errors.New("支付宝必填参数不能为空")
	}
	return nil
}

func (c *AlipayConfig) GetNotifyURL() string {
	return c.NotifyURL
}

func (c *AlipayConfig) GetReturnURL() string {
	return c.ReturnURL
}

func (c *AlipayConfig) IsProduction() bool {
	return c.Environment == "production"
}

// PaymentAccountService 支付账号服务
type PaymentAccountService struct{}

// GetPaymentAccountInfoList  服务方法
func (p *PaymentAccountService) GetPaymentAccountInfoList(req request.PaymentAccountSearch) (list []project.PaymentAccount, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.GVA_DB.Model(&project.PaymentAccount{}).Preload("Provider")

	// 添加搜索条件
	if req.ProviderCode != "" {
		db = db.Where("provider_code = ?", req.ProviderCode)
	}

	if req.AccountType != "" {
		db = db.Where("account_type = ?", req.AccountType)
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	if req.Group != "" {
		db = db.Where("group = ?", req.Group)
	}

	if req.Region != "" {
		db = db.Where("region = ?", req.Region)
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取数据
	err = db.Limit(limit).Offset(offset).Order("weight DESC, created_at DESC").Find(&list).Error
	return
}

// CreatePaymentAccount 创建支付账号
func (p *PaymentAccountService) CreatePaymentAccount(req *request.CreatePaymentAccountReq) error {
	// 检查支付服务商是否存在
	var provider project.PaymentProvider
	if err := global.GVA_DB.Where("code = ? AND deleted_at IS NULL", req.ProviderCode).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("支付服务商不存在")
		}
		return err
	}

	// 检查账号名称是否重复
	var count int64
	if err := global.GVA_DB.Model(&project.PaymentAccount{}).
		Where("name = ? AND deleted_at IS NULL", req.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("账号名称已存在")
	}

	// 验证配置参数
	err := p.validateAccountConfig(req.ProviderCode, req.Config)
	if err != nil {
		return fmt.Errorf("配置参数验证失败: %v", err)
	}

	// 转换配置为JSON字符串
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return fmt.Errorf("配置参数序列化失败: %v", err)
	}

	// 设置默认值
	if req.Status == "" {
		req.Status = "active"
	}
	if req.Weight == 0 {
		req.Weight = 1
	}

	account := project.PaymentAccount{
		Name:           req.Name,
		ProviderCode:   req.ProviderCode,
		AccountType:    req.AccountType,
		Config:         string(configBytes),
		Status:         req.Status,
		Weight:         req.Weight,
		MaxDailyAmount: req.MaxDailyAmount,
		DailyAmount:    0,
		TotalAmount:    0,
		TotalOrders:    0,
		Group:          req.Group,
		Region:         req.Region,
		Tags:           req.Tags,
		Remark:         req.Remark,
	}

	return global.GVA_DB.Create(&account).Error
}

// UpdatePaymentAccount 更新支付账号
func (p *PaymentAccountService) UpdatePaymentAccount(req *request.PaymentAccountUpdate) error {
	// 检查账号是否存在
	var account project.PaymentAccount
	if err := global.GVA_DB.Where("id = ? AND deleted_at IS NULL", req.ID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("账号不存在")
		}
		return err
	}

	// 检查名称重复
	var count int64
	if err := global.GVA_DB.Model(&project.PaymentAccount{}).
		Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, req.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("账号名称已存在")
	}

	// 验证配置参数
	var provider project.PaymentProvider
	if err := global.GVA_DB.Where("code = ? AND deleted_at IS NULL", account.ProviderCode).First(&provider).Error; err == nil {
		if err := p.validateAccountConfig(account.ProviderCode, req.Config); err != nil {
			return fmt.Errorf("配置参数验证失败: %v", err)
		}
	}

	// 转换配置为JSON字符串
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return fmt.Errorf("配置参数序列化失败: %v", err)
	}

	// 更新账号信息
	updates := map[string]interface{}{
		"name":             req.Name,
		"config":           string(configBytes),
		"weight":           req.Weight,
		"max_daily_amount": req.MaxDailyAmount,
		"group":            req.Group,
		"region":           req.Region,
		"tags":             req.Tags,
		"remark":           req.Remark,
		"updated_at":       time.Now(),
	}

	if req.Status != "" {
		updates["status"] = req.Status
	}

	return global.GVA_DB.Model(&account).Updates(updates).Error
}

// GetPaymentAccount 获取账号信息
func (p *PaymentAccountService) GetPaymentAccount(accountID uint) (interface{}, error) {
	var account project.PaymentAccount
	err := global.GVA_DB.Preload("Provider").Where("id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return account, nil
}

// UpdateAccountWeight 更新账号权重
func (p *PaymentAccountService) UpdateAccountWeight(id uint, weight int) error {
	return global.GVA_DB.Model(&project.PaymentAccount{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("weight", weight).Error
}

// ResetDailyAmount 重置日交易金额
func (p *PaymentAccountService) ResetDailyAmount(id uint) error {
	return global.GVA_DB.Model(&project.PaymentAccount{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("daily_amount", 0).Error
}

// GetAccountConfig 获取解析后的配置
func (p *PaymentAccountService) GetAccountConfig(accountID uint) (interface{}, error) {
	var account project.PaymentAccount
	err := global.GVA_DB.Preload("Provider").Where("id = ?", accountID).First(&account).Error
	if err != nil {
		return nil, err
	}

	return p.parseConfig(account.ProviderCode, account.Config)
}

// parseConfig 根据支付方式解析配置
func (p *PaymentAccountService) parseConfig(providerCode, configJSON string) (interface{}, error) {
	switch providerCode {
	case "wechat":
		var config WechatPayConfig
		err := json.Unmarshal([]byte(configJSON), &config)
		return &config, err

	case "alipay":
		var config AlipayConfig
		err := json.Unmarshal([]byte(configJSON), &config)
		return &config, err

	case "stripe":
		var config StripeConfig
		err := json.Unmarshal([]byte(configJSON), &config)
		return &config, err

	case "paypal":
		var config PayPalConfig
		err := json.Unmarshal([]byte(configJSON), &config)
		return &config, err

	default:
		// 对于未知的支付方式，返回通用map
		var config map[string]interface{}
		err := json.Unmarshal([]byte(configJSON), &config)
		return config, err
	}
}

// ValidateAccountConfig 验证账号配置
func (p *PaymentAccountService) validateAccountConfig(providerCode string, config map[string]interface{}) error {
	configJSON, _ := json.Marshal(config)
	parsedConfig, err := p.parseConfig(providerCode, string(configJSON))
	if err != nil {
		return err
	}
	// 如果配置实现了Validate接口，则调用验证
	if validator, ok := parsedConfig.(PaymentConfigInterface); ok {
		return validator.Validate()
	}
	return nil
}

// SelectBestAccount 选择最佳支付账号
func (p *PaymentAccountService) SelectBestAccount(providerCode string, amount float64, region string) (*project.PaymentAccount, error) {
	db := global.GVA_DB.Where("provider_code = ? AND status = ?", providerCode, "active")

	// 如果指定了地区，优先选择对应地区的账号
	if region != "" {
		db = db.Where("region = ? OR region = ''", region)
	}

	var accounts []project.PaymentAccount
	err := db.Order("weight DESC, created_at DESC").Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, errors.New("没有可用的支付账号")
	}

	// 过滤超过日限额的账号
	for _, account := range accounts {
		if account.MaxDailyAmount <= 0 || account.DailyAmount+amount <= account.MaxDailyAmount {
			return &account, nil
		}
	}

	return nil, errors.New("所有账号都已达到日限额")
}

func (p *PaymentAccountService) DeletePaymentAccount(id uint) error {
	return global.GVA_DB.Delete(&project.PaymentAccount{}, id).Error
}
