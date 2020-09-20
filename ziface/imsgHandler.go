package ziface

type IMsgHandler interface {
	// 新增路由
	AddRouter(flag uint32, router IRouter)
	// 查找 router 并执行
	DoHandler(request IRequest)
}
