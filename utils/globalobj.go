package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"Ngo/ziface"
)
//定义全局配置接口
type GlobalObj struct {
	//Server
	TcpServer ziface.IServer
	Host string
	TcpPort int
	Name string
	//Ngo
	Version string
	MaxPacketSize uint32
	MaxConn int        //主机允许的最大连接数
	WorkerPoolSize uint32 //工作池的worker数量
	MaxWorkerTaskLen uint32 //每个worker中任务队列的最大任务数
	MaxMsgChanLen uint32   //SendBuffMsg发消息的缓冲最大长度

	ConfFilePath string

}

//定义一个全局对象
var GlobalObject *GlobalObj

//判断文件是否存在 
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

//读取配置文件
func (g *GlobalObj) Reload() {
	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}

	//将json数据解析存到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
//提供默认加载
func init() {
	//初始化GlobalObject变量，设置初始值
	GlobalObject = &GlobalObj {
		Name: "NgoDemo",
		Version: "V1.0",
		TcpPort: 8888,
		Host: "0.0.0.0",
		MaxConn: 12000,
		MaxPacketSize: 4096,
		ConfFilePath: "conf/Ngo.json",
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen: 1024,
	}
	//从配置文件中加载
	GlobalObject.Reload()
}