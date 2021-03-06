package main

import (
	"fmt"
	"Ngo/ziface"
	"Ngo/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Println("recv from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println("Handle SendMsg err: ", err)
	}
}


type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call helloRouter Handle")
	fmt.Println("recevie from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("hello Ngo hello router"))
	if err != nil {
		fmt.Println("HelloRouter Handle SendMsg err: ", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnection is Called ...")

	//设置两个链接属性
	conn.SetProperty("Name", "saulliu")
	conn.SetProperty("Career", "后台开发")
	fmt.Println("Set conn Name, Career done!")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前查询连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if job, err := conn.GetProperty("Career"); err == nil {
		fmt.Println("Conn Property Career = ", job)
	}
	fmt.Println("DoConnectionLost is Called ...")
}

func main() {
	s := znet.NewServer()

	//设置hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	//开启服务
	s.Serve()
}