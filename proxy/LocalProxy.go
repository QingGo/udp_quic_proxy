package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/QingGo/udp_quic_proxy/config"
	quic "github.com/lucas-clemente/quic-go"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

// NewLocalProxy 生成一个LocalProxy
func NewLocalProxy() *LocalProxy {
	localUDPAddr := net.UDPAddr{
		Port: config.GConfig.ListenUDPPort,
		IP:   net.ParseIP(config.GConfig.ListenUDPAddress),
	}
	localUDPSocket, err := net.ListenUDP("udp", &localUDPAddr)
	if err != nil {
		log.Debugf("Fatal error: %s", err.Error())
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		// 这个应该和服务端一致
		NextProtos: []string{"quic-proxy"},
	}

	remoteAddr := fmt.Sprintf("%s:%d", config.GConfig.RemoteQUICAddress, config.GConfig.RemoteQUICPort)
	session, err := quic.DialAddr(remoteAddr, tlsConf, nil)
	if err != nil {
		// 这里也许该阻塞并加入重试逻辑等待连接远程quic服务器成功
		log.Debug(err)
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Debug(err)
	}

	dataChannelSend := make(chan []byte, 100)
	dataChannelReceive := make(chan []byte, 100)
	localProxy := LocalProxy{
		udpSocket:          localUDPSocket,
		quicStream:         stream,
		dataChannelSend:    dataChannelSend,
		dataChannelReceive: dataChannelReceive,
	}
	return &localProxy
}

// LocalProxy 接收本地udp包并通过quic包发送到指定的ip端口
type LocalProxy struct {
	udpSocket          *net.UDPConn
	quicStream         quic.Stream
	dataChannelSend    chan []byte
	dataChannelReceive chan []byte
	// 当第一次收到客户端信息时获取
	clientAddr *net.UDPAddr
}

// receiveFromLocal 从本地udp端口读数据
func (localProxy *LocalProxy) receiveFromLocal() {
	defer wg.Done()

	buf := make([]byte, 32768)
	for {
		// 阻塞
		n, clientAddr, err := localProxy.udpSocket.ReadFromUDP(buf)
		if err != nil {
			log.Warnf("ReceiveFromClient whih %s goroutine stop: %s", clientAddr, err.Error())
			break
		}
		if localProxy.clientAddr == nil {
			localProxy.clientAddr = clientAddr
		}
		newbuf := make([]byte, n)
		copy(newbuf, buf[:n])
		localProxy.dataChannelSend <- newbuf
		log.Debugf("receive local udp message %s from %s", newbuf, clientAddr)
	}
}

func (localProxy *LocalProxy) sendToServer() {
	defer wg.Done()

	for {
		select {
		case msg := <-localProxy.dataChannelSend:
			_, err := localProxy.quicStream.Write(msg)
			if err != nil {
				log.Warn(err)
			}
		}
	}
}

func (localProxy *LocalProxy) receiveFromServer() {
	defer wg.Done()

	buf := make([]byte, 32768)
	for {
		// 这里要注意，quic是否会像tcp一样把一个包拆成几部分
		n, err := localProxy.quicStream.Read(buf)
		if err != nil {
			log.Warn(err)
		}

		newbuf := make([]byte, n)
		copy(newbuf, buf[:n])
		localProxy.dataChannelReceive <- newbuf
		log.Debugf("receive remote quic message %s", newbuf)

	}
}

func (localProxy *LocalProxy) sendToLocal() {
	defer wg.Done()

	for {
		select {
		case msg := <-localProxy.dataChannelReceive:
			_, err := localProxy.udpSocket.WriteToUDP(msg, localProxy.clientAddr)
			if err != nil {
				log.Warn(err)
			}
		}
	}
}

// Run 启动local proxy的入口
func (localProxy *LocalProxy) Run() {
	wg.Add(4)
	go localProxy.receiveFromLocal()
	go localProxy.sendToServer()
	go localProxy.receiveFromServer()
	go localProxy.sendToLocal()
	log.Info("local proxy start")
	wg.Wait()
}
