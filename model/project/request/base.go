package request

import (
	"errors"
	"fmt"
)

// GetById Find by id structure
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type DeleteIDs struct {
	IDs []uint `json:"ids"` //删除的IDS
}

// Validate 验证删除ID请求
func (req *DeleteIDs) Validate() error {
	// 验证IDs不能为空
	if len(req.IDs) == 0 {
		return errors.New("删除的ID列表不能为空")
	}

	// 验证ID数量限制（防止一次性删除过多）
	if len(req.IDs) > 100 {
		return errors.New("单次删除数量不能超过100条")
	}

	// 验证每个ID的有效性
	for i, id := range req.IDs {
		if id == 0 {
			return fmt.Errorf("第%d个ID无效：ID不能为0", i+1)
		}
	}

	// 验证是否有重复ID
	if err := validateNoDuplicateIDs(req.IDs); err != nil {
		return err
	}

	return nil
}

// validateNoDuplicateIDs 验证没有重复的ID
func validateNoDuplicateIDs(ids []uint) error {
	seen := make(map[uint]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			return fmt.Errorf("ID列表中存在重复的ID: %d", id)
		}
		seen[id] = true
	}
	return nil
}
