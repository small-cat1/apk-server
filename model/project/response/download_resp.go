package response

type DownloadResp struct {
	CanDownload    bool   `json:"can_download"`    // 是否可以下载
	PackageUrl     string `json:"package_url"`     // 安装包地址
	PackageDetail  string `json:"package_detail"`  //安装包详情
	DownloadReason string `json:"download_reason"` //是否可以下载原因
}
