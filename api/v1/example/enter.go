package example

import "ApkAdmin/service"

type ApiGroup struct {
	FileUploadAndDownloadApi
	AttachmentCategoryApi
}

var (
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
	attachmentCategoryService    = service.ServiceGroupApp.ExampleServiceGroup.AttachmentCategoryService
)
