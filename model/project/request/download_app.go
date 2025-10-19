package request

import (
	"ApkAdmin/utils"
	"errors"
	"strings"
)

type DownloadAppRequest struct {
	AppId  uint   `json:"appId" binding:"required"`
	OsType string `json:"osType" binding:"required"`
}

func (r DownloadAppRequest) Validate() error {
	if r.AppId <= 0 || r.AppId > 9999 {
		return errors.New("APPID参数不正确")
	}
	osType := strings.ToLower(r.OsType)
	validTypes := []string{"ios", "android"}
	if !utils.Contains(validTypes, osType) {
		return errors.New("系统类型不支持！")
	}
	return nil
}
