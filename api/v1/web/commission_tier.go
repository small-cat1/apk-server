package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommissionTierApi struct {
}

func (a CommissionTierApi) GetCommissionTier(c *gin.Context) {
	tiers, err := commissionTierService.GetAllEnabledTiers()
	if err != nil {
		global.GVA_LOG.Error("获取分佣规则失败", zap.Error(err))
		response.OkWithMessage("获取分佣规则失败", c)
		return
	}
	withdrawalRules, err := systemConfigService.GetConfig("commission")
	if err != nil {
		global.GVA_LOG.Error("获取分佣提现规则失败", zap.Error(err))
		response.OkWithMessage("获取分佣提现规则失败", c)
		return
	}
	withdrawalRules["tiers"] = tiers
	response.OkWithDetailed(withdrawalRules, "success", c)
	return
}
