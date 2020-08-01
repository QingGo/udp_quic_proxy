package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/QingGo/udp_quic_proxy/config"
)

func main() {
	log.Infof("配置信息：%+v\n", config.GetConfig())
}

