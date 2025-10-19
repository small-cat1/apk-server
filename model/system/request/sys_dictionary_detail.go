package request

import (
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
