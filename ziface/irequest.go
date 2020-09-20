package ziface

// 将用户的请求包装到 Request 中
type IRequest interface {
	// 获取连接
	GetConnection() IConnection
	// 获取用户请求的数据
	GetData() []byte
	// 获取消息数据的id
	GetMsgId() uint32
}
