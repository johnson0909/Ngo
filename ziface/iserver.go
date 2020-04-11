package ziface
type IServer interface {
	Start()
	Stop()
	Serve()
	//路由功能，给当前服务注册一个路由方法，提供给客户端连接处理使用
	AddRouter(msgId uint32, router IRouter)
	//得到连接管理器
	GetConnMgr() IConnManager
	//设置server的链接创建时的hook函数
	SetOnConnStart(func(IConnection))
	//设置server的链接断开时的hook函数
	SetOnConnStop(func(IConnection))
	//连接创建时的hook函数
	CallOnConnStart(conn IConnection)
	//连接断开时的hook函数
	CallOnConnStop(conn IConnection)
}