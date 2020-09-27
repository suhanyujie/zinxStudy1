package ziface

type IMsgHandler interface {
	// 新增路由
	AddRouter(flag uint32, router IRouter)
	// 查找 router 并执行
	DoHandler(request IRequest)
	// 创建 work pool
	StartWorkPool(size uint32)

	StartOneWork(workId uint32, taskQueue chan IRequest)

	SendMsgToTaskQueue(request IRequest)
}
