package project

import (
	"ApkAdmin/global"
	common "ApkAdmin/model/common/request"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
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

func (a *CategoryApi) GetSelectCategory(c *gin.Context) {
	var pageInfo request.SelectCategoryRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	categoryList, err := CategoryService.FindSelectCategory(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取Select应用分类列表失败!", zap.Error(err))
		response.FailWithMessage("获取Select应用分类列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(categoryList, "获取Select应用分类列表成功", c)
}

func (a *CategoryApi) GetCategoryList(c *gin.Context) {
	var pageInfo request.CategoryPageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage("获取应用分类列表失败，"+err.Error(), c)
		return
	}
	categoryList, total, err := CategoryService.GetCategoryList(pageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取应用分类列表失败!", zap.Error(err))
		response.FailWithMessage("获取应用分类列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     categoryList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取应用分类列表成功", c)
}

func (a *CategoryApi) AddCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	if err := CategoryService.CreateCategory(req); err != nil {
		global.GVA_LOG.Error("创建/更新失败!", zap.Error(err))
		response.FailWithMessage("创建/更新失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建/更新成功", c)
}

func (a *CategoryApi) UpdateCategory(c *gin.Context) {
	var req request.UpdateCategoryRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	err = CategoryService.UpdateCategory(&req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *CategoryApi) DeleteCategory(c *gin.Context) {
	var info common.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(info, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = CategoryService.DeleteCategory(info.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败,"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
