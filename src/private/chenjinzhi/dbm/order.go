package dbm

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(PayOrderTB))
}

const (
	PAY_STATE_NIL       = iota // 空状态
	PAY_STATE_SCANED           // 收到扫码支付
	PAY_STATE_PREPAYING        // 发起预支付，等待支付回调
	PAY_STATE_PREPAY_FAIL
	PAY_STATE_SUCC      // 支付成功
	PAY_STATE_FAIL      // 支付失败
	PAY_STATE_REFUNDING // 正在退款
	PAY_STATE_REFUNDED  // 已经退款
)

type PayOrderTB struct {
	Id          int64 `orm:"pk;auto"`
	ProductId   string
	ProductDesc string
	TotalFee    int64
	PayId       string
	PayType     string
	ThirdId     string
	State       int
	QueryCount  int
	CreateAt    string
	ModifyAt    string
}

func (p *PayOrderTB) TableName() string {
	return "pay_order"
}

func ProductOrder(db orm.Ormer, inProductId string) (*PayOrderTB, error) {
	var order PayOrderTB

	err := db.QueryTable(&order).Filter("production_id", inProductId).One(&order)
	if nil != err {
		return nil, err
	}

	// 没有订单，更新订单
	if order.PayId == "" {
		curTm := time.Now().Format("20060102150405")
		order.PayId = fmt.Sprintf("%s%9d", curTm, order.Id)
		order.ModifyAt = curTm
		_, err = db.Update(&order, "pay_id", "modify_at")
		if nil != err {
			return nil, err
		}
	}

	return order, nil
}

func PrePayResp(db orm.Ormer, thirdId string, state int) error {
	order.ThirdId = thirdId
	order.State = state
	order.ModifyAt = time.Now().Format("20060102150405")

	_, err = db.Update(&order, "pay_id", "modify_at", "state")
	if nil != err {
		return nil, err
	}
}
