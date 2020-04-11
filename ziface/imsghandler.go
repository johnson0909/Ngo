package ziface

//定义消息管理接口
type IMsgHandle interface {
	DoMsgHandler(request IRequest)   //马上以非阻塞的方式处理消息
	AddRouter(msgId uint32, rouoter IRouter) //为消息添加具体的处理逻辑
	StartWorkerPool()      //启动worker工作池
	SendMsgToTaskQueue(request IRequest)  //将消息交给消息队列，由worker进行处理
}