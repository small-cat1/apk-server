package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"ApkAdmin/utils/crypto"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

type AppAccountService struct{}

// CreateAccount 创建账号
func (s *AppAccountService) CreateAccount(req request.CreateAppAccountRequest, userId uint) error {
	var app project.Application
	err := global.GVA_DB.Model(&project.Application{}).Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		global.GVA_LOG.Error("创建应用账号失败，查找应用错误", zap.Error(err))
		return err
	}
	ExtraInfo, err := json.Marshal(req.ExtraInfo)
	if err != nil {
		global.GVA_LOG.Error("创建应用账号失败，序列化ExtraInfo参数错误", zap.Error(err))
		return err
	}
	// 验证状态
	account := project.AppAccount{
		AppID:         req.AppID,
		AccountDetail: req.AccountDetail,
		ExtraInfo:     string(ExtraInfo),
		CategoryID:    *app.CategoryID,
		AccountStatus: req.AccountStatus,
		CreatedBy:     userId,
	}
	// 保存到数据库
	return global.GVA_DB.Create(&account).Error
}

func (s *AppAccountService) UpdateAccount(req request.UpdateAppAccountRequest, userID uint) error {
	var app project.Application
	err := global.GVA_DB.Model(&project.Application{}).Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		global.GVA_LOG.Error("创建应用账号失败，查找应用错误", zap.Error(err))
		return err
	}
	var account project.AppAccount
	err = global.GVA_DB.Model(&project.AppAccount{}).
		Where("id = ? and app_id = ?", req.ID, req.AppID).First(&account).Error
	if err != nil {
		global.GVA_LOG.Error("创建应用账号失败，查找应用账号错误", zap.Error(err))
		return err
	}
	encrypted, err := crypto.EncryptAccountDetail(req.AccountDetail)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"category_id":    app.CategoryID,
		"account_detail": encrypted,
		"extra_info":     req.ExtraInfo,
		"account_status": req.AccountStatus,
		"updated_by":     userID,
	}
	return global.GVA_DB.Model(&project.AppAccount{}).
		Where("id = ?", req.ID).
		Updates(updates).Error
}

// UpdateAccountStatus 更新账号状态
func (s *AppAccountService) UpdateAccountStatus(req request.UpdateAccountStatusRequest) (*response.BatchUpdateStatusResponse, error) {
	if len(req.IDs) == 0 {
		return nil, fmt.Errorf("账号ID列表不能为空")
	}
	var accounts []project.AppAccount
	err := global.GVA_DB.Model(&project.AppAccount{}).
		Where("id in ?", req).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	result := &response.BatchUpdateStatusResponse{
		FailDetails: make([]response.FailDetail, 0),
	}
	// 验证每个账号的状态转换是否合法
	var validAccountIDs []uint
	for _, account := range accounts {
		currentStatus := account.AccountStatus
		// 验证状态转换
		if err := currentStatus.ValidateTransition(req.Status); err != nil {
			result.FailCount++
			result.FailDetails = append(result.FailDetails, response.FailDetail{
				AccountID: account.ID,
				Reason:    err.Error(),
			})
			continue
		}
		validAccountIDs = append(validAccountIDs, account.ID)
	}
	// 批量更新合法的账号状态
	if len(validAccountIDs) > 0 {
		err = global.GVA_DB.Model(&project.AppAccount{}).
			Where("id in ?", validAccountIDs).
			Update("account_status", req.Status).Error
		if err != nil {
			return nil, fmt.Errorf("批量更新状态失败: %w", err)
		}
		result.SuccessCount = len(validAccountIDs)
	}
	return result, nil
}

// GetAccountList 获取账号列表
func (s *AppAccountService) GetAccountList(req request.GetAppAccountListRequest) ([]response.AppAccountResponse, int64, error) {
	var accounts []project.AppAccount
	var total int64
	db := global.GVA_DB.Model(&project.AppAccount{})
	// 条件筛选
	if req.AppID != "" {
		db = db.Where("app_id = ?", req.AppID)
	}
	if req.CategoryID > 0 {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.AccountNo != "" {
		db = db.Where("account_no = ?", req.AccountNo)
	}
	if req.AccountStatus > 0 {
		db = db.Where("account_status = ?", req.AccountStatus)
	}

	// 统计总数
	db.Count(&total)

	// 分页查询
	err := db.
		Preload("Application").
		Preload("Creator").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("created_at DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var list []response.AppAccountResponse
	for _, account := range accounts {
		list = append(list, response.AppAccountResponse{
			ID:                account.ID,
			AppID:             account.AppID,
			AppName:           account.Application.AppName,
			CategoryID:        account.CategoryID,
			AccountNo:         account.AccountNo,
			AccountStatus:     account.AccountStatus.Int(),
			AccountStatusText: account.AccountStatus.GetAccountStatusText(),
			CreatedAt:         account.CreatedAt,
			UpdatedAt:         account.UpdatedAt,
			CreatedBy:         account.CreatedBy,
			CreatorName:       account.Creator.Username,
		})
	}

	return list, total, nil
}

// GetAccountDetail 获取账号详情（已解密）
func (s *AppAccountService) GetAccountDetail(id uint) (*response.AppAccountResponse, error) {
	var account project.AppAccount
	// AfterFind 会自动解密
	if err := global.GVA_DB.Preload("Application").Preload("Creator").First(&account, id).Error; err != nil {
		return nil, err
	}
	return &response.AppAccountResponse{
		ID:            account.ID,
		AppID:         account.AppID,
		AppName:       account.Application.AppName,
		AccountNo:     account.AccountNo,
		AccountDetail: account.AccountDetail, // 已解密
		AccountStatus: account.AccountStatus.Int(),
		ExtraInfo:     account.ExtraInfo,
		CreatedAt:     account.CreatedAt,
		CreatorName:   account.Creator.Username,
	}, nil
}

// GetAppAccountOrderDetail 查看售出账号的订单详情
func (s AppAccountService) GetAppAccountOrderDetail(accountID uint) (*response.AppAccountOrderDetailResp, error) {
	var account project.AppAccount
	if err := global.GVA_DB.Where("account_status = ?", constants.AppAccountStatusSold).First(&account, accountID).Error; err != nil {
		return nil, err
	}
	return nil, nil
}
