package ztimer

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TimeWheel struct {
	//时间轮名称
	name string
	//刻度的时间间隔
	interval int64 
	//每个时间轮上的刻度数
	scales int
	//当前时间指针的指向
	curIndex int
	//每个刻度所存放的timer定时器的最大容量
	maxCap int
	//当前时间轮上所有的timer
	timerQueue map[int]map[uint32]*Timer //int表示当前时间轮的刻度
	//下一层时间轮
	nextTimeWheel *TimeWheel
	//互斥锁
	sync.RWMutex
}
/*
	创建一个时间轮
	name: 时间轮的名称
	intercal: 每个刻度之间的duration
	scales: 每个轮盘一共有多少个刻度
	maxCap: 每个刻度上的最大Timer定时器个数
*/
func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	tw := &TimeWheel{
		name: name,
		interval: interval,
		scales: scales,
		maxCap: maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
	}
	//初始化map
	for i := 0; i < scales; i++ {
		tw.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}
	fmt.Println("Init timerWheel name = ", tw.name, " is Done!")
	return tw
}
/*
	将一个timer定时器加入到分层时间轮中
	tid：每个定时器timer的唯一标识
	t: 当前被加入时间轮的定时器
	forceNext: 是否强制将定时器添加到下一层定时器
	算法：
	如果当前timer的超时时间间隔大于一个刻度，那么进行hash计算找到对应的刻度上添加
*/
func (tw *TimeWheel) addTimer(tid uint32, t *Timer, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errstr := fmt.Sprintf("addTimer function err: %s", err)
			fmt.Println(errstr)
			return errors.New(errstr)
		}
		return nil
	}()
	//得到当前的超时时间间隔(ms)毫秒为单位
	delayInterval := t.unixts - UnixMill()

	//如果当前的超时时间 大于一个刻度的时间间隔
	if delayInterval >= tw.interval {
		//计算需要跨越几个刻度
		dn := delayInterval / tw.interval
		tw.timerQueue[(tw.curIndex+int(dn)%tw.scales)][tid] = t
		return nil
	}

	//如果当前超时时间小于一个刻度的时间间隔，并且当前时间轮没有下一层精度更小的时间轮
	if delayInterval < tw.interval && tw.nextTimeWheel == nil {
		if forceNext == true {
			//如果设置为强制移至下一刻度，那么将定时器移至下一个刻度
			//这种情况。主要是时间轮自动轮转的情况
			//因为这是底层时间轮，该定时器在转动的时候如果没有被调度者取走的话，该定时器将不会再被发现
			//因为时间轮刻度已经过去，如果不强制把该定时器Timer移至下一时刻，将永远不会被取走并触发调用
			//所以这里强制将timer移至下个刻度的集合中，等待调用者再下次轮转之前取走该定时器
			tw.timerQueue[(tw.curIndex+1)%tw.scales][tid] = t 
		} else {
			//如果手动添加定时器，那么直接将timer添加到对应底层时间轮的当前刻度
			tw.timerQueue[tw.curIndex][tid] = t 
		}
	}

	//如果当前的超时时间，小于一个刻度的时间间隔， 并且有下一层时间轮
	if delayInterval < tw.interval {
		return tw.nextTimeWheel.AddTimer(tid, t)
	}

	return nil
}
//添加一个timer到一个时间轮中
func (tw *TimeWheel) AddTimer(tid uint32, t *Timer) error {
	tw.Lock()
	defer tw.Unlock()

	return tw.addTimer(tid, t, false)
}

//删除一个定时器，根据定时器id
func (tw *TimeWheel) RemoveTimer(tid uint32) {
	tw.Lock()
	defer tw.Unlock()

	for i := 0; i < tw.scales; i++ {
		if _, ok := tw.timerQueue[i][tid]; ok {
			delete(tw.timerQueue[i], tid)
		}
	}
}


//给一个时间轮添加下层时间轮，比如给小时时间轮添加分钟时间轮，给分钟时间轮添加秒钟时间轮
func (tw *TimeWheel) AddTimeWheel(next *TimeWheel) {
	tw.nextTimeWheel = next
	fmt.Println("Add timeWheel[", tw.name, "]'s next [", next.name, "] is succ!")


}

//启动时间轮
func (tw *TimeWheel) run() {
	for {
		//时间轮每间隔interval一刻度时间，触发一次转动
		time.Sleep(time.Duration(tw.interval) * time.Millisecond)

		tw.Lock()
		//取出挂载再当前刻度的全部定时器
		curTimers := tw.timerQueue[tw.curIndex]
		//当前定时器要重新添加 所给当前刻度重新开辟一个map Timer容器
		tw.timerQueue[tw.curIndex] = make(map[uint32]*Timer, tw.maxCap)
		for tid, timer := range curTimers {
			tw.addTimer(tid, timer, true)
		}
		//当前刻度指针 走一格
		tw.curIndex = (tw.curIndex + 1) % tw.scales
		tw.Unlock()
	}
}

//非阻塞的方式让时间轮
func (tw *TimeWheel) Run() {
	go tw.run()

	fmt.Println("timerwheel name = ", tw.name, " is running ...")
}

//获取定时器再一段时间间隔内的timer
func (tw *TimeWheel) GetTimerWithIn(duration time.Duration) map[uint32]*Timer {
	//最终触发的定时器一定是挂载再最底层时间轮上的定时器
	//找到底层时间轮
	leaftw :=  tw 
	for leaftw.nextTimeWheel != nil {
		leaftw = leaftw.nextTimeWheel 
	}

	leaftw.Lock()
	defer leaftw.Unlock()

	//返回Timer集合
	timerList := make(map[uint32]*Timer)

	now := UnixMill()
	//取出当前时间刻度内全部Timer
	for tid, timer := range leaftw.timerQueue[leaftw.curIndex] {
		if timer.unixts - now < int64(duration/1e6) {
			//当前定时器已经超时
			timerList[tid] = timer
			//定时器已经被超时取走，从当前时间轮上删除
			delete(leaftw.timerQueue[leaftw.curIndex], tid)
		}
	}

	return timerList
}