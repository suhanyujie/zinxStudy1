package znet

import (
	"errors"
	"log"
	"sync"
	"zinx_study1/ziface"
)

/// 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	// 对连接集合管理时的读写锁
	connLock sync.RWMutex
}

// 实例化连接管理器
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
}

// 新增连接
func (this *ConnManager) Add(conn ziface.IConnection) {
	// 写入时，加锁
	this.connLock.Lock()
	defer this.connLock.Unlock()
	this.connections[conn.GetConnId()] = conn
	log.Printf("add conn into connManager ok, manager len is: %d\n", this.Len())
}

// 删除连接
func (this *ConnManager) Remove(conn ziface.IConnection) {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	delete(this.connections, conn.GetConnId())
	log.Printf("Remove conn from connManager ok, connId is: %d\n", conn.GetConnId())
}

// 根据 connId 获取连接
func (this *ConnManager) GetConnByConnId(connId uint32) (ziface.IConnection, error) {
	if conn, ok := this.connections[connId]; ok {
		return conn, nil
	}

	return nil, errors.New("Could not found this conn! ")
}

// 获取连接总数
func (this *ConnManager) Len() int {
	return len(this.connections)
}

// 清除并终止所有连接
func (this *ConnManager) ClearConn() {
	this.connLock.Lock()
	defer this.connLock.Unlock()
	for connId, conn := range this.connections {
		conn.Stop()
		delete(this.connections, connId)
	}
	log.Printf("Clear all connections! ")
}
