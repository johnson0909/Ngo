package main
import (
	"fmt"
	"Ngo/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client0 Test ... start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client0 start err, exit!")
		return
	}

	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("client0 test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client0 write err: ", err)
			return
		}
		//先读流出中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("clent0 unpack head err: ", err)
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client0 unpack head err: ", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client0 unpack data err")
				return
			}
			fmt.Printf("==> Client0 receive Msg: Id = %d, len = %d , data = %s\n", msg.Id, msg.DataLen, msg.Data)
		}
		time.Sleep(1 * time.Second)
	}
}