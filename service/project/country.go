package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CountryService struct{}

func (a *CountryService) Exists(code string) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.CountryRegion{}).Where("country_code = ?", code).Count(&count).Error
	return count > 0, err
}

func (a *CountryService) ExistsExcludeID(code string, excludeID uint) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.CountryRegion{}).Where("country_code = ? AND id != ?", code, excludeID).Count(&count).Error
	return count > 0, err
}

func (a *CountryService) GetCountry(id uint) (country project.CountryRegion, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&country).Error
	return
}

func (a *CountryService) GetByID(id uint) (*project.CountryRegion, error) {
	var country project.CountryRegion
	err := global.GVA_DB.First(&country, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}

func (a *CountryService) CreateCountry(req request.CountryCreateRequest) (err error) {
	// 检查国家代码是否已存在
	exists, err := a.Exists(req.CountryCode)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("国家代码 %s 已存在", req.CountryCode)
	}
	country := req.ToCountryRegion()
	err = global.GVA_DB.Create(&country).Error
	return err
}

func (a *CountryService) UpdateCountry(req *request.CountryUpdateRequest) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}
	// 检查国家代码是否被其他记录使用
	if existing.CountryCode != req.CountryCode {
		exists, err := a.ExistsExcludeID(req.CountryCode, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("国家代码 %s 已被其他记录使用", req.CountryCode)
		}
	}
	existing.CountryCode = req.CountryCode
	existing.CountryName = req.CountryName
	existing.CountryNameEN = req.CountryNameEN
	existing.Region = req.Region
	existing.CurrencyCode = req.CurrencyCode
	existing.LanguageCodes = req.LanguageCodes
	existing.ContentRatingSystem = req.ContentRatingSystem
	existing.IsSupported = &req.IsSupported

	err = global.GVA_DB.Debug().Omit("created_at").Updates(existing).Error
	return err
}

func (a *CountryService) DeleteCountry(cid int) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(uint(cid))
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}

	var t1 int64
	err = global.GVA_DB.Model(&project.AppPackage{}).Where("country_code = ?", cid).Count(&t1).Error
	if err != nil {
		return err
	}
	if t1 > 0 {
		return errors.New("该国家代码下存在安装包，请先解除使用")
	}
	var total int64
	err = global.GVA_DB.Model(&project.AppScreenshot{}).Where("country_code = ?", cid).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该国家代码下存在安装包截图，请先解除使用")
	}
	return global.GVA_DB.Model(&project.CountryRegion{}).Where("id = ?", cid).Delete(&project.CountryRegion{}).Error
}

func (a *CountryService) GetCountryList(info request.CountryListRequest, order string, desc bool) (list interface{}, total int64, err error) {
	var countryLists []project.CountryRegion
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&project.CountryRegion{})

	if info.Region != "" {
		db = db.Where("region = ?", info.Region)
	}

	err = db.Count(&total).Error
	if err != nil {
		return countryLists, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	err = db.Order(OrderStr).Find(&countryLists).Error
	return countryLists, total, err
}
