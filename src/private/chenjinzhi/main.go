package main

import (
	"customlib/mservice/gincom"
	"flag"
	"fmt"
	"production/chenjinzhi/conf"
	"production/chenjinzhi/router"
)

var (
	conf_path = flag.String("c", "./conf/config.json", "the config file path")
)

func main() {
	flag.Parse()
	fmt.Printf("conf_path: %s\n", *conf_path)

	if err := conf.Init(*conf_path); nil != err {
		fmt.Printf("conf.Init().%v]n", err)
		return
	}

	engine := gincom.GinServer(conf.G_conf.SrvId)
	router.Init(engine)

	if err := engine.Run(conf.G_conf.Listen); nil != err {
		fmt.Printf("engine.Run(%s).%v\n", conf.G_conf.Listen, err)
	}
}
