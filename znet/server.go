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
		//0 启动协程持工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error, err: ", err)
			return 
		}
		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen err: ", err)
			return
		}
		//监听成功
		fmt.Println("start Ngo server ", s.Name, " succ! now listenning...")
		//TODO 自动生成id的方法
		var cid uint32
		cid = 0
		
		//3 启动server网络连接业务 
		for {
			//3.1 阻塞等待网络连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())
			//3.2 设置服务器最大连接控制，如果超过最大连接数，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			
			//3.3 处理该新连接请求的业务方法，此时应该有handler和conn时绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			
			//3.4 启动当前连接的处理业务
			go dealConn.Start()
		}

	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] Ngo server , name ", s.Name)
	//执行关闭相关的逻辑
	s.ConnMgr.ClearConn()
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
