package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTcpConnection() *net.TCPConn
	GetConnId() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
	// 向连接中存储一些属性值
	SetProperty(key string, value interface{})
	// 获取连接中存储的属性值
	GetProperty(key string) (interface{}, error)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
