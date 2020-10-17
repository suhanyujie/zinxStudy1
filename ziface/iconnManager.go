package ziface

/// 连接管理模块抽象
type IConnManager interface {
	// 新增连接
	Add(conn IConnection)
	// 删除连接
	Remove(conn IConnection)
	// 根据 connId 获取连接
	GetConnByConnId(connId uint32) (IConnection, error)
	// 获取连接总数
	Len() int
	// 清除并终止所有连接
	ClearConn()
}
