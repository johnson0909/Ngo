package ziface

/*
	定义路由接口
	路由是指 使用框架者给该链接自定义的业务处理方法
*/

type IRouter interface {
	PreHandle(request IRequest) //处理conn业务之前的hook方法
	Handle(requets IRequest)   //处理conn业务的方法
	PostHandle(request IRequest) //处理conn业务之后的hook方法
}