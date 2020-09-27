package znet

import (
	"log"
	"zinx_study1/utils"
	"zinx_study1/ziface"
)

type MsgHandler struct {
	// 路由存储
	Apis map[uint32]ziface.IRouter
	// 添加 work 协程池属性
	WorkPoolSize uint32
	// 取任务的队列
	TaskQueue []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utils.GlobalObject.WorkPoolSize, // 通过配置获取
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.MaxTaskLenForOneWork),
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

// 实例化 work pool
func (this *MsgHandler) StartWorkPool(size uint32) {
	for i := uint32(0); i < size; i++ {
		this.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxTaskLenForOneWork)
		go this.StartOneWork(i, this.TaskQueue[i])
	}
}

// 实现基于 work pool 的工作流
func (this *MsgHandler) StartOneWork(workId uint32, taskQueue chan ziface.IRequest) {
	log.Printf("Work [%d] is start\n", workId)
	// 不断阻塞等待对应消息队列中的消息
	for {
		select {
		case request := <-taskQueue:
			this.DoHandler(request)
		}
	}
}

// 将请求消息发送给某一个 work 队列
func (_this *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 均衡地分配给这些 work
	// 可以根据客户端链接的 ConnID 来分配
	oneWorkId := request.GetConnection().GetConnId() % _this.WorkPoolSize
	log.Printf("send request: %d to work id: %d\n", request.GetConnection().GetConnId(), oneWorkId)
	_this.TaskQueue[oneWorkId] <- request
}
