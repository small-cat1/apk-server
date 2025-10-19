package example

import (
	api "ApkAdmin/api/v1"
)

type RouterGroup struct {
	FileUploadAndDownloadRouter
	AttachmentCategoryRouter
}

var (
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
	attachmentCategoryApi       = api.ApiGroupApp.ExampleApiGroup.AttachmentCategoryApi
)
