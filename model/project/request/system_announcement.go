package request

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type MarkAnnouncementReadRequest struct {
	AnnouncementID int64 `json:"announcement_id" binding:"required"`
}

type ListAnnouncementRequest struct {
	PageInfo
	Type int `json:"type" form:"type" `
}

type CreateAnnouncementRequest struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	DisplayType int    `json:"display_type" binding:"required"`
	TargetUsers string `json:"target_users" binding:"required"`
	LinkUrl     string `json:"link_url"`
	StartTime   string `json:"start_time" `
	EndTime     string `json:"end_time" `
	IsClosable  int    `json:"is_closable" binding:"required"`
	Status      *int   `json:"status" `
	// 内部字段，存储转换后的时间
	startTimeParsed *time.Time
	endTimeParsed   *time.Time
}

// Validate 验证创建公告请求
func (req *CreateAnnouncementRequest) Validate() error {
	// 验证标题
	if err := validateTitle(req.Title); err != nil {
		return err
	}

	// 验证内容
	if err := validateContent(req.Content); err != nil {
		return err
	}

	// 验证类型
	if err := validateType(req.Type); err != nil {
		return err
	}

	// 验证展示类型
	if err := validateDisplayType(req.DisplayType); err != nil {
		return err
	}

	// 验证目标用户
	if err := validateTargetUsers(req.TargetUsers); err != nil {
		return err
	}

	// 验证跳转链接（如果提供）
	if req.LinkUrl != "" {
		if err := validateLinkUrl(req.LinkUrl); err != nil {
			return err
		}
	}

	// 验证并转换时间，存储到内部字段
	if req.StartTime != "" || req.EndTime != "" {
		startTime, endTime, err := validateAndConvertTimeRange(req.StartTime, req.EndTime)
		if err != nil {
			return err
		}
		req.startTimeParsed = startTime
		req.endTimeParsed = endTime
	}

	// 验证是否可关闭
	if err := validateIsClosable(req.IsClosable); err != nil {
		return err
	}

	// 验证状态
	if err := validateStatus(*req.Status); err != nil {
		return err
	}

	return nil
}

// GetStartTime 获取转换后的开始时间
func (req *CreateAnnouncementRequest) GetStartTime() *time.Time {
	return req.startTimeParsed
}

// GetEndTime 获取转换后的结束时间
func (req *CreateAnnouncementRequest) GetEndTime() *time.Time {
	return req.endTimeParsed
}

// UpdateAnnouncementRequest 更新公告请求
type UpdateAnnouncementRequest struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	DisplayType int    `json:"display_type" binding:"required"`
	TargetUsers string `json:"target_users" binding:"required"`
	LinkUrl     string `json:"link_url"`
	StartTime   string `json:"start_time" `
	EndTime     string `json:"end_time" `
	IsClosable  int    `json:"is_closable" binding:"required"`
	Status      *int   `json:"status"`
	// 内部字段，存储转换后的时间
	startTimeParsed *time.Time
	endTimeParsed   *time.Time
}

// Validate 验证更新公告请求
func (req *UpdateAnnouncementRequest) Validate() error {
	// 如果字段被提供，则进行验证
	if req.Title != "" {
		if err := validateTitle(req.Title); err != nil {
			return err
		}
	}

	if req.Content != "" {
		if err := validateContent(req.Content); err != nil {
			return err
		}
	}

	if req.Type != 0 {
		if err := validateType(req.Type); err != nil {
			return err
		}
	}

	if req.DisplayType != 0 {
		if err := validateDisplayType(req.DisplayType); err != nil {
			return err
		}
	}

	if req.TargetUsers != "" {
		if err := validateTargetUsers(req.TargetUsers); err != nil {
			return err
		}
	}

	if req.LinkUrl != "" {
		if err := validateLinkUrl(req.LinkUrl); err != nil {
			return err
		}
	}

	if req.StartTime != "" || req.EndTime != "" {
		if req.StartTime != "" || req.EndTime != "" {
			startTime, endTime, err := validateAndConvertTimeRange(req.StartTime, req.EndTime)
			if err != nil {
				return err
			}
			req.startTimeParsed = startTime
			req.endTimeParsed = endTime
		}
	}

	if req.IsClosable != 0 {
		if err := validateIsClosable(req.IsClosable); err != nil {
			return err
		}
	}

	if req.Status != nil {
		if err := validateStatus(*req.Status); err != nil {
			return err
		}
	}

	return nil
}

