package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ... start")
	time.Sleep(3*time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Client start err, exit!")
		return 
	}

	for {
		_, err := conn.Write([]byte("hello Ngo"))
		if err != nil {
			fmt.Println("send buf err: ", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err: ", err)
			return 
		}
		fmt.Printf("server call back: %s, cnt: %d\n", buf, cnt)
		time.Sleep(1*time.Second)
	}

}

//单元测试
func main(t *testing.T) {
	s := NewServer("Ngo v0.1")
	go ClientTest()
	s.Serve()  
}