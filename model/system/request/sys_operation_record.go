package request

import (
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