// GetStartTime 获取转换后的开始时间
func (req *UpdateAnnouncementRequest) GetStartTime() *time.Time {
	return req.startTimeParsed
}

// GetEndTime 获取转换后的结束时间
func (req *UpdateAnnouncementRequest) GetEndTime() *time.Time {
	return req.endTimeParsed
}

// validateTitle 验证标题
func validateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("标题不能为空")
	}
	if len(title) > 200 {
		return errors.New("标题长度不能超过200个字符")
	}
	return nil
}

// validateContent 验证内容
func validateContent(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("内容不能为空")
	}
	// 根据数据库TEXT类型，最大约65535字节
	if len(content) > 65535 {
		return errors.New("内容长度超出限制")
	}
	return nil
}

// validateType 验证类型
func validateType(announcementType int) error {
	validTypes := map[int]string{
		1: "紧急",
		2: "重要",
		3: "普通",
	}
	if _, ok := validTypes[announcementType]; !ok {
		return fmt.Errorf("无效的公告类型，必须是1(紧急)、2(重要)或3(普通)")
	}
	return nil
}

// validateDisplayType 验证展示类型
func validateDisplayType(displayType int) error {
	validDisplayTypes := map[int]string{
		1: "横幅",
		2: "弹窗",
		3: "卡片",
		4: "仅消息中心",
	}
	if _, ok := validDisplayTypes[displayType]; !ok {
		return fmt.Errorf("无效的展示类型，必须是1(横幅)、2(弹窗)、3(卡片)或4(仅消息中心)")
	}
	return nil
}

// validateTargetUsers 验证目标用户
func validateTargetUsers(targetUsers string) error {
	targetUsers = strings.TrimSpace(targetUsers)
	if targetUsers == "" {
		return errors.New("目标用户不能为空")
	}

	validTargets := map[string]bool{
		"all": true,
		"vip": true,
		"new": true,
	}

	if !validTargets[targetUsers] {
		return fmt.Errorf("无效的目标用户，必须是all(全部)、vip(VIP)或new(新用户)")
	}
	return nil
}

// validateLinkUrl 验证跳转链接
func validateLinkUrl(linkUrl string) error {
	linkUrl = strings.TrimSpace(linkUrl)
	if linkUrl == "" {
		return nil // 链接可以为空
	}

	if len(linkUrl) > 500 {
		return errors.New("链接长度不能超过500个字符")
	}

	// 验证URL格式
	parsedUrl, err := url.ParseRequestURI(linkUrl)
	if err != nil {
		return fmt.Errorf("链接格式无效: %v", err)
	}

	// 验证协议
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return errors.New("链接协议必须是http或https")
	}

	return nil
}

// validateAndConvertTimeRange 验证并转换时间范围
func validateAndConvertTimeRange(startTime, endTime string) (*time.Time, *time.Time, error) {
	if startTime == "" || endTime == "" {
		return nil, nil, errors.New("开始时间和结束时间不能为空")
	}

	// 支持多种时间格式
	timeFormats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+08:00",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}

	var start, end time.Time
	var err error

	// 解析开始时间
	for _, format := range timeFormats {
		start, err = time.Parse(format, startTime)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("开始时间格式无效，支持格式: RFC3339 或 '2006-01-02 15:04:05'")
	}

	// 解析结束时间
	for _, format := range timeFormats {
		end, err = time.Parse(format, endTime)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("结束时间格式无效，支持格式: RFC3339 或 '2006-01-02 15:04:05'")
	}

	// 验证时间逻辑
	if !end.After(start) {
		return nil, nil, errors.New("结束时间必须晚于开始时间")
	}

	// 验证时间不能是过去的时间（可选，根据业务需求）
	now := time.Now()
	if end.Before(now) {
		return nil, nil, errors.New("结束时间不能早于当前时间")
	}

	return &start, &end, nil
}

// validateIsClosable 验证是否可关闭
func validateIsClosable(isClosable int) error {
	if isClosable != 0 && isClosable != 1 {
		return errors.New("是否可关闭字段必须是0(否)或1(是)")
	}
	return nil
}

// validateStatus 验证状态
func validateStatus(status int) error {
	if status != 0 && status != 1 {
		return errors.New("状态字段必须是0(草稿)或1(发布)")
	}
	return nil
}

// ValidateAnnouncementID 验证公告ID
func ValidateAnnouncementID(id int64) error {
	if id <= 0 {
		return errors.New("无效的公告ID")
	}
	return nil
}
