package response

type AliOssConfigResponse struct {
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	StsToken        string `json:"stsToken"`
	Bucket          string `json:"bucket"`
	Dir             string `json:"dir"` // 上传目录前缀
	Expiration      string `json:"expiration"`
}

type UploadSignatureResponse struct {
	SignedUrl  string `json:"signedUrl"`  // 签名的上传 URL
	ObjectName string `json:"objectName"` // OSS 对象路径
	Url        string `json:"url"`        // 最终访问 URL
	ExpireTime int64  `json:"expireTime"` // 过期时间戳
}
