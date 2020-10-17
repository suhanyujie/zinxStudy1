package znet

import (
	"fmt"
	"log"
	"net"
	"zinx_study1/utils"
	"zinx_study1/ziface"
)

//IServer 的接口实现
type Server struct {
	Name        string
	IpVersion   string
	Ip          string
	Port        int
	MsgHandler  ziface.IMsgHandler
	ConnManager ziface.IConnManager
}

// 对客户端的业务处理 暂时固定，后续优化
func CallbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	var err error
	log.Printf("[Conn Handle] CallbackToClient...\n")
	// 回写操作
	cnt, err = conn.Write(data)
	if err != nil {
		log.Printf("write data into client error: %s\n", err)
		conn.Close()
	}
	return err
}

// 启动
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s, is starting\n", s.Ip)
	// 获取一个 tcp 的addr
	addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		log.Fatalf("resolve tcp addr error: %s\n", err)
		return
	}
	//监听服务器的地址
	listenner, err := net.ListenTCP(s.IpVersion, addr)
	if err != nil {
		log.Fatalf("listen tcp error: %s\n", err)
	}
	// 设定一个连接id
	var cid uint32 = 1
	//阻塞等待客户端链接
	for {
		conn, err := listenner.AcceptTCP()
		if err != nil {
			log.Printf("Accpet connection error: %s\n", err)
			continue
		}
		// 当前连接数是否过大
		if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
			// todo 给客户端返回一个，服务连接数过大的提示
			log.Printf("Too many connection. Can't deal with it. \n")
			conn.Close()
			continue
		}
		userConn := NewConnection(s, conn, cid, s.MsgHandler)
		cid++
		// 开始处理当前请求的业务
		go userConn.Start()
	}
}

// 停止
func (s *Server) Stop() {
	// 进行资源回收等操作
	s.ConnManager.ClearConn()
	log.Printf("server stop, clear all the connections.\n")
	return
}

// 服务
func (s *Server) Serve() {
	s.Start()
	// 启动服务器后的一些操作
	// 阻塞
	select {}
}

// 实例化 server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        name,
		IpVersion:   "tcp4",
		Ip:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	s.MsgHandler.StartWorkPool(utils.GlobalObject.WorkPoolSize)
	return s
}

func (s *Server) AddRoute(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

// 获取连接管理器
func (_this *Server) GetConnManager() ziface.IConnManager {
	return _this.ConnManager
}
