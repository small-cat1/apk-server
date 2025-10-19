package project

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// ==================== 条件构建器 ====================

// WithPhone 手机号条件（非空才添加）
func WithPhone(phone string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if phone != "" {
			return db.Where("phone = ?", phone)
		}
		return db
	}
}

func WithUuid(uuid uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if &uuid != nil {
			return db.Where("uuid = ?", uuid)
		}
		return db
	}
}

// WithEmail 邮箱条件（非空才添加）
func WithEmail(email string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if email != "" {
			return db.Where("email = ?", email)
		}
		return db
	}
}

// WithPhoneOrEmail 手机号或邮箱条件（至少有一个非空才添加）
func WithPhoneOrEmail(phone, email string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if phone == "" && email == "" {
			return db
		}

		if phone != "" && email != "" {
			return db.Where("phone = ? OR email = ?", phone, email)
		}

		if phone != "" {
			return db.Where("phone = ?", phone)
		}

		return db.Where("email = ?", email)
	}
}

// WithStatus 状态条件（状态值有效才添加）
func WithStatus(status int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 如果你的状态值包含 0，这里需要调整判断逻辑
		// 方案1: 如果 0 是有效状态，去掉这个判断
		// 方案2: 使用指针类型 *int，判断 status != nil
		return db.Where("status = ?", status)
	}
}

// WithStatusPtr 状态条件（使用指针，nil 表示不过滤）
func WithStatusPtr(status *int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != nil {
			return db.Where("status = ?", *status)
		}
		return db
	}
}

// WithDeleted 是否包含已删除（这个不需要判断）
func WithDeleted() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}
}

// WithInviteCode 邀请码条件（非空才添加）
func WithInviteCode(inviteCode string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inviteCode != "" {
			return db.Where("referral_code = ?", inviteCode)
		}
		return db
	}
}

// WithID ID条件（大于0才添加）
func WithID(id uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id > 0 {
			return db.Where("id = ?", id)
		}
		return db
	}
}

// ==================== 更复杂的查询 ====================

// WithCustomCondition 自定义条件（value 非 nil 才添加）
func WithCustomCondition(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || value == nil {
			return db
		}

		// 如果是字符串，还要判断是否为空
		if str, ok := value.(string); ok && str == "" {
			return db
		}

		return db.Where(field+" = ?", value)
	}
}

// WithLike 模糊查询（非空才添加）
func WithLike(field, value string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || value == "" {
			return db
		}
		return db.Where(field+" LIKE ?", "%"+value+"%")
	}
}

// WithTimeRange 时间范围（时间有效才添加）
func WithTimeRange(field string, start, end time.Time) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || start.IsZero() || end.IsZero() {
			return db
		}
		return db.Where(field+" BETWEEN ? AND ?", start, end)
	}
}

// WithIn IN 查询（切片非空才添加）
func WithIn(field string, values []interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" || len(values) == 0 {
			return db
		}
		return db.Where(field+" IN ?", values)
	}
}
