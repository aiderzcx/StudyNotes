/*
	订单的管理
	1. 创建订单
	2. 获取订单
	3. 更新订单
*/

package dbm

import (
	"bytes"
	"customlib/tool"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(PayOrderTB))
}

const (
	PAY_STATE_NIL         = iota // 空状态
	PAY_STATE_SCANED             // 收到扫码支付
	PAY_STATE_PREPAYED           // 发起预支付，等待支付回调
	PAY_STATE_PREPAY_FAIL        // 预支付失败
	PAY_STATE_PAY_SUCC           // 支付成功
	PAY_STATE_PAY_FAIL           // 支付失败
	PAY_STATE_REFUNDING          // 正在退款
	PAY_STATE_REFUNDED           // 已经退款
	PAY_STATE_CLOSED             // 支付已经关闭
)

type PayOrderTB struct {
	Id          int64 `orm:"pk;auto"`
	ProductId   string
	ProductDesc string
	TotalFee    int64
	PayId       string
	PayType     string
	ThirdPreId  string // 第三方的预支付ID
	ThirdId     string // 第三方的支付ID
	State       int
	Remarks     string // 支付失败的原因
	CreateAt    string // 订单创建时间
	PayAt       string // 支付时间
	RefundAt    string // 退款时间
	QueryCount  int    // 查询的次数
}

func (p *PayOrderTB) TableName() string {
	return "pay_order"
}

func CreateOrder(db orm.Ormer, productId, productDesc string, totalFee int64) (*PayOrderTB, error) {
	var order PayOrderTB

	if productId != "" {
		err := db.QueryTable(&order).Filter("product_id", productId).One(&order)
		if err != nil && err != orm.ErrNoRows {
			return nil, err
		}

		if err == nil {
			if order.TotalFee != totalFee {
				return nil, fmt.Errorf("totalFee(%d) != inPut(%d)", order.TotalFee, totalFee)
			}

			return &order, nil
		}
	}

	order.ProductId = productId
	order.ProductDesc = productDesc
	order.TotalFee = totalFee
	order.CreateAt = tool.CurTimeNormal()

	lastId, err := db.Insert(&order)
	if nil != err {
		return nil, err
	}

	order.Id = lastId
	return &order, nil
}

// 根据支付ID获取订单
func OrderInfo(db orm.Ormer, payId string) (*PayOrderTB, error) {
	var order PayOrderTB

	err := db.QueryTable(&order).Filter("pay_id", payId).One(&order)
	if nil != err {
		return nil, err
	}

	return &order, nil
}

func StartPrepay(db orm.Ormer, inProductId string) (*PayOrderTB, error) {
	var order PayOrderTB

	err := db.QueryTable(&order).Filter("product_id", inProductId).One(&order)
	if nil != err {
		return nil, err
	}

	// 没有订单，更新订单
	if order.PayId == "" {
		order.PayId = fmt.Sprintf("%s%09d", transDbTm2Wx(order.CreateAt), order.Id)
		_, err = db.Update(&order, "pay_id")
		if nil != err {
			return nil, err
		}
	}

	return &order, nil
}

func PrePayResp(db orm.Ormer, order *PayOrderTB) error {
	_, err := db.Update(order, "pay_id", "pay_type", "third_pre_id", "state", "remarks")
	return err
}

func PayCallback(db orm.Ormer, inOrder *PayOrderTB) error {
	_, err := db.Update(inOrder, "third_id", "state", "pay_at", "remarks")
	return err
}

func UnPayedOrder(db orm.Ormer, startId int64, startTm string) ([]PayOrderTB, error) {
	var orders []PayOrderTB
	var err error

	if startId > 0 {
		_, err = db.QueryTable(&PayOrderTB{}).
			Filter("id__gte", startId).
			Filter("state", PAY_STATE_PREPAYED).
			Limit(100).
			All(&orders)
	} else {
		_, err = db.QueryTable(&PayOrderTB{}).
			Filter("create_at__gte", startTm).
			Filter("state", PAY_STATE_PREPAYED).
			Limit(100).
			All(&orders)
	}

	return orders, err
}

func UpdateQueryOrder(db orm.Ormer, orders *PayOrderTB) error {
	return nil
}

const (
	std_tm_len = len("2006-01-02 15:04:05")
)

func transDbTm2Wx(inTm string) string {
	if len(inTm) < std_tm_len {
		return time.Now().Format("20060102150405")
	}

	var buf bytes.Buffer

	buf.WriteString(inTm[0:4])   // yyyy
	buf.WriteString(inTm[5:7])   // MM
	buf.WriteString(inTm[8:10])  // DD
	buf.WriteString(inTm[11:13]) // HH
	buf.WriteString(inTm[14:16]) // MM
	buf.WriteString(inTm[17:])   // SS

	return buf.String()
}
