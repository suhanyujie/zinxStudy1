package znet

import "zinx_study1/ziface"

// 实现一个 base router
// 后续在业务中使用时，可以直接嵌入 BaseRouter，再根据实际需要重写对应的方法
type BaseRouter struct{}

// 一些业务中，有些钩子不需要，因此，在 base router 中，方法的不具体实现，而是空方法体
// before business handler
func (_this *BaseRouter) PreHandle(request ziface.IRequest) {

}

// doing business handler
func (_this *BaseRouter) DoingHandle(request ziface.IRequest) {

}

// after business handler
func (_this *BaseRouter) AfterHandle(request ziface.IRequest) {

}
