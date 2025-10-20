package utils

import (
	"path/filepath"
	"strings"
)

// GetFileType 根据文件扩展名判断类型
func GetFileType(filename string) string {
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
	if imageExts[ext] {
		return "image"
	} else if videoExts[ext] {
		return "video"
	} else if docExts[ext] {
		return "document"
	}

	return "file"
}
