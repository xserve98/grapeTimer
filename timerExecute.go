/// Author:koangel
/// jackliu100@gmail.com
/// 每个单独的执行模块，结构中越简单越好
package grapeTimer

import (
	"encoding/json"
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	// 通过字符串转换下一次的日期
	timerDateMode = iota
	// 通过转成的Tick日期进行下一次时间计算
	timerTickMode
)

type GrapeExecFn func(timerId int, args interface{})

// 每个执行周期的执行者
// 每个时间执行器
// 可Json化并反Json化
type GrapeTimer struct {
	Id        int    `json:"TimerId"`
	NextTime  int64  `json:"nextUnix"`
	RunMode   int    `json:"Mode"`
	TimeData  string `json:"timeData"`
	LoopCount int32  `json:"loopCount"` // -1为无限执行

	exectFunc  GrapeExecFn // 执行的函数 不保存
	execArgs   interface{} // 执行参数 外部不可访问和改变
	tickSecond int         // 临时使用 不可保存
}

/// 直接创建一个timer 内部函数
func newTimer(Id, Mode, Count int, timeData string, fn GrapeExecFn, args interface{}) *GrapeTimer {
	newValue := &GrapeTimer{
		Id:         Id,
		RunMode:    Mode,
		TimeData:   timeData,
		LoopCount:  int32(Count),
		NextTime:   0,
		tickSecond: 0,

		execArgs:  args,
		exectFunc: fn,
	}

	newValue.makeNextTime()

	return newValue
}

/// 通过Json创建一个Timer 内部函数
func newTimerFromJson(s string, fn GrapeExecFn, args interface{}) *GrapeTimer {
	newValue := &GrapeTimer{}
	newValue = newValue.ParserJson(s)
	if newValue != nil {
		newValue.exectFunc = fn
		newValue.execArgs = args
		newValue.makeNextTime()
	}

	return newValue
}

/// 执行Timer
func (c *GrapeTimer) Execute() {
	if c.IsDestroy() {
		return // 销毁中的不可执行
	}

	if c.IsExpired() {
		if CDebugMode {
			log.Printf("[grapeTimer] Timer Execute:%v |time:%v| Begin", c.Id, time.Now())
		}
		// 执行一下
		if UseAsyncExec {
			go c.exectFunc(c.Id, c.execArgs)
		} else {
			c.exectFunc(c.Id, c.execArgs)
		}

		if CDebugMode {
			log.Printf("[grapeTimer] Timer Execute:%v |time:%v| End", c.Id, time.Now())
		}

		c.nextTime() // 计数
	}
}

/// 是否到时间
func (c *GrapeTimer) IsExpired() bool {
	if time.Now().Unix() >= c.NextTime {
		if CDebugMode {
			log.Printf("[grapeTimer] Timer Expired:%v |time:%v|", c.Id, time.Now())
		}
		return true
	}

	return false
}

/// 是否可销毁
func (c *GrapeTimer) IsDestroy() bool {
	val := atomic.LoadInt32(&c.LoopCount)
	if val == -1 {
		return false
	}

	if val == 0 {
		return true
	}

	return false
}

func (c *GrapeTimer) makeNextTime() {
	if CDebugMode {
		log.Printf("[grapeTimer] Timer NextTime:%v |time:%v|LoopCount:%v|", c.Id, time.Now(), c.LoopCount)
	}

	switch c.RunMode {
	case timerDateMode:
		vtime, err := Parser(c.TimeData)
		if err != nil {
			atomic.StoreInt32(&c.LoopCount, 0) // 出错销毁
			log.Printf("[grapeTimer] Timer NextTime:%v |time:%v|Error:%v|", c.Id, time.Now(), err)
			return
		}

		c.NextTime = vtime.Unix() // 生成下一次时间
		break
	case timerTickMode:
		if c.tickSecond == 0 {
			c.tickSecond, _ = strconv.Atoi(c.TimeData)
		}

		c.NextTime = time.Now().Add(time.Duration(c.tickSecond) * time.Microsecond).Unix()
		break
	}
}

/// 计算下一次时间
func (c *GrapeTimer) nextTime() {
	if c.IsDestroy() {
		return
	}

	val := atomic.LoadInt32(&c.LoopCount)
	// 先进行计数
	if val != -1 {
		atomic.AddInt32(&c.LoopCount, -1)
		val = atomic.LoadInt32(&c.LoopCount)
		if val <= 0 {
			val = 0
		}
	}

	c.makeNextTime()
}

func (c *GrapeTimer) Stop() {
	atomic.StoreInt32(&c.LoopCount, 0) // 停止这个timer继续执行
}

func (c *GrapeTimer) toJson() string {
	bJson, err := json.Marshal(c)
	if err != nil {
		log.Printf("[grapeTimer] toJson:%v|Error:%v|", c.Id, err)
		return ""
	}

	return string(bJson)
}

func (c *GrapeTimer) ParserJson(s string) *GrapeTimer {
	err := json.Unmarshal([]byte(s), c)
	if err != nil {
		log.Printf("[grapeTimer] ParserJson:%v|Error:%v|", c.Id, err)
		return nil
	}

	return c
}
