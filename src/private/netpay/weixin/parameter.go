/*
	微信的参数管理
*/

package wx

import (
	"fmt"
)

const (
	RESP_SUCC = "SUCCESS"
	RESP_FAIL = "FAIL"
	RESP_DROP = "DROP" // SYSTEMERROR	系统错误	系统超时	系统异常，请用相同参数重新调用
)

const (
	min_param_len = 2 // 至少2个参数 return_code,return_msg
)

// 检测微信的Resp消息
func checkWxRespMsg(inParam map[string]string, inConf ConfItemST) (string, error) {
	if len(inParam) <= min_param_len {
		return RESP_FAIL, fmt.Errorf("param not enough")
	}

	if inParam["appid"] != inConf.AppId ||
		inParam["mch_id"] != inConf.MchId {
		return RESP_FAIL, fmt.Errorf("apid or mchid not matched")
	}

	returnCode, exist := inParam["return_code"]
	if !exist {
		return RESP_FAIL, fmt.Errorf("no return_code")
	}

	if returnCode != RESP_SUCC {
		return RESP_FAIL, fmt.Errorf("return Code(%s), msg(%s)", returnCode, inParam["return_msg"])
	}

	_, exist = inParam["sign"]
	if !exist {
		return RESP_FAIL, fmt.Errorf("no sign")
	}

	resultCode, exist := inParam["result_code"]
	if !exist {
		return RESP_FAIL, fmt.Errorf("no result_code")
	}

	// 消息鉴权
	err := verifyMd5Sign(inParam, inConf.SecretKey)
	if nil != err {
		return RESP_FAIL, err
	}

	// SYSTEMERROR	系统错误	系统超时	系统异常，请用相同参数重新调用
	if resultCode != RESP_SUCC {
		if resultCode == "SYSTEMERROR" {
			return RESP_DROP, nil
		} else {
			return inParam["err_code"], nil
		}
	}

	return RESP_SUCC, nil
}
