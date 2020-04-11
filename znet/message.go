package znet

type Message struct {
	DataLen uint32 //消息长度
	Id 		uint32 //消息id
	Data 	[]byte//消息数据
}
//创建一个Message消息
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id: 	 id,
		Data:	 data,
	}
}
//创建一个Message
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

//获取消息id
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

//获取消息数据
func (msg *Message) GetData() []byte {
	return msg.Data
}

//设置消息长度
func (msg *Message) SetDataLen(len uint32)  {
	msg.DataLen = len
}

//设置消息id
func (msg *Message) SetMsgId(msgId uint32) {
	msg.Id = msgId
}

//设置消息数据
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}