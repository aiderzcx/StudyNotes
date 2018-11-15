package wx

import (
	"customlib/log"
	"customlib/mservice/httpc"
	"customlib/tool"
	"encoding/xml"
	"fmt"
)

const (
	pre_pay_url = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

const (
	NTF_APP_ID     = "appid"
	NTF_MCH_ID     = "mch_id"
	NTF_TRADE_TYPE = "trade_type"
	NTF_PREPAY_ID  = "prepay_id"
)

type PrePayReq struct {
	XMLName   xml.Name `xml:"xml"`
	AppId     string   `xml:"appid"`
	MchId     string   `xml:"mch_id"`
	OrderDesc string   `xml:"body"`             // 商品描述
	OrderId   string   `xml:"out_trade_no"`     // 商户订单号 要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一
	TotalFee  int64    `xml:"total_fee"`        // 订单总金额，单位为分
	ServerIp  string   `xml:"spbill_create_ip"` // APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP
	NotifyUrl string   `xml:"notify_url"`       // 异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数
	TradeType string   `xml:"trade_type"`       // JSAPI 公众号支付,NATIVE 扫码支付,APP APP支付
	ProductId string   `xml:"product_id"`       // trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID
	OpenId    string   `xml:"openid"`           // trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识
	NonceStr  string   `xml:"nonce_str"`
	Sign      string   `xml:"sign"`
}

func (p *PrePayReq) ToMap() map[string]string {
	result := map[string]string{
		"appid":            p.AppId,
		"mch_id":           p.MchId,
		"body":             p.OrderDesc,
		"out_trade_no":     p.OrderId,
		"total_fee":        fmt.Sprintf("%d", p.TotalFee),
		"spbill_create_ip": p.ServerIp,
		"notify_url":       p.NotifyUrl,
		"trade_type":       p.TradeType,
		"nonce_str":        p.NonceStr,
		"sign":             p.Sign,
	}

	if p.TradeType == "NATIVE" {
		result["product_id"] = p.ProductId
	}

	if p.TradeType == "JSAPI" {
		result["openid"] = p.OpenId
	}

	return result
}

type PrePayParamST struct {
	OrderId   string
	OrderDesc string
	TotalFee  int64
	TradeType string
	ProductId string
	OpenId    string
	Logger    *log.LoggerST
}

type PrePayResultST struct {
	Code      string
	TradeType string
	PrepayId  string
	CodeUrl   string
}

// 返回与支付ID
func PrepayToWx(inConf ConfItemST, inParam *PrePayParamST) (PrePayResultST, error) {
	var req PrePayReq
	req.AppId = inConf.AppId
	req.MchId = inConf.MchId
	req.ServerIp = inConf.ServerIp
	req.NotifyUrl = inConf.CallBackUrl

	req.OrderDesc = inParam.OrderDesc
	req.OrderId = inParam.OrderId
	req.TotalFee = inParam.TotalFee
	req.TradeType = inParam.TradeType
	req.ProductId = inParam.ProductId
	req.OpenId = inParam.OpenId
	req.NonceStr = nonceString()
	req.Sign = md5Sign(req.ToMap(), inConf.SecretKey)

	reqData := req.ToMap()
	inParam.Logger.Debug("prePay.req: %+v", reqData)

	var result PrePayResultST
	respBuf, err := httpc.TlsPostXml(pre_pay_url, tool.NewXmlTrans(reqData), 0)
	if nil != err {
		return result, fmt.Errorf("TlsPostXml().%v", err)
	}

	inParam.Logger.Debug("prePay.Resp: %s", string(respBuf))

	respData, err := tool.XmlTransToMap(respBuf)
	if nil != err {
		return result, fmt.Errorf("XmlTransToMap().%v", err)
	}

	code, err := checkWxRespMsg(respData, inConf)
	if nil != err {
		return result, fmt.Errorf("checkWxRespMsg(): %s %v", code, err)
	}

	var ok bool
	result.Code = code
	result.TradeType, ok = respData["trade_type"]
	if !ok {
		return result, fmt.Errorf("no trade_type")
	}

	result.PrepayId, ok = respData["prepay_id"]
	if !ok {
		return result, fmt.Errorf("no prepay_id")
	}

	result.CodeUrl = respData["code_url"] // 可选参数
	return result, nil
}
