package ztimer

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	//默认缓冲触发函数队列的大小
	MAX_CHAN_BUFF = 2048
	//默认最大误差时间
	MAX_TIME_DELAY = 100
)

type TimerScheduler struct {
	//当前调度的最高级时间轮
	tw *TimeWheel
	//定时器编号累加器
	idGen uint32
	//已经触发定时器的channel
	triggerChan chan *DelayFunc

	sync.RWMutex
}

//返回一个定时器调度器
func NewTimerScheduler() *TimerScheduler {
	//秒级时间轮
	second_tw := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	//分钟级时间轮
	minute_tw := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	//小时级时间轮
	hour_tw := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	second_tw.Run()
	minute_tw.Run()
	hour_tw.Run()

	return &TimerScheduler {
		tw: hour_tw,
		triggerChan: make(chan *DelayFunc, MAX_CHAN_BUFF),
	}

}

//创建一个定点Timer 并将Timer添加到分层的时间轮中，返回Timer的tid
func (this *TimerScheduler) CreateTimerAt(df *DelayFunc, unixNano int64) (uint32, error) {
	this.Lock()
	defer this.Unlock()

	this.idGen++
	return this.idGen, this.tw.AddTimer(this.idGen, NewTimerAt(df, unixNano))
}

//创建一个延迟timer 并将Timer添加到分层时间轮中，返回Timer的tid
func (this *TimerScheduler) CreateTimerAfter(df *DelayFunc, duration time.Duration) (uint32, error){
	this.Lock()
	defer this.Unlock()
	this.idGen++
	return this.idGen, this.tw.AddTimer(this.idGen, NewTimerAfter(df, duration))
}

//删除timer
func (this *TimerScheduler) CancelTimer(tid uint32) {
	this.Lock()
	this.Unlock()

	this.tw.RemoveTimer(tid)
}

//获取计时结束的延迟执行函数通道
func (this *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return this.triggerChan
}

//非阻塞的方式启动timerScheduler
func (this *TimerScheduler) Start() {
	go func() {
		for {
			now := UnixMill()
			//获取最近的MAX_TIME_DELAY毫秒的超时定时器集合
			timerList := this.tw.GetTimerWithIn(MAX_TIME_DELAY * time.Millisecond)
			for _, timer := range timerList {
				if math.Abs(float64(now - timer.unixts)) > MAX_TIME_DELAY {
					//已经超出的定时器报警
					fmt.Println("want call at ", timer.unixts, "; real call at", now, "; delay ", now-timer.unixts)
				}
				this.triggerChan <- timer.delayFunc
			}
			time.Sleep(MAX_TIME_DELAY / 2 * time.Millisecond)
		}
	}()
}

//自动调度
func NewAutoExecTimerScheduler() *TimerScheduler {
	//创建一个调度器
	autoExecScheduler := NewTimerScheduler()
	//启动调度器
	autoExecScheduler.Start()

	//永久从调度中获取超时，触发的函数并执行
	go func() {
		delayFuncChan := autoExecScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			go df.Call()
		}
	}()

	return autoExecScheduler
}