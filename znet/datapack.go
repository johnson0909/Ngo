package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"Ngo/ziface"
	"Ngo/utils"
)

//封包拆包实例，暂不需要成员
type DataPack struct {

}

func NewDataPack() *DataPack {
	return &DataPack{}
}
//获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) + DataLen uint32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil

}

//拆包方法
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	
	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超过我们允许的最大包长度
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data recevied")
	}

	//这里只需要把head的数据拆包出来就可以了，然后根据head的长度，从conn中读取一次长度
	return msg, nil
}