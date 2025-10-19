package router

import (
	"ApkAdmin/router/example"
	"ApkAdmin/router/project"
	"ApkAdmin/router/system"
	"ApkAdmin/router/web"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Project project.RouterGroup
	Web     web.RouterGroup
}
