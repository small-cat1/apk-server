package utils

import (
	"path/filepath"
	"strings"
)

// GetFileType 根据文件扩展名判断类型
func GetFileType(filename string) (string, string) {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".webp": true, ".svg": true,
	}
	videoExts := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".wmv": true,
		".flv": true, ".mkv": true,
	}
	docExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".xls": true,
		".xlsx": true, ".ppt": true, ".pptx": true, ".txt": true,
	}
	packageExts := map[string]bool{
		".apk": true,
		".ipa": true,
		".exe": true,
		".dmg": true,
		".pkg": true,
		".zip": true,
		".rar": true,
		".7z":  true,
	}
	if imageExts[ext] {
		return "image", ext
	} else if videoExts[ext] {
		return "video", ext
	} else if docExts[ext] {
		return "document", ext
	} else if packageExts[ext] {
		return "package", ext
	}

	return "file", ext
}
