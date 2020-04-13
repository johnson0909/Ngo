package ztimer

import (
	"fmt"
	"reflect"
)

/*
	定义一个延迟调用函数
	当时间定时器超时的时候，触发事先注册好的回调函数
*/

type DelayFunc struct {
	f  func(...interface{}) 
	args []interface{}
}

//创建一个延迟调用函数
func NewDelayFunc( f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f: f,
		args: args,
	}
}

//打印当前延迟函数的信息，用于记录日志
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFunc:%s, args:%v}", reflect.TypeOf(df.f).Name(), df.args)
}

//执行延迟函数
func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Call err:", err)
		}
	}()

	df.f(df.args...)
}