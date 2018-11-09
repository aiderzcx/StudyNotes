/*
	扫码支付相关的功能集合
*/
package wx

import (
	"customlib/tool"
	"encoding/xml"
	"fmt"
	"time"
)

const (
	scan_code_link = "weixin://wxpay/bizpayurl?sign=%s&appid=%s&mch_id=%s&product_id=%s&time_stamp=%s&nonce_str=%s"
)

// 二维码连接, 输出产品ID
func CreateCodeLink(inProductId string) string {
	params := map[string]string{
		"appid":      g_conf.Scan.AppId,
		"mch_id":     g_conf.Scan.MchId,
		"product_id": inProductId,
		"time_stamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonce_str":  nonceString(),
	}

	sign := md5Sign(params, g_conf.Scan.SecretKey)

	return fmt.Sprintf(
		scan_code_link,
		sign, params["appid"], params["mch_id"],
		params["product_id"], params["time_stamp"], params["nonce_str"],
	)
}

//扫码回调的处理函数
type ScanCallBackReq struct {
	XMLName     xml.Name `xml:"xml"`
	AppId       string   `xml:"appid"`
	MchId       string   `xml:"mch_id"`
	OpenId      string   `xml:"openid"`
	IsSubScribe string   `xml:"is_subscribe"` // 是否关注公众账号 Y或N;Y-关注;N-未关注
	NonceStr    string   `xml:"nonce_str"`
	ProductId   string   `xml:"product_id"` // 商户定义的商品id 或者订单号
	Sign        string   `xml:"sign"`
}

func (s *ScanCallBackReq) ParamCheck() error {
	var merr tool.MultiError

	if s.AppId != g_conf.Scan.AppId {
		merr.AddErr("appid not equal")
	}

	if s.MchId != g_conf.Scan.MchId {
		merr.AddErr("mchid not equal")
	}

	if s.NonceStr == "" {
		merr.AddErr("noncestr is nil")
	}

	if s.ProductId == "" {
		merr.AddErr("productid is nil")
	}

	if s.Sign == "" {
		merr.AddErr("sign is nil")
	}

	err := verifyMd5Sign(s.ToMap(), g_conf.Scan.SecretKey)
	if nil != err {
		merr.AddErr(fmt.Sprintf("sign verify fail: %v", err))
	}

	return merr.Result()
}

func (s *ScanCallBackReq) ToMap() map[string]string {
	return map[string]string{
		"appid":        s.AppId,
		"mch_id":       s.MchId,
		"openid":       s.OpenId,
		"is_subscribe": s.IsSubScribe,
		"nonce_str":    s.NonceStr,
		"product_id":   s.ProductId,
		"sign":         s.Sign,
	}
}

type ScanCallBackResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	AppId      string   `xml:"appid"`
	MchId      string   `xml:"mch_id"`
	NonceStr   string   `xml:"nonce_str"`
	PrepayId   string   `xml:"prepay_id"`
	ResultCode string   `xml:"result_code"`
	ResultMsg  string   `xml:"err_code_des"`
	Sign       string   `xml:"sign"`
}

// 收到扫码通知
func ScanCallBack(inReq *ScanCallBackReq, outResp *ScanCallBackResp) error {
	// 1. 根据产品ID，获取订单信息

	// 2. 发起与支付流程

	// 3. 响应结果

	outResp.ResultCode = RESP_SUCC
	outResp.ReturnMsg = "ok"
	outResp.AppId = `xml:"appid"`
	outResp.MchId = `xml:"mch_id"`
	outResp.NonceStr = `xml:"nonce_str"`
	outResp.PrepayId = `xml:"prepay_id"`
	outResp.ResultCode = `xml:"result_code"`
	outResp.ResultMsg = `xml:"err_code_des"`
	outResp.Sign = `xml:"sign"`

	return nil
}

func ScanPrePay(inParam *PrePayParamST) (string, error) {
	return PrepayToWx(g_conf.Scan, inParam)
}

func ScanPayNotif(inParam map[string]string, outResp *PayNotifyResp) error {
	return PayNotify(g_conf.Scan, inParam, outResp)
}
