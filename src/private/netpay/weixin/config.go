/*
	微信配置的相关功能
*/

package wx

import (
	"encoding/json"
)

type ConfigST struct {
	Scan ConfItemST `json:"Scan" yaml:"Scan"` // 扫码支付的配置信息
}

// 微信业务的配置信息
type ConfItemST struct {
	AppId       string `json:"AppId" yaml:"AppId"`             //微信分配的公众账号ID
	MchId       string `json:"MchId" yaml:"MchId"`             //用户在商户appid下的唯一标识
	SecretKey   string `json:"SecretKey" yaml:"SecretKey"`     //秘钥
	CallBackUrl string `json:"CallBackUrl" yaml:"CallBackUrl"` //通知回调接口
	ServerIp    string `json:"ServerIp" yaml:"ServerIp"`       //server Ip，微信接口需要带给微信
}

var (
	g_conf ConfigST
)

func Init(inConfText string) error {
	return g_conf.Deserialize(inConfText)
}

func (c *ConfigST) Serialize() string {
	return ""
}

func (c *ConfigST) Deserialize(inConfText string) error {
	return json.Unmarshal([]byte(inConfText), c)
}
