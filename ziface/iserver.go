package ziface

//定义一个服务器接口
type IServer interface {
	Start()
	Stop()
	Serve()

	// func for router
	// regist a route in this server
	AddRoute(router IRouter)
}
