// 精确到s的时间范围计算
package datetime

import (
	Time "time"

	"github.com/khan-lau/kutils/container/klists"
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

type TimeSlice struct {
	Start uint64 // 时间戳, 精确到s
	End   uint64 // 时间戳, 精确到s
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
		weekDay = 6
	} else {
		weekDay -= 1
	}

	year, month, day := time.Date()
	start := Time.Date(year, month, day, 0, 0, 0, 0, time.Location())
	start = start.Add(-(Time.Nanosecond * 1000 * 1000 * 1000 * 60 * 60 * 24 * Time.Duration(weekDay)))
	return uint64(start.Unix())
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

// 取1个月的开始时间与结束时间
func FirstAndLastMonth(time Time.Time) (uint64, uint64) {
	year, month, _ := time.Date()
	end := Time.Date(year, month+1, 1, 0, 0, 0, 0, time.Location())
	end = end.Add(-Time.Nanosecond)
	start := Time.Date(year, month, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix()), uint64(end.Unix())
}

// 当前时间所在 年 的起始秒时间戳
func YearFirst(time Time.Time) uint64 {
	year, _, _ := time.Date()
	start := Time.Date(year, 1, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix())
}

// 取1年的开始时间与结束时间
func FirstAndLastYear(time Time.Time) (uint64, uint64) {
	year, _, _ := time.Date()
	end := Time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Location())
	end = end.Add(-Time.Nanosecond)
	start := Time.Date(year, 1, 1, 0, 0, 0, 0, time.Location())
	return uint64(start.Unix()), uint64(end.Unix())
}

// 将起止时间按指定周期分割, 返回每个周期的起止时间
//   - @param time.Time start    开始时间
//   - @param time.Time end      结束时间
//   - @param Duration  duration 分割周期
//   - @return []*TimeSlice  每个分段的起止时间
func SplitDuration(start, end Time.Time, duration Duration) []*TimeSlice {
	list := klists.New[*TimeSlice]()
	if duration == NONE {
		list.PushBack(&TimeSlice{Start: uint64(start.Unix()), End: uint64(end.Unix())})
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == SECOND {
		first := start.Unix()
		last := end.Unix()
		for i := first; i <= last; i++ {
			list.PushBack(&TimeSlice{Start: uint64(i), End: uint64(i)})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MINUTE {
		first := MinuteFirst(start)
		last := uint64(end.Unix())
		step := uint64(60)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := FirstAndLastMinute(Time.Unix(int64(i), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == HOUR {
		first := HourFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := FirstAndLastHour(Time.Unix(int64(i), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == DAY {
		first := DayFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := FirstAndLastDay(Time.Unix(int64(i), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == WEEK {
		first := WeekFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * 7)
		// fmt.Printf("-- datetime: %s, weekFirst: %s\n", start.Format(DATETIME_FORMATTER),
		// 	Time.Unix(int64(first), 0).Format(DATETIME_FORMATTER))

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := FirstAndLastWeek(Time.Unix(int64(i), 0))

			// fmt.Printf("-- datetime: %s, start: %s, end: %s\n", Time.Unix(int64(i), 0).Format(DATETIME_FORMATTER),
			// 	Time.Unix(int64(valBegin), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(valEnd), 0).Format(DATETIME_FORMATTER))

			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MONTH {
		first := MonthFirst(start)
		last := uint64(end.Unix())

		step := uint64(60 * 60 * 24 * MonthDays(start))

		for i := first; i <= last; {
			valBegin, valEnd := FirstAndLastMonth(Time.Unix(int64(i), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			monthDays := MonthDays(Time.Unix(int64(i), 0))
			step = uint64(60 * 60 * 24 * monthDays)
			i = i + step
			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == YEAR {
		first := YearFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * YearDays(start))

		for i := first; i <= last; {
			valBegin, valEnd := FirstAndLastYear(Time.Unix(int64(i), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			yearDays := YearDays(Time.Unix(int64(i), 0))
			step := uint64(60 * 60 * 24 * yearDays)
			i = i + step

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	return klists.ToKSlice[*TimeSlice](list)
}
