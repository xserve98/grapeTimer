/// Author:koangel
/// jackliu100@gmail.com
/// 负责分析并解析字符串格式并转换成一个time
package grapeTimer

import "time"

const (
	grapeTimeFormat = "2006-01-02 15:04:05"
)

/// 基本格式 Day or Week or Month or Year
/// like : Type DayNum Time
/// format: Day 1 00:00:00
/// 返回一个处理好的下一个时间
func Parser(fmt string) time.Time {
	// 默认使用上海时区

	return ParserLoc()
}

func ParserLoc(fmt string, loc *time.Location) time.Time {

}
