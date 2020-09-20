package znet

import "zinx_study1/ziface"

// request 消息/请求封装
type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

// 获取请求的连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取消息的id
func (_this *Request) GetMsgId() uint32 {
	return _this.msg.GetMsgId()
}
