## 项目简介
在udp协议和quic协议之间做一层代理，能够把收到的udp协议的数据包通过quic协议转发出去，反之亦可。设想的用途是可以让现有的基于udp连接的网络程序在不用修改源码的前提下，就能通过使用quic协议获得网络传输上的可靠性。

架构设想如下：
```
origin udp client <-> local quic proxy <- 不可靠的网络环境 -> server quic proxy <->  origin udp server
```
* origin udp client和origin udp server即现有的基于udp连接的网络程序。
* local quic proxy和origin udp client应该在同一机器或同一内网环境。作用是把从origin udp client收到的udp协议数据包通过quic协议转发出去。而且从把server quic proxy传回来的quic协议数据包转换成udp协议数据包传回origin udp client。
* server quic proxy应该和origin udp server在同一机器或同一内网环境。作用是把local quic proxy传过来的quic协议数据包转换成udp协议数据包传给origin udp server。而且从把origin udp server传回来的udp协议数据包转换成quic协议数据包传给local quic proxy。
* 本程序包含local quic proxy和server quic proxy两个功能。

## 运行方法
### 手动测试
go run tool/quic_echo_server/quic_echo_server.go 
go run main.go -LogLevel=debug
nc -u 127.0.0.1 11114