// grapeTimer project main.go
package main

import (
	"fmt"
	"time"

	"github.com/koangel/grapeTimer"
)

var RunTick = 0

func OnOnceTick(timerId int, args interface{}) {
	fmt.Printf("OnceTick:%v\n", args)

	RunTick++
}

func OnOnceDayTick(timerId int, args interface{}) {
	fmt.Printf("OnOnceDayTick:%v\n", args)

	RunTick++
}

func main() {
	grapeTimer.InitGrapeScheduler(500*time.Microsecond, true)
	grapeTimer.CDebugMode = false
	grapeTimer.UseAsyncExec = true
	// 10万的单次TIMER测试
	for i := 0; i < 1000000; i++ {
		grapeTimer.NewTickerOnce(1000, OnOnceTick, i+1)
	}

	for i := 0; i < 50000; i++ {
		grapeTimer.NewTimeDataOnce("Day 21:38:00", OnOnceDayTick, i+1)
	}

	for {
		time.Sleep(1 * time.Second)
	}

}
