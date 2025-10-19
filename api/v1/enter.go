package v1

import (
	"ApkAdmin/api/v1/example"
	"ApkAdmin/api/v1/project"
	"ApkAdmin/api/v1/system"
	"ApkAdmin/api/v1/web"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	ProjectApiGroup project.ApiGroup
	WebApiGroup     web.ApiGroup
}
