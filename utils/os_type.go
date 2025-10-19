package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// OSType 操作系统类型
type OSType string

const (
	OSTypeIOS     OSType = "ios"
	OSTypeAndroid OSType = "android"
	OSTypeUnknown OSType = "unknown"
)

// GetOSType 获取操作系统类型
func GetOSType(c *gin.Context) OSType {
	userAgent := c.Request.UserAgent()
	return ParseOSType(userAgent)
}

// ParseOSType 解析 User-Agent 字符串获取操作系统类型
func ParseOSType(userAgent string) OSType {
	userAgent = strings.ToLower(userAgent)

	// 判断 iOS（包括 iPhone、iPad、iPod）
	if strings.Contains(userAgent, "iphone") ||
		strings.Contains(userAgent, "ipad") ||
		strings.Contains(userAgent, "ipod") {
		return OSTypeIOS
	}

	// 判断 Android
	if strings.Contains(userAgent, "android") {
		return OSTypeAndroid
	}

	return OSTypeUnknown
}

// IsIOS 判断是否为 iOS
func IsIOS(c *gin.Context) bool {
	return GetOSType(c) == OSTypeIOS
}

// IsAndroid 判断是否为 Android
func IsAndroid(c *gin.Context) bool {
	return GetOSType(c) == OSTypeAndroid
}

// IsMobile 判断是否为移动设备
func IsMobile(c *gin.Context) bool {
	osType := GetOSType(c)
	return osType == OSTypeIOS || osType == OSTypeAndroid
}
