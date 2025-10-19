package request

import (
	"ApkAdmin/utils"
	"errors"
)

type FilterAppRequest struct {
	PageInfo
	CategoryId uint   `json:"category_id" form:"categoryId" `
	SortType   string `json:"sort_type" form:"sortType" binding:"required,oneof=hot new rating downloads"`
	FilterType string `json:"filter_type" form:"filterType" binding:"required,oneof=all free paid premium"`
}

func (r FilterAppRequest) Validate() error {
	if r.CategoryId < 0 || r.CategoryId > 9999 {
		return errors.New("分类参数不正确")
	}
	SortTypes := []string{"hot", "new", "rating", "downloads"}
	if !utils.Contains(SortTypes, r.SortType) {
		return errors.New("排序类型不正确")
	}
	return nil
}

type FilterAccountAppRequest struct {
	PageInfo
	CategoryId uint   `json:"category_id" form:"categoryId" `
	SortType   string `json:"sort_type" form:"sortType" binding:"required"`
}

func (r FilterAccountAppRequest) Validate() error {
	if r.CategoryId < 0 || r.CategoryId > 9999 {
		return errors.New("分类参数不正确")
	}
	SortTypes := []string{"default", "price_asc", "price_desc", "time_desc", "time_asc", "sales_desc", "sales_asc"}
	if !utils.Contains(SortTypes, r.SortType) {
		return errors.New("排序类型不正确")
	}
	return nil
}

type SearchAppRequest struct {
	PageInfo
	Keyword  string `json:"keyword" form:"keyword" binding:"required"`
	SortType string `json:"sort_type" form:"sortType" binding:"required,oneof=hot new rating downloads"`
}
