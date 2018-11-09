package router

import (
	"private/chenjinzhi/logic"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) error {
	// serialize config text for safe
	r.POST("/config/serialize", logic.ConfigSerialize)

	// 微信扫码支付，订单支付通知
	r.GET("/wx/scan/link", logic.WxScanLink)
	r.POST("/wx/scan/notify", logic.WxScanNotify)
	r.POST("/wx/scan/pay/notify", logic.WxScanPayNotify)

	return nil
}
