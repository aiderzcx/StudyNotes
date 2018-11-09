package logic

import (
	"customlib/log"
	"customlib/mservice/gincom"
	"customlib/mservice/mtrace"
	"customlib/mservice/srvmsg"
	"customlib/tool"
	"fmt"
	"io/ioutil"
	"production/chenjinzhi/dbm"
	"production/chenjinzhi/wapi"
	"production/netpay/weixin"

	"github.com/astaxie/beego/orm"
	"github.com/gin-gonic/gin"
)

func xmlReqToMap(ctx *gin.Context, inFunc string) (*log.LoggerST, map[string]string, error) {
	// 生成日志对象
	logger := log.WithFields(
		map[string]interface{}{
			mtrace.TraceKey: ctx.GetString(mtrace.TraceKey),
			"function":      inFunc,
		})

	data, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()

	if nil != err {
		return &logger, nil, fmt.Errorf("ioutil.ReadAll().%v", err)
	}

	dataMap, err := tool.XmlTransToMap(data)
	if nil != err {
		return &logger, nil, fmt.Errorf("tool.XmlTransToMap().%v", err)
	}

	return &logger, dataMap, nil
}

func WxScanLink(ctx *gin.Context) {
	var req wapi.ScanLinkReq
	var respData wapi.ScanLinkResp
	var resp srvmsg.RespMsgST

	logger, err := gincom.Prepare("WxScanNotif", ctx, &req, &resp)
	defer gincom.RespJsonData(ctx, &resp, &respData, &logger)

	if nil != err {
		logger.Warning("gincom.Prepare().%v", err)
		return
	}

	respData.Link = wx.CreateCodeLink(req.ProductId)

	resp.SetCode(srvmsg.RESP_CODE_SUCC)
	return
}

func WxScanNotify(ctx *gin.Context) {
	var req wx.ScanCallBackReq
	var respData wx.ScanCallBackResp
	var resp srvmsg.RespMsgST

	logger, err := gincom.Prepare("WxScanNotif", ctx, &req, &resp)
	defer gincom.RespXmlData(ctx, &respData, &logger)

	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "param error"
		logger.Warning("gincom.Prepare().%v", err)
		return
	}

	logger.Info("WxScanNotify req: %+v", req)

	db := orm.NewOrmer()
	// 1 获取产品信息
	order, err := dbm.ProductOrder(db, req.ProductId)
	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "query db error"
		logger.Warning("dbm.ProductOrder(%s).%v", req.ProductId, err)
		return
	}

	// 如果状态是等待支付通知以上，则直接发回成功
	if order.State >= db.PAY_STATE_WAIT_NTF {
		respData.ResultCode = wx.RESP_SUCC
		respData.ResultMsg = "repeat notify"
		logger.Info("order state(%d) not init", order.State)
		return
	}

	// 2 发起与支付
	var param wx.PrePayParamST
	param.OrderId = order.PayId
	param.OrderDesc = order.ProductDesc
	param.TotalFee = order.TotalFee
	param.TradeType = order.TradeType
	param.ProductId = order.ProductId
	param.Logger = &logger

	prePayId, err := wx.ScanPrePay(&param)
	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "prepay error"
		logger.Warning("wx.ScanPrePay(%+v).%v", param, err)
		return
	}

	// 3. 响应结果
	err = wx.ScanCallBack(&req, &respData)
	if err != err {
		logger.Warning("wx.ScanCallBack().%v", err)
	}

	return

}

func WxScanPayNotify(ctx *gin.Context) {
	var respData wx.PayNotifyResp

	logger, dataMap, err := xmlReqToMap(ctx, "WxScanPayNotify")
	defer gincom.RespXmlData(ctx, &respData, logger)

	if nil != err {
		logger.Warning("gincom.Prepare().%v", err)
		return
	}

	logger.Info("WxScanPayNotify req: %+v", dataMap)
	respData.ReturnCode = wx.RESP_SUCC
	respData.ReturnMsg = "ok"

	return
}
