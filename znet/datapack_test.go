package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestPack(t *testing.T) {
	// 模拟一个服务器
	ls, err := net.Listen("tcp", "127.0.0.1:3002")
	if err != nil {
		t.Fatalf("listen error: %s\n", err)
	}
	go func() {
		for {
			conn, err := ls.Accept()
			if err != nil {
				t.Fatalf("accept error: %s\n", err)
			}
			go func() {
				dp := NewDataPacker()
				// 处理客户端请求
				for {
					headBuff := make([]byte, dp.GetHeaderLen())
					// 第 1 次从 conn 读，读出包的 header
					_, err := io.ReadFull(conn, headBuff)
					if err != nil {
						t.Logf("read headbuff error: %s\n", err)
						return
					}
					// 数据在服务端解包
					msgHead, err := dp.UnPack(headBuff)
					if err != nil {
						t.Logf("UnPack headbuff error: %s\n", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// 此时，data 是有数据的
						// 第 2 次读，根据 header 中的 dataLen，读出数据
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Logf("read data error: %s\n", err)
							return
						}
						// 一个完整的消息读取完毕
						fmt.Printf("---> Recv MsgId: %d len is: %d, data is: %v\n", msg.Id, msg.DataLen, string(msg.Data))
					}
				}
			}()
		}
	}()
	// 模拟一个客户端
	c1Conn, err := net.Dial("tcp", "127.0.0.1:3002")
	if err != nil {
		t.Logf("c1 connect error: %s\n", err)
		return
	}
	index := 1
	dp := NewDataPacker()
	for {
		// 数据在客户端封包
		// 客户端写入数据
		data := []byte("Hello world " + string(index))
		msg1 := &Message{
			Id:      1,
			DataLen: uint32(len(data)),
			Data:    data,
			Headers: nil,
			Body:    nil,
		}
		data1, err := dp.Pack(msg1)
		if err != nil {
			t.Logf("c1 pack data error: %s\n", err)
			return
		}
		data = []byte("一起学习 zinx")
		msg2 := &Message{
			Id:      1,
			DataLen: uint32(len(data)),
			Data:    data,
			Headers: nil,
			Body:    nil,
		}
		data2, err := dp.Pack(msg2)
		if err != nil {
			t.Logf("c1 pack data error: %s\n", err)
			return
		}
		// 模拟粘包
		dataBytes := append(data1, data2...)
		_, err = c1Conn.Write(dataBytes)
		if err != nil {
			t.Logf("c1 write data error: %s\n", err)
			return
		}
		// 客户端写读取数据
		buff := make([]byte, 512)
		_, err = c1Conn.Read(buff)
		if err != nil {
			t.Logf("c1 read data error: %s\n", err)
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func TestUnPack(t *testing.T) {

}
