package znet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"zinx_study1/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnId   uint32
	IsClosed bool
	// handleAPI ziface.HandleFunc
	ExitChan chan bool
	// router
	// router 和 handleAPI 是二选一的，作用是类似的
	Router ziface.IRouter
}

// 实例化自定义的链接
func NewConnection(conn *net.TCPConn, cid uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnId:   cid,
		Router:   router,
		IsClosed: false,
		ExitChan: make(chan bool, 1),
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
		//buf := make([]byte, utils.GlobalObject.MaxPkgSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	log.Printf("server receive buf err: %s\n", err)
		//	// 连接出问题了，需关闭连接，这里如何主动关闭连接？
		//	c.Stop()
		//	break
		//}
		// 1.实例化封包解包器
		dp := NewDataPacker()
		headBuff := make([]byte, dp.GetHeaderLen())
		// io.LimitReader()
		// 2.从连接中读取 msg head 信息
		_, err := io.ReadFull(c.Conn, headBuff)
		if err != nil {
			log.Printf("server receive head buf err: %s\n", err)
			break
		}
		// 3.将 msg head 解包
		msg, err := dp.UnPack(headBuff)
		if err != nil {
			log.Printf("server unpack head buf err: %s\n", err)
			break
		}
		headMsgObj := msg.(*Message)
		var dataBuff []byte
		if headMsgObj.GetMsgLen() > 0 {
			// 根据 msg head 中的信息，读取数据
			dataBuff = make([]byte, headMsgObj.DataLen)
			_, err = io.ReadFull(c.Conn, dataBuff)
			if err != nil {
				log.Printf("server receive data buf err: %s\n", err)
				break
			}
			if err != nil {
				log.Printf("server unpack data buf err: %s\n", err)
				break
			}
		}
		msg.SetData(dataBuff)

		// 根据读取到的数据，封装成 request
		req := Request{
			conn: c,
			msg:  msg,
		}
		go func(request ziface.IRequest) {
			// pre handle of route
			c.Router.PreHandle(request)
			c.Router.DoingHandle(request)
			c.Router.AfterHandle(request)
		}(&req)
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

// 封包，然后将数据写到客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("conn is closed! ")
	}
	dp := NewDataPacker()
	msg := NewMessage(msgId, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		return err
	}
	_, err = c.GetTcpConnection().Write(binaryMsg)
	return err
}
