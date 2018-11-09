/*
	微信的参数管理
*/

package wx

import (
	"fmt"
)

const (
	msg_succ         = iota // 消息check成功
	msg_error               // 消息错误丢弃消息
	return_code_fail        // returnCode错误，重试
	result_code_fail        // resultCode错误，流程结束
)

const (
	RESP_SUCC = "SUCCESS"
	RESP_FAIL = "FAIL"
)

const (
	min_param_len = 2 // 至少2个参数 return_code,return_msg
)

// 检测微信的Resp消息
func checkWxRespMsg(inParam map[string]string, inKey string) (int, error) {
	if len(inParam) <= min_param_len {
		return msg_error, fmt.Errorf("param not enough")
	}

	returnCode, exist := inParam["return_code"]
	if !exist {
		return msg_error, fmt.Errorf("no return_code")
	}

	if returnCode != "SUCCESS" {
		return return_code_fail, fmt.Errorf("return Code(%s), msg(%s)", returnCode, inParam["return_msg"])
	}

	_, exist = inParam["sign"]
	if !exist {
		return msg_error, fmt.Errorf("no sign")
	}

	resultCode, exist := inParam["result_code"]
	if !exist {
		return msg_error, fmt.Errorf("no result_code")
	}

	if resultCode != "SUCCESS" {
		return result_code_fail, fmt.Errorf("result Code(%s), msg(%s)", resultCode, inParam["err_code_des"])
	}

	// 消息鉴权
	err := verifyMd5Sign(inParam, inKey)
	if nil != err {
		return msg_error, err
	}

	return msg_succ, nil
}
