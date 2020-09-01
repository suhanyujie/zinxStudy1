package znet

import (
	"fmt"
	"log"
	"net"
	"zinx_study1/ziface"
)

//IServer 的接口实现
type Server struct {
	Name string
	IpVersion string
	Ip string
	Port int
}

// 启动
func (s *Server) Start()  {
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
	//阻塞等待客户端链接
	for {
		conn, err := listenner.Accept()
		if err != nil {
			log.Printf("Accpet connection error: %s\n", err)
			continue
		}
		// 对客户端连接进行一些读写操作
		// 基于 tcp 的连接是全双工的，可以收信息也能发信息
		for  {
			buffer := make([]byte, 512)
			cnt, err := conn.Read(buffer)
			if err != nil {
				log.Printf("read client data error: %s\n", err)
				conn.Close()
				break
			}
			// debug 打印一些信息
			log.Println(string(buffer))
			// 回写操作
			cnt, err = conn.Write(buffer[:cnt])
			if err != nil {
				log.Printf("write data into client error: %s\n", err)
				conn.Close()
				break
			}
		}
	}
}

// 停止
func (s *Server) Stop()  {
	// todo 进行将资源回收等操作
	return
}

// 服务
func (s *Server) Serve()  {
	s.Start()
	// 启动服务器后的一些操作
	// 阻塞
	select {}
}

// 实例化 server
func NewServer(name string) ziface.IServer {
	s := &Server {
		Name: name,
		IpVersion: "tcp4",
		Ip: "0.0.0.0",
		Port: 3001,
	}
	return s
}

