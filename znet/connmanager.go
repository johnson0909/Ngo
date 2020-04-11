package znet

import (
	"errors"
	"fmt"
	"Ngo/ziface"
	"sync"
)
/*
	连接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex 
}

//创建一个连接管理
func NewConnManager() *ConnManager {
	return &ConnManager {
		connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map,加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到connManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager succ! conn num = ", connMgr.Len())

}

//删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection Remove ConnID = ", conn.GetConnID(), " succ! conn num = ", connMgr.Len())


}

//利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接管理器中连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

//消除并停止所有服务
func (connMgr *ConnManager) ClearConn() {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除所有连接
	for connID, conn := range connMgr.connections {
		conn.Stop()

		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All Connections succ!!! conn num = ", connMgr.Len())
}