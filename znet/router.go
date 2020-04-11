package znet

import "Ngo/ziface"
//实现Router时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {}

/*
	所有方法为空，因为有的Router可能不需要PreHandle和PostHandle
	让用户根据需要进行重写
*/
func (br *BaseRouter) PreHandle(req ziface.IRequest) {}
func (br *BaseRouter) Handle(req ziface.IRequest) {}
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
