package task

import (
	"ApkAdmin/model/common"
	"ApkAdmin/model/project"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: ClearTable
//@description: 清理数据库表数据
//@param: db(数据库对象) *gorm.DB, tableName(表名) string, compareField(比较字段) string, interval(间隔) string
//@return: error

func ClearTable(db *gorm.DB) error {
	var ClearTableDetail []common.ClearDB

	ClearTableDetail = append(ClearTableDetail, common.ClearDB{
		TableName:    "sys_operation_records",
		CompareField: "created_at",
		Interval:     "2160h",
	})

	ClearTableDetail = append(ClearTableDetail, common.ClearDB{
		TableName:    "jwt_blacklists",
		CompareField: "created_at",
		Interval:     "168h",
	})

	if db == nil {
		return errors.New("db Cannot be empty")
	}
	for _, detail := range ClearTableDetail {
		duration, err := time.ParseDuration(detail.Interval)
		if err != nil {
			return err
		}
		if duration < 0 {
			return errors.New("parse duration < 0")
		}
		err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s < ?", detail.TableName, detail.CompareField), time.Now().Add(-duration)).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// 每天或每小时执行
func UpdateActiveMembers(db *gorm.DB) {
	// 统计每个用户的活跃下级数（近30天有订单）
	sql := `
        UPDATE team_statistics ts
        SET active_members = (
            SELECT COUNT(DISTINCT u.id)
            FROM users u
            INNER JOIN user_statistics us ON u.id = us.user_id
            WHERE u.referrer_id = ts.user_id
            AND us.last_order_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
        )
    `
	db.Exec(sql)
}

// 每天0点执行
func ResetTodayNew(db *gorm.DB) {
	db.Model(&project.TeamStatistics{}).
		Where("today_new > 0").
		Update("today_new", 0)
}
