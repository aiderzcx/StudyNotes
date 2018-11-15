package wx

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

const (
	PAY_RESULT_SUCC = iota // 支付成功
	PAY_RESULT_FAIL        // 支付失败
	PAY_RESULT_DROP        // 消息错误，丢弃
)

type PayNotifyResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

type PayNtfResultST struct {
	Code      string
	WxOrderId string
	SpOrderId string
	TotalFee  int64
	PayCache  int64
	PayTime   string
	TradeType string
}

func (p *PayNtfResultST) check() bool {
	if p.WxOrderId == "" ||
		p.SpOrderId == "" ||
		p.TotalFee <= 0 ||
		p.PayCache < 0 ||
		p.TradeType == "" {
		return false
	}

	return true
}

const (
	WX_TM_LEN = 14 // 20060102150405
)

func PayNotify(inPayCfg ConfItemST, inParam map[string]string) (PayNtfResultST, error) {
	var result PayNtfResultST

	code, err := checkWxRespMsg(inParam, inPayCfg)
	if nil != err {
		return result, err
	}

	result.Code = code

	var ok bool
	result.SpOrderId, ok = inParam["out_trade_no"]
	if !ok {
		return result, fmt.Errorf("not out_trade_no field")
	}

	result.WxOrderId, ok = inParam["transaction_id"]
	if !ok {
		return result, fmt.Errorf("not transaction_id field")
	}

	result.TradeType, ok = inParam["trade_type"]
	if !ok {
		return result, fmt.Errorf("not trade_type field")
	}

	result.PayTime, ok = inParam["time_end"]
	if !ok || len(result.PayTime) != WX_TM_LEN {
		return result, fmt.Errorf("time_end('%s') error", result.PayTime)
	}

	result.PayCache, err = strconv.ParseInt(inParam["cash_fee"], 10, 64)
	if nil != err {
		return result, fmt.Errorf("CashFee.ParseInt(%s).%v", inParam["cash_fee"], err)
	}

	result.TotalFee, err = strconv.ParseInt(inParam["total_fee"], 10, 64)
	if nil != err {
		return result, fmt.Errorf("TotalFee.ParseInt(%s).%v", inParam["total_fee"], err)
	}

	if !result.check() {
		return result, fmt.Errorf("result is invalid, value(%+v)", result)
	}

	return result, nil
}
