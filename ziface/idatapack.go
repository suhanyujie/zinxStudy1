package ziface

// 封包，拆包模块
type IDataPack interface {
	// 获取包长度
	GetHeaderLen()
	// 封包
	Pack(msg IMessage) ([]byte, error)
	// 拆包
	UnPack([]byte) (IMessage, error)
}
