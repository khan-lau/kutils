// 精确到s的时间范围计算
package datetime

import (
	Time "time"
)

type Duration int64

const (
	DATETIME_FORMATTER          = "2006-01-02 15:04:05"
	DATETIME_FORMATTER_Mill     = "2006-01-02 15:04:05.000"
	DATETIME_TIMEZONE_FORMATTER = "2006-01-02 15:04:05 -0700"
)

const (
	NONE   Duration = 0
	SECOND Duration = 1
	MINUTE Duration = 2
	HOUR   Duration = 3
	DAY    Duration = 4
	WEEK   Duration = 5
	MONTH  Duration = 6
	YEAR   Duration = 7
)

type RecordItem struct {
	Timestamp uint64  // 时间戳, 精确到s
	Value     float64 // 值
}

type TimeSlice struct {
	Start uint64 // 时间戳, 精确到s
	End   uint64 // 时间戳, 精确到s
	Tag   any    // 结果集

	//Tag *klists.KList[*RecordItem]
}

// 获取指定时间所在月份的天数
func MonthDays(time Time.Time) uint {
	return uint(Time.Date(time.Year(), time.Month()+1, 1, 0, 0, 0, 0, time.Location()).Sub(Time.Date(time.Year(), time.Month(), 1, 0, 0, 0, 0, time.Location())).Hours() / 24)
}

// 获取指定时间所在年份的天数
func YearDays(time Time.Time) uint {
	return uint(Time.Date(time.Year()+1, 1, 1, 0, 0, 0, 0, time.Location()).Sub(Time.Date(time.Year(), 1, 1, 0, 0, 0, 0, time.Location())).Hours() / 24)
}

// 当前时间所在 分 的起始秒时间戳
func MinuteFirst(time Time.Time) uint64 {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, time.Hour(), time.Minute(), 0, 0, time.Location())
	return uint64(start.Unix())
}

// 当前时间所在 分 的结束秒时间戳
func MinuteLast(time Time.Time) uint64 {
	year, month, day := time.Date()
	end := Time.Date(year, month, day, time.Hour(), time.Minute(), 59, 999999999, time.Location())
	return uint64(end.Unix())
}

// 取1分钟的开始时间与结束时间
func FirstAndLastMinute(time Time.Time) (uint64, uint64) {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, time.Hour(), time.Minute(), 0, 0, time.Location())
	end := start.Add(Time.Nanosecond*1000*1000*1000*60 - 1)
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 小时 的起始秒时间戳
func HourFirst(time Time.Time) uint64 {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, time.Hour(), 0, 0, 0, time.Location())
	return uint64(start.Unix())
}

// 当前时间所在 小时 的结束秒时间戳
func HourLast(time Time.Time) uint64 {
	year, month, day := time.Date()
	end := Time.Date(year, month, day, time.Hour(), 59, 59, 999999999, time.Location())
	return uint64(end.Unix())
}

// 取1小时的开始时间与结束时间
func FirstAndLastHour(time Time.Time) (uint64, uint64) {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, time.Hour(), 0, 0, 0, time.Location())
	end := start.Add(Time.Nanosecond*1000*1000*1000*60*60 - 1)
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 天 的起始秒时间戳
func DayFirst(time Time.Time) uint64 {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix())
}

// 当前时间所在 天 的结束秒时间戳
func DayLast(time Time.Time) uint64 {
	year, month, day := time.Date()
	end := Time.Date(year, month, day, 23, 59, 59, 999999999, time.Location())
	return uint64(end.Unix())
}

// 取1天的开始时间与结束时间
func FirstAndLastDay(time Time.Time) (uint64, uint64) {
	year, month, day := time.Date()
	start := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	end := start.Add(Time.Nanosecond*1000*1000*1000*60*60*24 - 1)
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 周 的起始秒时间戳
func WeekFirst(time Time.Time) uint64 {
	weekDay := int(time.Weekday())
	if weekDay == 0 {
		weekDay = 6 // 周日
	} else {
		weekDay -= 1
	}
	// 0 = Monday, 1 = Tuesday, 2 = Wednesday, 3 = Thursday, 4 = Friday, 5 = Saturday, 6 = Sunday

	year, month, day := time.Date()
	curr := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	start := curr.Add(-(Time.Nanosecond * 1000 * 1000 * 1000 * 60 * 60 * 24 * Time.Duration(weekDay)))
	return uint64(start.Unix())
}

func WeekLast(time Time.Time) uint64 {
	weekDay := int(time.Weekday())
	if weekDay == 0 {
		weekDay = 6
	} else {
		weekDay -= 1
	}

	year, month, day := time.Date()
	curr := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	start := curr.Add(-(Time.Nanosecond * 1000 * 1000 * 1000 * 60 * 60 * 24 * Time.Duration(weekDay)))
	end := start.Add(Time.Nanosecond*1000*1000*1000*60*60*24*7 - 1)
	return uint64(end.Unix())
}

// 取1周的开始时间与结束时间
func FirstAndLastWeek(time Time.Time) (uint64, uint64) {
	weekDay := int(time.Weekday())
	if weekDay == 0 {
		weekDay = 6
	} else {
		weekDay -= 1
	}

	year, month, day := time.Date()
	start := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	// fmt.Printf("datetime: %s, week: %d, start: %s\n", time.Format(DATETIME_FORMATTER), weekDay, start.Format(DATETIME_FORMATTER))

	start = start.Add(-(Time.Nanosecond * 1000 * 1000 * 1000 * 60 * 60 * 24 * Time.Duration(weekDay)))
	end := start.Add(Time.Nanosecond*1000*1000*1000*60*60*24*7 - 1)
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 月 的起始秒时间戳
func MonthFirst(time Time.Time) uint64 {
	year, month, _ := time.Date()
	start := Time.Date(year, month, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix())
}

// 当前时间所在 月 的结束秒时间戳
func MonthLast(time Time.Time) uint64 {
	year, month, _ := time.Date()
	end := Time.Date(year, month, 1, 0, 0, 0, 0, time.Location()).AddDate(0, 1, 0).Add(-Time.Nanosecond)
	return uint64(end.Unix())
}

// 取1个月的开始时间与结束时间
func FirstAndLastMonth(time Time.Time) (uint64, uint64) {
	year, month, _ := time.Date()
	end := Time.Date(year, month, 1, 0, 0, 0, 0, time.Location()).AddDate(0, 1, 0).Add(-Time.Nanosecond)
	start := Time.Date(year, month, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 年 的起始秒时间戳
func YearFirst(time Time.Time) uint64 {
	year, _, _ := time.Date()
	start := Time.Date(year, 1, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix())
}

// 当前时间所在 年 的起始秒时间戳
func YearLast(time Time.Time) uint64 {
	year, _, _ := time.Date()
	curr := Time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Location())
	return uint64(curr.Add(-Time.Nanosecond).Unix())
}

// 取1年的开始时间与结束时间
func FirstAndLastYear(time Time.Time) (uint64, uint64) {
	year, _, _ := time.Date()
	end := Time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Location())
	end = end.Add(-Time.Nanosecond)
	start := Time.Date(year, 1, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix()), uint64(end.Unix())
}
