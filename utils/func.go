package utils

import (
	"fmt"
	"time"
)

// 工具函数
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GenerateFlowNo(prefix string) string {
	now := time.Now()
	// 格式：COM20251022001234567890
	return fmt.Sprintf("%s%s%09d",
		prefix,
		now.Format("20060102"),
		now.UnixNano()%1000000000,
	)
}
