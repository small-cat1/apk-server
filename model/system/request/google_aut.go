package request

import (
	"ApkAdmin/constants"
	"ApkAdmin/utils"
	"errors"
	"strings"
)

type GoogleAuthRequest struct {
	Code   string `json:"code" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type GoogleAuthVerifyReq struct {
	Code   string `json:"code" binding:"required"`
	Action string `json:"action" binding:"required"`
}

func (r GoogleAuthVerifyReq) Validate() error {
	actions := strings.Split(constants.GoogleVerifyAction, ",")
	if !utils.Contains(actions, r.Action) {
		return errors.New("操作行为不允许")
	}
	return nil
}

type ViewSensitiveConfigRequest struct {
	Id   int    `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}
