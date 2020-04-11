package ziface

//定义封包拆包接口
//为解决TCP连接中的数据流粘包问题

type IDataPack interface {
	GetHeadLen() uint32  //获取包头长度
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error) //拆包方法
}