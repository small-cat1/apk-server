package example

import (
	api "ApkAdmin/api/v1"
)

type RouterGroup struct {
	FileUploadAndDownloadRouter
}

var (
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
)
