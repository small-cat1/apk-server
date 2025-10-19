package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type PaymentProviderService struct{}

// GetPaymentProvider 获取单个支付服务商详情
func (service *PaymentProviderService) GetPaymentProvider(id uint) (provider project.PaymentProvider, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider, errors.New("支付服务商不存在")
		}
	}
	return
}

func (service *PaymentProviderService) GetAllPaymentProviders() (list interface{}, err error) {
	var paymentProviders []project.PaymentProvider
	err = global.GVA_DB.Model(&project.PaymentProvider{}).
		Where("status = ?", "active").
		Order("sort_order desc").
		Find(&paymentProviders).Error
	if err != nil {
		return nil, err
	}
	if len(paymentProviders) <= 0 {
		return nil, errors.New("暂无可以用的支付服务商")
	}
	var result []response.PaymentProviderResp
	for _, v := range paymentProviders {
		result = append(result, response.PaymentProviderResp{
			Id:    v.ID,
			Name:  v.Name,
			Value: v.Code,
			Icon:  v.Icon,
		})
	}
	return result, nil
}

// GetPaymentProviderList 获取支付服务商列表
func (service *PaymentProviderService) GetPaymentProviderList(info request.PaymentProviderListRequest, orderKey string, desc bool) (list []project.PaymentProvider, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&project.PaymentProvider{})

	// 构建查询条件
	if info.Code != "" {
		db = db.Where("code LIKE ?", "%"+info.Code+"%")
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.CreatedAtGte != "" {
		db = db.Where("created_at >= ?", info.CreatedAtGte)
	}
	if info.CreatedAtLte != "" {
		db = db.Where("created_at <= ?", info.CreatedAtLte)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 排序
	if orderKey != "" {
		orderStr := orderKey
		if desc {
			orderStr += " desc"
		}
		db = db.Order(orderStr)
	} else {
		db = db.Order("sort_order asc, created_at desc")
	}

	// 分页查询
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

// CreatePaymentProvider 创建支付服务商
func (service *PaymentProviderService) CreatePaymentProvider(req request.PaymentProviderCreateRequest) (err error) {
	// 检查代码是否已存在
	var count int64
	err = global.GVA_DB.Model(&project.PaymentProvider{}).Where("code = ?", req.Code).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("服务商代码已存在")
	}

	// 创建支付服务商
	provider := project.PaymentProvider{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
	}

	err = global.GVA_DB.Create(&provider).Error
	return err
}

// UpdatePaymentProvider 更新支付服务商
func (service *PaymentProviderService) UpdatePaymentProvider(req *request.PaymentProviderUpdateRequest) (err error) {
	// 检查记录是否存在
	var provider project.PaymentProvider
	err = global.GVA_DB.Where("id = ?", req.ID).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("支付服务商不存在")
		}
		return err
	}

	// 检查代码是否与其他记录重复
	var count int64
	err = global.GVA_DB.Model(&project.PaymentProvider{}).Where("code = ? AND id != ?", req.Code, req.ID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("服务商代码已存在")
	}

	// 更新数据
	updates := map[string]interface{}{
		"code":          req.Code,
		"name":          req.Name,
		"description":   req.Description,
		"status":        req.Status,
		"icon":          req.Icon,
		"config_schema": req.ConfigSchema,
		"sort_order":    req.SortOrder,
	}

	err = global.GVA_DB.Model(&provider).Updates(updates).Error
	return err
}

// DeletePaymentProvider 删除支付服务商
func (service *PaymentProviderService) DeletePaymentProvider(id uint) (err error) {
	// 检查是否存在
	var provider project.PaymentProvider
	err = global.GVA_DB.Where("id = ?", id).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("支付服务商不存在")
		}
		return err
	}
	// 这里可以添加业务逻辑检查，比如是否有关联的订单等
	// TODO: 检查是否有关联数据
	var total int64
	err = global.GVA_DB.Model(&project.PaymentAccount{}).Where("provider_code = ?", provider.Code).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该服务商下存在账户，无法删除")
	}
	err = global.GVA_DB.Delete(&provider).Error
	return err
}

// BatchUpdatePaymentProviderStatus 批量更新支付服务商状态
func (service *PaymentProviderService) BatchUpdatePaymentProviderStatus(ids []uint, status string) (err error) {
	if len(ids) == 0 {
		return errors.New("请选择要更新的数据")
	}

	if status != "active" && status != "inactive" {
		return errors.New("无效的状态值")
	}

	err = global.GVA_DB.Model(&project.PaymentProvider{}).Where("id IN ?", ids).Update("status", status).Error
	return err
}

// CheckProviderCodeAvailable 检查服务商代码是否可用
func (service *PaymentProviderService) CheckProviderCodeAvailable(code string, excludeID uint) (available bool, err error) {
	if strings.TrimSpace(code) == "" {
		return false, errors.New("代码不能为空")
	}

	var count int64
	db := global.GVA_DB.Model(&project.PaymentProvider{}).Where("code = ?", code)

	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}

	err = db.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// UpdateProviderSortOrder 更新服务商排序
func (service *PaymentProviderService) UpdateProviderSortOrder(id uint, sortOrder *int) (err error) {
	// 检查记录是否存在
	var provider project.PaymentProvider
	err = global.GVA_DB.Where("id = ?", id).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("支付服务商不存在")
		}
		return err
	}

	err = global.GVA_DB.Model(&provider).Update("sort_order", sortOrder).Error
	return err
}
