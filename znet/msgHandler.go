package znet

import (
	"log"
	"zinx_study1/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 新增路由
func (this *MsgHandler) AddRouter(flag uint32, router ziface.IRouter) {
	if _, ok := this.Apis[flag]; ok {
		return
	}
	this.Apis[flag] = router
	log.Printf("add route success!\n")
}

// 通过 path 查找 router
func (this *MsgHandler) DoHandler(request ziface.IRequest) {
	if router, ok := this.Apis[request.GetMsgId()]; ok {
		router.PreHandle(request)
		router.DoingHandle(request)
		router.AfterHandle(request)
		return
	}
	log.Printf("couldn't found the router for the msgId: %d\n", request.GetMsgId())
}
