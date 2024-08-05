package socket

import "net"

// Server is a socket server.
type Server struct {
	Name         string
	SendChanSize int          //发型消息缓冲区大小
	Listener     net.Listener // 服务器监听器
}
