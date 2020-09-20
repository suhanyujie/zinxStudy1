package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
	Headers map[string]string
	Body    []byte
}

func NewMessage(msgId uint32, data []byte) *Message {
	return &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
		Headers: nil,
		Body:    nil,
	}
}

func (this *Message) SetMsgId(msgId uint32) {
	this.Id = msgId
}

func (this *Message) GetMsgId() uint32 {
	return this.Id
}

func (this *Message) SetMsgLen(msgLen uint32) {
	this.DataLen = msgLen
}

func (this *Message) GetMsgLen() uint32 {
	return this.DataLen
}

func (this *Message) SetData(data []byte) {
	this.Data = data
}

func (this *Message) GetData() []byte {
	return this.Data
}
