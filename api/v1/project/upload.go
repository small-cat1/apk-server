package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	projectReq "ApkAdmin/model/project/request"
	projectResp "ApkAdmin/model/project/response"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
	"time"
)

type UploadApi struct {
}

// GetOssConfig 获取 OSS 上传配置
func (u *UploadApi) GetOssConfig(c *gin.Context) {
	// 1. 创建 STS 客户端
	client, err := sts.NewClientWithAccessKey(
		global.GVA_CONFIG.AliyunOSS.Endpoint,
		global.GVA_CONFIG.AliyunOSS.AccessKeyId,
		global.GVA_CONFIG.AliyunOSS.AccessKeySecret,
	)
	if err != nil {
		response.FailWithMessage("创建STS客户端失败", c)
		return
	}

	// 2. 构建 AssumeRole 请求
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = ""                       // RAM 角色 ARN global.GVA_CONFIG.AliyunOSS.RoleArn
	request.RoleSessionName = "upload-session" // 会话名称
	request.DurationSeconds = "3600"           // 凭证有效期（秒）1小时

	// 可选：限制权限策略
	policy := `{
		"Version": "1",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"oss:PutObject",
					"oss:GetObject"
				],
				"Resource": [
					"acs:oss:*:*:` + global.GVA_CONFIG.AliyunOSS.BucketName + `/apk/*",
					"acs:oss:*:*:` + global.GVA_CONFIG.AliyunOSS.BucketName + `/ipa/*"
				]
			}
		]
	}`
	request.Policy = policy

	// 3. 获取临时凭证
	stsResponse, err := client.AssumeRole(request)
	if err != nil {
		response.FailWithMessage("获取临时凭证失败: "+err.Error(), c)
		return
	}

	// 4. 返回配置信息
	ossConfig := projectResp.AliOssConfigResponse{
		Region:          global.GVA_CONFIG.AliyunOSS.Region,
		AccessKeyId:     stsResponse.Credentials.AccessKeyId,
		AccessKeySecret: stsResponse.Credentials.AccessKeySecret,
		StsToken:        stsResponse.Credentials.SecurityToken,
		Bucket:          global.GVA_CONFIG.AliyunOSS.BucketName,
		Dir:             "apk", // 或根据业务需要动态设置
		Expiration:      stsResponse.Credentials.Expiration,
	}
	response.OkWithData(ossConfig, c)
}

// GetUploadSignature 生成 OSS 上传签名 URL
func (u *UploadApi) GetUploadSignature(c *gin.Context) {
	var req projectReq.UploadSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 1. 生成唯一的对象名称
	timestamp := time.Now().UnixNano() / 1e6
	objectName := fmt.Sprintf("private/package/%d/%s", timestamp, req.FileName)

	// 2. 设置过期时间（1小时）
	expireTime := time.Now().Add(1 * time.Hour).Unix()

	// 3. 构建签名 URL
	bucket := global.GVA_CONFIG.AliyunOSS.BucketName
	region := global.GVA_CONFIG.AliyunOSS.Region
	accessKeyId := global.GVA_CONFIG.AliyunOSS.AccessKeyId
	accessKeySecret := global.GVA_CONFIG.AliyunOSS.AccessKeySecret

	// OSS 域名
	host := fmt.Sprintf("https://%s.%s.aliyuncs.com", bucket, region)

	// 构建要签名的字符串
	method := "PUT"
	contentType := req.FileType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	date := fmt.Sprintf("%d", expireTime)

	// 规范化的资源路径
	canonicalizedResource := fmt.Sprintf("/%s/%s", bucket, objectName)

	// 构建待签名字符串
	stringToSign := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		method,
		contentType,
		date,
		canonicalizedResource,
	)

	// 4. 生成签名
	signature := generateSignature(stringToSign, accessKeySecret)

	// 5. 构建完整的签名 URL
	signedUrl := fmt.Sprintf("%s/%s?OSSAccessKeyId=%s&Expires=%s&Signature=%s",
		host,
		objectName,
		accessKeyId,
		date,
		signature,
	)

	// 6. 返回结果
	result := projectResp.UploadSignatureResponse{
		SignedUrl:  signedUrl,
		ObjectName: objectName,
		Url:        fmt.Sprintf("%s/%s", host, objectName),
		ExpireTime: expireTime,
	}

	response.OkWithData(result, c)
}

// generateSignature 生成 OSS 签名
func generateSignature(stringToSign, accessKeySecret string) string {
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(stringToSign))
	signedBytes := h.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}
