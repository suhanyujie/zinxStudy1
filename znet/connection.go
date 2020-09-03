package znet

import (
	"fmt"
	"log"
	"net"
	"zinx_study1/ziface"
)

type Connection struct {
	Conn      *net.TCPConn
	ConnId    uint32
	IsClosed  bool
	handleAPI ziface.HandleFunc
	ExitChan  chan bool
}

// 实例化自定义的链接
func NewConnection(conn *net.TCPConn, cid uint32, callback ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnId:    cid,
		IsClosed:  false,
		handleAPI: callback,
		ExitChan:  make(chan bool, 1),
	}
}

// todo
func (c *Connection) Start() {
	fmt.Println("new connection connected", c.ConnId)
	// 启动读数据逻辑
	go c.StartReader()
	// 启动写数据逻辑
	// go c.StartWriter()
}

// 客户端链接的读逻辑
func (c *Connection) StartReader() {
	fmt.Println("start reader goroutine")
	defer c.Stop()
	for true {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			log.Printf("server receive buf err: %s\n", err)
			// 连接出问题了，需关闭连接，这里如何主动关闭连接？
			c.Stop()
			break
		}
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			log.Printf("handle api err: %s\n", err)
		}
	}
}

// 客户端链接的写逻辑
func (c *Connection) StartWriter() {
	fmt.Println("start writer goroutine")
}

// todo
func (c *Connection) Stop() {
	fmt.Println("conn stop ConnID=", c.ConnId)
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	// 关闭 socket 链接
	c.Conn.Close()
	// 回收管道资源
	close(c.ExitChan)
}

// todo
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

// todo
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// todo
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// todo
func (c *Connection) Send(data []byte) error {

	return nil
}
