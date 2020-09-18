package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"zinx_study1/utils"
	"zinx_study1/ziface"
)

type DataPack struct {
}

// 获取消息长度和消息id
func (_this *DataPack) GetHeaderLen() uint32 {
	// 长度（4字节） + ID（4字节）
	return 8
}

// 数据的打包
func (_this *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	// 将长度信息写入 buffer
	// 以二进制的方式写入，注意大小端字节序
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return []byte{}, err
	}
	// 将信息 id 写入 buffer
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return []byte{}, err
	}
	// 将信息主体写入 buffer
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// 数据的解包
func (_this *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	// 拆解数据
	// 只解析出数据 head 信息，得到 dataLen 和 msgId
	msg := &Message{}
	// 读取长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读取 msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	log.Printf("%v\n", msg)
	// 判断包长度超出允许的最大值
	if utils.GlobalObject.MaxPkgSize > 0 && msg.DataLen > utils.GlobalObject.MaxPkgSize {
		return nil, errors.New("package too large!")
	}
	return msg, nil
}

// 封包、拆包器
func NewDataPacker() *DataPack {
	return &DataPack{}
}
