package web

import (
	"ApkAdmin/global"
	common "ApkAdmin/model/common/request"
	"ApkAdmin/model/common/response"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryApi struct {
}

func (a *CategoryApi) FirstCategory(c *gin.Context) {
	var idInfo common.GetById
	err := c.ShouldBindQuery(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	id := idInfo.Uint()
	data, err := CategoryService.GetCategory(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *CategoryApi) GetTrendingCategory(c *gin.Context) {
	categoryList, err := CategoryService.GetTrendingCategory()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(categoryList, "获取成功", c)
}

func (a *CategoryApi) FindCategory(c *gin.Context) {
	data, err := CategoryService.GetAllCategory()
	if err != nil {
		global.GVA_LOG.Error("获取分类列表的应用分类失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表的应用分类失败", c)
		return
	}
	response.OkWithDetailed(data, "获取分类列表的应用分类成功", c)
}

func (a *CategoryApi) FindAccountAppCategory(c *gin.Context) {
	data, err := CategoryService.GetAccountListAppCategory()
	if err != nil {
		global.GVA_LOG.Error("获取账户列表的应用分类失败!", zap.Error(err))
		response.FailWithMessage("获取账户列表的应用分类失败", c)
		return
	}
	response.OkWithDetailed(data, "获取账户列表的应用分类成功", c)
}
