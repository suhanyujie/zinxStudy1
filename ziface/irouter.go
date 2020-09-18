package ziface

//路由抽象层接口
type IRouter interface {
	// before business hook
	PreHandle(request IRequest)
	// doing business hook
	DoingHandle(request IRequest)
	// after business hook
	AfterHandle(request IRequest)
}
