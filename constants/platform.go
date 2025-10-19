package constants

import (
	"strings"
)

// Platform 平台类型枚举
type Platform string

const (
	// 移动端平台
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformHarmony Platform = "harmony" // 鸿蒙

	// 桌面端平台
	PlatformWindows Platform = "windows"
	PlatformMacOS   Platform = "macos"
	PlatformLinux   Platform = "linux"

	// Web 平台
	PlatformWeb Platform = "web"

	// 其他平台
	PlatformUnknown Platform = "unknown"
)

// String 实现 Stringer 接口
func (p Platform) String() string {
	return string(p)
}

// IsValid 判断平台类型是否有效
func (p Platform) IsValid() bool {
	switch p {
	case PlatformAndroid, PlatformIOS, PlatformHarmony,
		PlatformWindows, PlatformMacOS, PlatformLinux,
		PlatformWeb, PlatformUnknown:
		return true
	default:
		return false
	}
}

// IsMobile 判断是否为移动端平台
func (p Platform) IsMobile() bool {
	return p == PlatformAndroid || p == PlatformIOS || p == PlatformHarmony
}

// IsDesktop 判断是否为桌面端平台
func (p Platform) IsDesktop() bool {
	return p == PlatformWindows || p == PlatformMacOS || p == PlatformLinux
}

// IsWeb 判断是否为 Web 平台
func (p Platform) IsWeb() bool {
	return p == PlatformWeb
}

// GetDisplayName 获取平台的显示名称（中文）
func (p Platform) GetDisplayName() string {
	switch p {
	case PlatformAndroid:
		return "安卓"
	case PlatformIOS:
		return "苹果"
	case PlatformHarmony:
		return "鸿蒙"
	case PlatformWindows:
		return "Windows"
	case PlatformMacOS:
		return "macOS"
	case PlatformLinux:
		return "Linux"
	case PlatformWeb:
		return "网页版"
	default:
		return "未知"
	}
}

// GetIcon 获取平台图标（可用于前端显示）
func (p Platform) GetIcon() string {
	switch p {
	case PlatformAndroid:
		return "android"
	case PlatformIOS:
		return "apple"
	case PlatformHarmony:
		return "harmony"
	case PlatformWindows:
		return "windows"
	case PlatformMacOS:
		return "apple"
	case PlatformLinux:
		return "linux"
	case PlatformWeb:
		return "global"
	default:
		return "question"
	}
}

// GetDownloadFileExtension 获取下载文件扩展名
func (p Platform) GetDownloadFileExtension() string {
	switch p {
	case PlatformAndroid:
		return ".apk"
	case PlatformIOS:
		return ".ipa"
	case PlatformHarmony:
		return ".hap"
	case PlatformWindows:
		return ".exe"
	case PlatformMacOS:
		return ".dmg"
	case PlatformLinux:
		return ".deb" // 或 .rpm、.AppImage
	default:
		return ""
	}
}

// ParsePlatform 从字符串解析平台类型（不区分大小写）
func ParsePlatform(s string) Platform {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "android":
		return PlatformAndroid
	case "ios", "iphone", "ipad":
		return PlatformIOS
	case "harmony", "harmonyos", "hongmeng":
		return PlatformHarmony
	case "windows", "win":
		return PlatformWindows
	case "macos", "mac", "osx", "darwin":
		return PlatformMacOS
	case "linux":
		return PlatformLinux
	case "web", "browser":
		return PlatformWeb
	default:
		return PlatformUnknown
	}
}

// GetAllPlatforms 获取所有平台列表
func GetAllPlatforms() []Platform {
	return []Platform{
		PlatformAndroid,
		PlatformIOS,
		PlatformHarmony,
		PlatformWindows,
		PlatformMacOS,
		PlatformLinux,
		PlatformWeb,
	}
}

// GetMobilePlatforms 获取所有移动端平台
func GetMobilePlatforms() []Platform {
	return []Platform{
		PlatformAndroid,
		PlatformIOS,
		PlatformHarmony,
	}
}

// GetDesktopPlatforms 获取所有桌面端平台
func GetDesktopPlatforms() []Platform {
	return []Platform{
		PlatformWindows,
		PlatformMacOS,
		PlatformLinux,
	}
}

// MarshalJSON 实现 JSON 序列化
func (p Platform) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(p) + `"`), nil
}

// UnmarshalJSON 实现 JSON 反序列化
func (p *Platform) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	*p = ParsePlatform(str)
	return nil
}

// PlatformInfo 平台详细信息
type PlatformInfo struct {
	Platform    Platform `json:"platform"`
	DisplayName string   `json:"displayName"`
	Icon        string   `json:"icon"`
	IsMobile    bool     `json:"isMobile"`
	IsDesktop   bool     `json:"isDesktop"`
	Extension   string   `json:"extension"`
}

// GetInfo 获取平台详细信息
func (p Platform) GetInfo() PlatformInfo {
	return PlatformInfo{
		Platform:    p,
		DisplayName: p.GetDisplayName(),
		Icon:        p.GetIcon(),
		IsMobile:    p.IsMobile(),
		IsDesktop:   p.IsDesktop(),
		Extension:   p.GetDownloadFileExtension(),
	}
}
