package ziface

type IMessage interface {
	SetMsgId(msgId uint32)
	GetMsgId() uint32
	SetMsgLen(msgLen uint32)
	GetMsgLen() uint32
	SetData([]byte)
	GetData() []byte
}
