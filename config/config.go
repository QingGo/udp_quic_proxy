package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var gConfig Config

func init() {
	// log.DebugLevel = 5, log.InfoLevel = 4
	logLevel := flag.String("LogLevel", "info", "-LogLevel [debug|info|warn...]")
	listenUDPAddress := flag.String("ListenUDPAddress", "0.0.0.0", "-ListenUDPAddress 0.0.0.0")
	listenUDPPort := flag.Int("ListenUDPPort", 11114, "-ListenUDPPort 11114")
	listenQUICAddress := flag.String("ListenQUICAddress", "0.0.0.0", "-ListenQUICAddress 0.0.0.0")
	listenQUICPort := flag.Int("ListenQUICPort", 51444, "-ListenQUICPort 0.0.0.0")
	flag.Parse()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	logLevelParsed, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal("日志等级参数解析错误，请使用[debug|info|warn...]，退出程序。")
	}
	log.SetLevel(logLevelParsed)

	gConfig = Config{
		ListenUDPAddress:  *listenUDPAddress,
		ListenUDPPort:     *listenUDPPort,
		ListenQUICAddress: *listenQUICAddress,
		ListenQUICPort:    *listenQUICPort,
	}
}

// Config 用来维护传入的命令行参数
type Config struct {
	ListenUDPAddress  string
	ListenUDPPort     int
	ListenQUICAddress string
	ListenQUICPort    int
}

// GetConfig 获取gConfig单例
func GetConfig() Config {
	return gConfig
}
