# **grapeTimer 时间调度器**

一款粗粒度的时间调度器，可以帮你通过一些字符串快速并简单的创建时间任务。

用于游戏服务端的优化设计，大量并行的时间调度方式。

- Author: Koangel
- Blog: [http://koangel.github.com](http://koangel.github.com)
- Weibo: [@koangel](http://weibo.com/koangel)
- Homepage: [个人博客](http://grapec.me)

#### 简单介绍：
- 通过命令格式创建`time.Time`
- 简洁的Api格式，轻度且可拆分的函数库
- 快速创建调度器
- 可控的调度器时间粒度
- 高性能的并发调度
- 时间周期，次数多模式可控`[支持每天，每周，每月]`
- *可选择对调度器保存或内存执行[待实现]
- *生成可保存的调度器字符串并反向分析他生成调度器[待实现]
- 不依赖第三方库

## **简单测试**

100W个TIMER压入4秒钟，每个执行1S的话，完成执行大概不到4S，基本上达到性能需求。
如果超过100W个并发TIMER，建议切分服务，所有执行采用异步行为。

## **安装方法**

```
go get -u -v github.com/koangel/grapeTimer
```

## **基本用法**

``` Go
// 初始化一个1秒钟粒度的调度器，ars代表是否自动设置运行为并行模式
grapeTimer.InitGrapeScheduler(1*time.Second, true)
// 启动一个单次执行的调度器，1秒时间，基本tick单位为毫秒
Id := grapeTimer.NewTickerOnce(1000, exec100, nil)
// 启动一个1秒为周期的 循环timer 循环100次 -1为永久循环
Id = grapeTimer.NewTickerLoop(1000,100, exec100Loop, nil)
// 启动一个每日规则的定时器，参数为args data
Id = grapeTimer.NewTimeDataOnce("Day 13:59:59", exeDayTime, "args data")
// 启动一个每日循环规则的定时器，参数为args data 循环100次 -1为永久循环
Id = grapeTimer.NewTimeDataLoop("Day 13:59:59",100, exeDayTime, "args data")
```

所有执行函数请保持以下函数格式：
```
func exec100(timerId int, args interface{}) {
	fmt.Printf("exec100")
}

func exec100Loop(timerId int, args interface{}) {
	fmt.Printf("exec100Loop")
}

func exeDayTime(timerId int, args interface{}) {
	fmt.Println("exeDayTime:", args)
}
```
## **停止计时器**

```Go
// 将自动返回的ID作为参数传入可停止持续循环的TIMER
grapeTimer.StopTimer(Id)
```

## **参数设置**

```
// 设置启用日志调试模式，建议正式版本关闭他
grapeTimer.CDebugMode = true
// 调用分析器使用的时区，可以个根据不同国家地区设置 
grapeTimer.LocationFormat = "Asia/Shanghai"
// 开启异步调度模式，在此模式下 timer执行时会建立一个go，不会阻塞其他timer执行调度，建议开启
grapeTimer.UseAsyncExec = true
```

## **保存调度器**

```
待添加
```

## **生成调度器字符串**

```
待添加
```

## **可用格式说明**

调度器有轻度的日期模式分析体系，可提供每日，每周，每月的时间日期生成方式，具体格式如下：

|关键字|格式|说明|
|:----------:|:-------:|:----------:|
|Day|Day 00:00:00|生成每日的日期时间|
|Week|Week 1 00:00:00|生成每周的日期时间， 0~6 分别代表周日到周六|
|Month|Month 1 00:00:00|生成每月该日期的时间，建议不要使用20日之后的日期|

以上日期如果超过时间自动生成下一个时间，每月时间使用时，如本月不存在该日期，则直接返回错误。

代码示例：

```go
	vtime, err := grapeTimer.Parser("Day 13:59:59") // 返回值为标准的*time.Time
	if err != nil {
		// 处理错误...
	}

        vtime, err = grapeTimer.Parser("Week 6 23:59:59")
	if err != nil {
		// 处理错误...
	}

	vtime, err = grapeTimer.Parser("Month 26 13:59:59")
	if err != nil {
		// 处理错误...
	}
```
代码输出：

```
2017-05-24 13:59:59 +0800 CST
```