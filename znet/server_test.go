package znet

import (
	"Ngo/ziface"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func ClientTest(i uint32) {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client start err, exit!")
		return
	}

	for {
		dp := NewDataPack()
		msg, _ := dp.Pack(NewMsgPackage(i, []byte("client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client write err: ", err)
			return
		}
		//先读流出中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("clent unpack head err: ", err)
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack head err: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client unpack data err")
				return
			}
			fmt.Printf("==> Client receive Msg: Id = %d, len = %d , data = %s\n", msg.Id, msg.DataLen, msg.Data)
		}
		time.Sleep(1 * time.Second)
	}

}

type PingRouter struct {
	BaseRouter
}

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	err := request.GetConnection().SendMsg(1, []byte("before ping ...\n"))
	if err != nil {
		fmt.Println("preHandle SendMsg err: ", err)
	}
}

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Println("recv from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("Handle SendMsg err: ", err)
	}
}

func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	err := request.GetConnection().SendMsg(1, []byte("after ping ...\n"))
	if err != nil {
		fmt.Println("Post SendMsg err: ", err)
	}
}

type HelloRouter struct {
	BaseRouter
}

func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call helloRouter Handle")
	fmt.Println("recevie from client: msgId = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(2, []byte("hello Ngo hello router\n"))
	if err != nil {
		fmt.Println("HelloRouter Handle SendMsg err: ", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnection is Called ...")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called ...")
}

//单元测试
func TestServer(t *testing.T) {
	//创建一个server句柄
	s := NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//多路由
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	//客户端测试
	go ClientTest(1)
	go ClientTest(2)

	//开启服务
	go s.Serve()

	select {
	case <-time.After(time.Second * 10):
		return
	}
}
