package request

import "errors"

// CreateCommissionTierRequest 创建阶梯等级请求
type CreateCommissionTierRequest struct {
	Name            string  `json:"name" binding:"required"` // 等级名称
	MinSubordinates int     `json:"minSubordinates" `        // 最低直属下级人数
	Rate            float64 `json:"rate" binding:"required"` // 分佣比例(%)
	Color           string  `json:"color"`                   // 等级颜色
	Icon            string  `json:"icon"`                    // 等级图标
	Sort            int     `json:"sort"`                    // 排序
	Status          int     `json:"status"`                  // 状态：1-启用, 0-禁用
}

// Validate 验证创建请求
func (r *CreateCommissionTierRequest) Validate() error {
	if len(r.Name) < 2 || len(r.Name) > 50 {
		return errors.New("等级名称长度必须在2-50个字符之间")
	}
	if r.MinSubordinates < 0 {
		return errors.New("最低下级人数不能为负数")
	}
	if r.Rate <= 0 || r.Rate > 100 {
		return errors.New("分佣比例必须在0-100之间")
	}
	if r.Sort < 0 {
		return errors.New("排序值不能为负数")
	}
	if r.Status != 0 && r.Status != 1 {
		return errors.New("状态值必须为0或1")
	}
	return nil
}

// UpdateCommissionTierRequest 更新阶梯等级请求
type UpdateCommissionTierRequest struct {
	ID              int     `json:"id" binding:"required"`   // 等级ID
	Name            string  `json:"name" binding:"required"` // 等级名称
	MinSubordinates int     `json:"minSubordinates" `        // 最低直属下级人数
	Rate            float64 `json:"rate" binding:"required"` // 分佣比例(%)
	Color           string  `json:"color"`                   // 等级颜色
	Icon            string  `json:"icon"`                    // 等级图标
	Sort            int     `json:"sort"`                    // 排序
	Status          int     `json:"status"`                  // 状态：1-启用, 0-禁用
}

// Validate 验证更新请求
func (r *UpdateCommissionTierRequest) Validate() error {
	if r.ID <= 0 {
		return errors.New("等级ID无效")
	}
	if len(r.Name) < 2 || len(r.Name) > 50 {
		return errors.New("等级名称长度必须在2-50个字符之间")
	}
	if r.MinSubordinates < 0 {
		return errors.New("最低下级人数不能为负数")
	}
	if r.Rate <= 0 || r.Rate > 100 {
		return errors.New("分佣比例必须在0-100之间")
	}
	if r.Sort < 0 {
		return errors.New("排序值不能为负数")
	}
	if r.Status != 0 && r.Status != 1 {
		return errors.New("状态值必须为0或1")
	}
	return nil
}

// DeleteCommissionTiersRequest 删除阶梯等级请求（支持批量）
type DeleteCommissionTiersRequest struct {
	IDs []int `json:"ids" binding:"required,min=1"` // 等级ID列表
}

// Validate 验证删除请求
func (r *DeleteCommissionTiersRequest) Validate() error {
	if len(r.IDs) == 0 {
		return errors.New("请选择要删除的等级")
	}
	for _, id := range r.IDs {
		if id <= 0 {
			return errors.New("等级ID无效")
		}
	}
	return nil
}

// GetCommissionTierRequest 获取单个阶梯等级请求
type GetCommissionTierRequest struct {
	ID int `json:"id" form:"id" binding:"required"` // 等级ID
}

// Validate 验证获取请求
func (r *GetCommissionTierRequest) Validate() error {
	if r.ID <= 0 {
		return errors.New("等级ID无效")
	}
	return nil
}

// GetCommissionTierListRequest 获取阶梯等级列表请求
type GetCommissionTierListRequest struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页数量
	Name     string `json:"name" form:"name"`         // 等级名称（模糊搜索）
	Status   *int8  `json:"status" form:"status"`     // 状态：1-启用, 0-禁用
}

// Validate 验证列表查询请求
func (r *GetCommissionTierListRequest) Validate() error {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 10
	}
	if r.Status != nil && *r.Status != 0 && *r.Status != 1 {
		return errors.New("状态值必须为0或1")
	}
	return nil
}

// UpdateCommissionTierStatusRequest 更新等级状态请求
type UpdateCommissionTierStatusRequest struct {
	ID     int  `json:"id" binding:"required"`     // 等级ID
	Status *int `json:"status" binding:"required"` // 状态：1-启用, 0-禁用
}

// Validate 验证状态更新请求
func (r *UpdateCommissionTierStatusRequest) Validate() error {
	if r.ID <= 0 {
		return errors.New("等级ID无效")
	}
	if *r.Status != 0 && *r.Status != 1 {
		return errors.New("状态值必须为0或1")
	}
	return nil
}

// UpdateCommissionTierSortRequest 批量更新等级排序请求
type UpdateCommissionTierSortRequest struct {
	Sorts []map[string]interface{} `json:"sorts" binding:"required,min=1"` // 排序数据列表
}

// Validate 验证排序更新请求
func (r *UpdateCommissionTierSortRequest) Validate() error {
	if len(r.Sorts) == 0 {
		return errors.New("排序数据不能为空")
	}
	for _, item := range r.Sorts {
		id, ok1 := item["id"]
		sort, ok2 := item["sort"]
		if !ok1 || !ok2 {
			return errors.New("排序数据格式错误，需要包含id和sort字段")
		}
		// 验证类型
		switch id.(type) {
		case float64, int, int64:
			// 合法类型
		default:
			return errors.New("id必须为数字类型")
		}
		switch sort.(type) {
		case float64, int, int64:
			// 合法类型
		default:
			return errors.New("sort必须为数字类型")
		}
	}
	return nil
}
