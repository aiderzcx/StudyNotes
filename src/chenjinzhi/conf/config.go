package conf

import (
	"customlib/database/mysql/bgorm"
	"customlib/log"
	"customlib/tool"
	"fmt"
	"production/netpay/weixin"
)

type ConfigST struct {
	SrvId  string          `json:"SrvId" yaml:"SrvId"`
	Listen string          `json:"Listen" yaml:"Listen"`
	Log    log.LogConfigST `json:"Log" yaml:"Log"`
	Mysql  bgorm.ConfigST  `json:"Mysql" yaml:"Mysql"`
	Wx     string          `json:"Wx" yaml:"Wx"`
}

var (
	G_conf ConfigST
)

func Init(filePath string) error {
	if err := tool.ParseConfig(filePath, &G_conf); nil != err {
		return fmt.Errorf("ParseConfig(%s).%v\n", filePath, err)
	}

	fmt.Printf("config: %+v\n", G_conf)

	if err := log.Init(G_conf.Log); nil != err {
		return fmt.Errorf("log.Init(%+v).%v", G_conf.Log, err)
	}

	if err := wx.Init(G_conf.Wx); nil != err {
		return fmt.Errorf("wx.Init().%v]", err)
	}

	if err := bgorm.Init(G_conf.Mysql); nil != err {
		return fmt.Errorf("bgorm.Init().%v", err)
	}

	return nil
}
