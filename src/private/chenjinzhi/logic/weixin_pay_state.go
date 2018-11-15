/*
	微信查看支付状态
*/
package logic

import (
	"customlib/log"
	"customlib/tool"
	"time"
	//"private/chenjinzhi/dbm"
)

const (
	QUERY_TIME_RANGE = 30 * time.Second
)

func QueryWxPayState() {
	defer func() {
		if err := recover(); nil != err {
			log.Error("QueryWxPayState.panic: %v", err)
			log.Error("%s", tool.PanicTrace())

			// 不停的把自己拉起来
			go QueryWxPayState()
		}
	}()

	log.Info("QueryWxPayState start")

	//	bStart := true

	//	//	var startId int64
	//	//	var startTm string

	//	for {
	//		if bStart {
	//			startId = 0
	//			// 只处理前一天的订单
	//			startTm := time.Now().Add(-24 * 60 * 60).Format(tool.TM_FMT_NORMAL)
	//		}

	//		time.Sleep(QUERY_TIME_RANGE)

	//	}

}
