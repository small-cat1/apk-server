package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentApi struct {
}

func (p PaymentApi) GetPaymentProviders(c *gin.Context) {
	providers, err := paymentProviderService.GetAllPaymentProviders()
	if err != nil {
		global.GVA_LOG.Error("获取支付服务商列表失败", zap.Error(err))
		response.OkWithMessage("获取支付服务商列表失败", c)
		return
	}
	response.OkWithDetailed(providers, "success", c)
	return
}
