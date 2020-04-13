package ztimer

import (
	"time"
)

const (
	HOUR_NAME = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3
	HOUR_SCALES = 12

	MINUTE_NAME = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES = 60

	SECOND_NAME = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES = 60

	TIME_MAX_CAP = 2048		//每个时间轮刻度上挂载定时器的最大个数
)


//定时器实现
type Timer struct {
	//延迟调用函数
	delayFunc *DelayFunc 
	//调用时间
	unixts int64 
}

//返回1970-1-1至今经历的毫秒数
func UnixMill() int64 {
	return time.Now().UnixNano/1e6
}

/*
	创建一个定时器，在指定的时间触发定时器方法
*/

func NewTimerAt(df *DelayFunc, unixNano int64) *Timer {
	return &Timer {
		delayFunc: df,
		unixts: unixNano/1e6,  //转换为毫秒
	}

}

//创建一个定时器，在当前时间延迟duration之后触发，定时器方法
func NewTimerAfter(df *DelayFunc, duration time.Duration) *Timer {
	return NewTimerAt(df, time.Now().UnixNano()+int64(duration))
}

//启动定时器
func (t *Timer) Run() {
	go func() {
		now := UnixMill()
		//设置的时间是否在当前时间之后
		if t.unixts > now {
			time.Sleep(time.Duration(t.unixts - now) * time.Millisecond)
		}
		//调用事先注册好的超市延迟方法
		t.delayFunc.Call()
	}()
}