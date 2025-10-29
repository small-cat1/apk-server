package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WithdrawApi struct {
}

func (u *WithdrawApi) Withdraw(c *gin.Context) {
	// 1. 获取当前用户ID
	userID := utils.GetUserID(c)
	if userID <= 0 {
		response.FailWithMessage("用户未登录", c)
		return
	}
	var req request.UserWithdrawRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户提现参数不正确！", zap.Error(err))
		response.OkWithMessage("参数不正确！", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.OkWithMessage(err.Error(), c)
		return
	}
	// 3. 调用服务层
	if err := UserService.ApplyWithdraw(userID, req); err != nil {
		global.GVA_LOG.Error("申请提现失败",
			zap.Uint("userID", userID),
			zap.Float64("amount", req.Amount),
			zap.String("type", req.WithdrawType),
			zap.Error(err),
		)
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 4. 返回成功
	global.GVA_LOG.Info("申请提现成功",
		zap.Uint("userID", uint(userID)),
		zap.Float64("amount", req.Amount),
		zap.String("type", req.WithdrawType),
	)
	response.OkWithMessage("提现申请已提交，请等待审核", c)
}

func (u *WithdrawApi) GetWithdrawConfig(c *gin.Context) {
	config, err := systemConfigService.GetConfig("commission")
	if err != nil {
		global.GVA_LOG.Error("获取提现规则失败!", zap.Error(err))
		response.FailWithMessage("获取提现规则失败", c)
		return
	}
	response.OkWithData(config, c)
}

func (u *WithdrawApi) GetRecords(c *gin.Context) {
	var req request.WithdrawRecordRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取提现记录请求参数错误!", zap.Error(err))
		response.FailWithMessage("获取提现记录请求参数错误", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	result, err := UserService.GetWithdrawRecord(userId, req)
	if err != nil {
		global.GVA_LOG.Error("获取提现记录失败!", zap.Error(err))
		response.FailWithMessage("获取提现记录失败", c)
		return
	}
	response.OkWithDetailed(result, "获取提现记录列表成功", c)
}
