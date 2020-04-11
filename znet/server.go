package znet
import (
	"Ngo/utils"
	"fmt"
	"net"
	"Ngo/ziface"
)
type Server struct {
	//服务器名称
	Name string
	IPVersion string
	IP string
	Port int
	//当前server的消息管理模块，用来绑定MsgID和对应处理方法
	msgHandler ziface.IMsgHandle
	//当前server的链接管理模块
	ConnMgr ziface.IConnManager
	//该server链接创建时的hook函数
	OnConnStart func(conn ziface.IConnection)
	//该server链接断开时的hook函数
	OnConnStop func(conn ziface.IConnection)

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

//路由功能给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

//得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}
//设置server创建时的hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置server销毁时的hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用链接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart ...")
		s.OnConnStart(conn)
	}
}

//调用链接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop ...")
		s.OnConnStop(conn)
	}
}

func NewServer() ziface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}
