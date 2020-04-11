package znet

import (
	"fmt"
	"Ngo/utils"
	"Ngo/ziface"
	"strconv"
)

type MsgHandle struct {
	Apis	map[uint32]ziface.IRouter  //存放每个msgId对应的业务处理方法
	WorkerPoolSize uint32	//业务工作Worker池的数量
	TaskQueue	[]chan ziface.IRequest //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:	make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个Worker对应一个消息队列
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//将消息交给消息队列，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//使用RR轮询分配，根据ConnID来分配哪个worker负责处理
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request msgID = ", request.GetMsgID(), " to workerID = ", workerID)
	mh.TaskQueue[workerID] <- request
}

//马上以非阻塞的方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not Found！")
		return
	}

	//执行对应的处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1.判断当前msg绑定的API方法是否存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}

	//添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = " ,msgId)
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	for {
		select {
		case request := <- taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要worker的数量
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
