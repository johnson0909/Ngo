package ziface

/* 
	定义请求接口
	封装了客户端请求的链接信息和请求的数据到Request
*/

type IRequest interface {
	GetConnection() IConnection 
	GetData() []byte
	GetMsgID() uint32
}