package service

import (
	"ApkAdmin/service/example"
	"ApkAdmin/service/project"
	"ApkAdmin/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	ProjectServiceGroup project.ServiceGroup
}
