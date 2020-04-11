package ziface

//定义消息接口，将请求的一个消息封装到message中
type IMessage interface {
	GetDataLen() uint32 
	GetMsgId() uint32
	GetData() []byte

	SetMsgId(uint32) //设置消息id
	SetData([]byte) //设置消息内容
	SetDataLen(uint32) //设置消息数据段长度
}