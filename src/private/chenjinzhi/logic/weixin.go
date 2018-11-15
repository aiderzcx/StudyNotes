package logic

import (
	"bytes"
	"customlib/log"
	"customlib/mservice/gincom"
	"customlib/mservice/mtrace"
	"customlib/mservice/srvmsg"
	"customlib/tool"
	"fmt"
	"io/ioutil"
	"private/chenjinzhi/dbm"
	"private/chenjinzhi/wapi"
	"private/netpay/weixin"

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

	db := orm.NewOrm()
	// 1 获取产品信息
	order, err := dbm.StartPrepay(db, req.ProductId)
	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "query db error"
		logger.Warning("dbm.ProductOrder(%s).%v", req.ProductId, err)
		return
	}

	// 如果状态是等待支付通知以上，则直接发回成功
	if order.State >= dbm.PAY_STATE_PREPAY_FAIL {
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
	param.TradeType = order.PayType
	param.ProductId = order.ProductId
	param.Logger = &logger

	result, err := wx.ScanPrePay(&param)
	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "prepay error"
		logger.Warning("wx.ScanPrePay(%+v).%v", param, err)
		return
	}

	if result.Code == wx.RESP_DROP {
		respData.ResultCode = wx.RESP_SUCC
		respData.ResultMsg = "drop msg"
		logger.Warning("wx.ScanPrePay(%+v) drop msg", param)
		return
	}

	order.ThirdId = result.PrepayId
	order.PayType = result.TradeType
	if result.Code == wx.RESP_SUCC {
		order.State = dbm.PAY_STATE_PREPAYED
	} else {
		order.State = dbm.PAY_STATE_PREPAY_FAIL
		order.Remarks = fmt.Sprintf("[%s]", result.Code)
	}

	err = dbm.PrePayResp(db, order)
	if nil != err {
		respData.ResultCode = wx.RESP_FAIL
		respData.ResultMsg = "update db fail"
		logger.Warning("dbm.PrePaySucc().%v", err)
		return
	}

	// 3. 响应结果
	respData.Update(result.PrepayId)
	return

}

func WxScanPayNotify(ctx *gin.Context) {
	var respData wx.PayNotifyResp
	respData.ReturnCode = wx.RESP_FAIL

	logger, dataMap, err := xmlReqToMap(ctx, "WxScanPayNotify")
	defer gincom.RespXmlData(ctx, &respData, logger)

	if nil != err {
		logger.Warning("gincom.Prepare().%v", err)
		respData.ReturnMsg = "msg err"
		return
	}

	logger.Info("WxScanPayNotify req: %+v", dataMap)

	// 1 调用微信的支付通知函数处理
	result, err := wx.ScanPayNotif(dataMap)
	if nil != err {
		respData.ReturnMsg = "msg err"
		logger.Warning("wx.ScanPayNotif().%v", err)
		return
	}

	logger.Debug("ScanPayNotif.Result(%+v)", result)

	// 2 更新数据库信息
	db := orm.NewOrm()
	order, err := dbm.OrderInfo(db, result.SpOrderId)
	if nil != err {
		respData.ReturnMsg = "db.query error"
		logger.Warning("dbm.OrderInfo(%s).%v", result.SpOrderId, err)
		return
	}

	if order.PayType != result.TradeType ||
		order.TotalFee != result.TotalFee ||
		order.ThirdId != result.WxOrderId {
		respData.ReturnMsg = "params not matched"
		logger.Warning("payType,totalFee,ThirdId not matched")
		return
	}

	if order.State >= dbm.PAY_STATE_PAY_SUCC {
		logger.Warning("the order(%s).state(%d) recv pay notify", order.PayId, order.State)
		respData.ReturnCode = wx.RESP_SUCC
		respData.ReturnMsg = "ok"
		return
	}

	order.PayAt = transWxTm2Db(result.PayTime)
	if result.Code == wx.RESP_SUCC {
		order.State = dbm.PAY_STATE_PAY_SUCC
	} else {
		order.State = dbm.PAY_STATE_PAY_FAIL
		order.Remarks += fmt.Sprintf("[%s]", result.Code)
	}

	err = dbm.PayCallback(db, order)
	if nil != err {
		respData.ReturnMsg = "db.update error"
		logger.Warning("dbm.PayCallback(%s).%v", result.SpOrderId, err)
		return
	}

	respData.ReturnCode = wx.RESP_SUCC
	respData.ReturnMsg = "ok"
	return
}

func transWxTm2Db(inTm string) string {
	if inTm == "" {
		return tool.CurTimeNormal()
	}

	var buf bytes.Buffer

	buf.WriteString(inTm[0:4]) // yyyy
	buf.WriteString("-")
	buf.WriteString(inTm[4:6]) //MM
	buf.WriteString("-")
	buf.WriteString(inTm[6:8]) // DD
	buf.WriteString(" ")
	buf.WriteString(inTm[8:10]) // HH
	buf.WriteString(":")
	buf.WriteString(inTm[10:12]) // MM
	buf.WriteString(":")
	buf.WriteString(inTm[12:14]) // SS

	return buf.String()
}
