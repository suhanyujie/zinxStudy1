package znet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"zinx_study1/utils"
	"zinx_study1/ziface"
)

type Connection struct {
	// 当前连接 属于哪个 server
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnId    uint32
	IsClosed  bool
	// handleAPI ziface.HandleFunc
	ExitChan chan bool
	// router
	// router 和 handleAPI 是二选一的，作用是类似的
	// Router ziface.IRouter
	// 集成消息路由管理
	MsgHandler ziface.IMsgHandler
	// 无缓冲通道，goroutine 之间的消息通信
	msgChan chan []byte
	// 为了给开发者提供更大的灵活性，遂给每个连接定义一个属性集合，开发者使用时，可以在对应的链接中设定属性、删除属性
	property map[string]interface{}
	// 因为属性需要读写操作，为了数据安全性，加上一个锁进行保护
	propertyLock sync.RWMutex
}

// 实例化自定义的链接
func NewConnection(server ziface.IServer, conn *net.TCPConn, cid uint32, msgHandler ziface.IMsgHandler) *Connection {
	userConn := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnId:       cid,
		IsClosed:     false,
		ExitChan:     make(chan bool, 1),
		MsgHandler:   msgHandler,
		msgChan:      make(chan []byte),
		property:     make(map[string]interface{}),
		propertyLock: sync.RWMutex{},
	}
	server.GetConnManager().Add(userConn)
	return userConn
}

// 开始处理连接
func (c *Connection) Start() {
	fmt.Println("new connection connected", c.ConnId)
	// 调用设定的"连接开始 hook"
	c.TcpServer.CallOnConnStartFunc(c)
	// 启动读数据逻辑
	go c.StartReader()
	// 启动写数据逻辑
	go c.StartWriter()
}

// 客户端链接的读逻辑
func (c *Connection) StartReader() {
	fmt.Println("start reader goroutine")
	defer c.Stop()
	for true {
		// 1.实例化封包解包器
		dp := NewDataPacker()
		headBuff := make([]byte, dp.GetHeaderLen())
		// io.LimitReader()
		// 2.从连接中读取 msg head 信息
		_, err := io.ReadFull(c.Conn, headBuff)
		if err != nil {
			log.Printf("server receive head buf err: %s\n", err)
			// 连接出问题了，需关闭连接，这里如何主动关闭连接？
			c.Stop()
			break
		}
		// 3.将 msg head 解包
		msg, err := dp.UnPack(headBuff)
		if err != nil {
			log.Printf("server unpack head buf err: %s\n", err)
			// 连接出问题了，需关闭连接，这里如何主动关闭连接？
			c.Stop()
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
				// 连接出问题了，需关闭连接，这里如何主动关闭连接？
				c.Stop()
				break
			}
			if err != nil {
				log.Printf("server unpack data buf err: %s\n", err)
				// 连接出问题了，需关闭连接，这里如何主动关闭连接？
				c.Stop()
				break
			}
		}
		msg.SetData(dataBuff)

		// 根据读取到的数据，封装成 request
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkPoolSize > 0 {
			// 已经配置了工作池机制
			// 将请求消息发送给 work 队列
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 寻找对应的处理 handler 并执行
			go c.MsgHandler.DoHandler(&req)
		}
	}
}

// 客户端链接的写逻辑
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit!]")
	// 阻塞等待 channel 消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Printf("writer goroutine send data error: %s\n", err)
				return
			}
		case <-c.ExitChan:
			// 代表 reader 退出
			return
		}
	}
}

// 关闭服务时的处理
func (c *Connection) Stop() {
	fmt.Println("conn stop ConnID=", c.ConnId)
	if c.IsClosed == true {
		return
	}
	c.IsClosed = true
	// 告诉 writer，连接关闭
	c.ExitChan <- true
	// 关闭 socket 链接
	c.Conn.Close()
	// 回收管道资源
	close(c.ExitChan)
	// 将当前连接从连接管理器中移除
	c.TcpServer.GetConnManager().Remove(c)
	// 调用设定的"关闭连接 hook"
	c.TcpServer.CallOnConnStopFunc(c)
}

// 获取 tcp 的连接对象
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

// 获取链接的id
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 返回远程的连接地址
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
	// _, err = c.GetTcpConnection().Write(binaryMsg)
	c.msgChan <- binaryMsg
	return err
}

// 向连接中存储属性
func (_this *Connection) SetProperty(key string, value interface{}) {
	_this.propertyLock.Lock()
	defer _this.propertyLock.Unlock()
	_this.property[key] = value
	if val, ok := _this.property[key]; ok {
		log.Printf("set property ok: %v\n", val)
	} else {
		log.Printf("set property not ok\n")
	}
}

// 获取属性值
func (_this *Connection) GetProperty(key string) (interface{}, error) {
	_this.propertyLock.Lock()
	defer _this.propertyLock.Unlock()
	if value, ok := _this.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("value not found. ")
}
