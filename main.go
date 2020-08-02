package main

import (
	"github.com/QingGo/udp_quic_proxy/config"
	"github.com/QingGo/udp_quic_proxy/proxy"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("配置信息：%+v\n", config.GConfig)
	localProxy := proxy.NewLocalProxy()
	localProxy.Run()
}
