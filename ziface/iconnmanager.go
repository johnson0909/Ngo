package ziface

//定义连接管理接口

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, error) 
	Len() int				//获取当前连接个数
	ClearConn()				//删除并停止所有连接
}