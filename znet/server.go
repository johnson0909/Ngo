package znet
import (
	"fmt"
	"net"
	"Ngo/ziface"
)
type Server struct {
	Name string
	IPVersion string
	IP string
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server listen on IP: %s, Port: %d is starting.\n", s.IP, s.Port);

	//开启一个协程等待连接
	go func() {
		//获取一个TCP地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error, err: ", err)
			return 
		}
		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen err: ", err)
			return
		}
		fmt.Println("start Ngo server ", s.Name, " succ! now listenning...")
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}
			// 执行业务逻辑, 执行一个读取输入回显的操作
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err: ", err)
						continue
					} 
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("send buf err: ", err)
						continue
					}
				}
			}()
		}

	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Ngo server , name ", s.Name)
	//执行关闭相关的逻辑
}

func (s *Server) Serve() {
	s.Start()
	//服务在启动时需要处理的逻辑
	//select阻塞，防止主进程退出，协程自动结束
	select {

	}

}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8888,
	}
	return s
}
