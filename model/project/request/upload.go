package request

type UploadSignatureRequest struct {
	FileName string `json:"fileName" binding:"required"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
}
