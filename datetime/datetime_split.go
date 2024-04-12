package datetime

import (
	Time "time"

	"github.com/khan-lau/kutils/container/klists"
)

// 按指定时间间隔 进行1个自然周期的切割
//
// 参数
//   - start: 开始时间
//   - end: 结束时间
//   - duration: 时间间隔周期
//
// 返回值
//   - []*TimeSlice: 切割后的时间片段
func SplitDuration(start, end Time.Time, duration Duration) []*TimeSlice {
	return SplitNaturalDuration(start, end, 1, duration)
}

// 按指定时间间隔 进行自然周期的切割
//
// 参数
//   - start: 开始时间
//   - end: 结束时间
//   - unit: 间隔数量, 时间间隔等于 unit * duration
//   - duration: 时间间隔周期
//
// 返回值
//   - []*TimeSlice: 切割后的时间片段
func SplitNaturalDuration(start, end Time.Time, unit uint, duration Duration) []*TimeSlice {
	if unit == 0 {
		unit = 1
	}
	list := klists.New[*TimeSlice]()
	if duration == NONE {
		list.PushBack(&TimeSlice{Start: uint64(start.Unix()), End: uint64(end.Unix()), Tag: nil}) //  Records: klists.New[*RecordItem]()
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == SECOND {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(unit)
		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(i), uint64(i+step-1)
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MINUTE {
		first := MinuteFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := MinuteFirst(Time.Unix(int64(i), 0)), MinuteLast(Time.Unix(int64(i+step-(1*60)), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == HOUR {
		first := HourFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := HourFirst(Time.Unix(int64(i), 0)), HourLast(Time.Unix(int64(i+step-(1*60*60)), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == DAY {
		first := DayFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := DayFirst(Time.Unix(int64(i), 0)), DayLast(Time.Unix(int64(i+step-(1*60*60*24)), 0))
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == WEEK {
		first := WeekFirst(start)
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * 7 * unit)

		// kstrings.Debug("datetime: {}, weekFirst: {}", start.Format(DATETIME_FORMATTER), Time.Unix(int64(first), 0).Format(DATETIME_FORMATTER))

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := WeekFirst(Time.Unix(int64(i), 0)), WeekLast(Time.Unix(int64(i+step-(1*60*60*24*7)), 0))

			// kstrings.Debug("datetime: {}, start: {}, end: {} - step: {}, days: {}", Time.Unix(int64(i), 0).Format(DATETIME_FORMATTER),
			// 	Time.Unix(int64(valBegin), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(valEnd), 0).Format(DATETIME_FORMATTER), int64(step), int64(step/60/60/24))

			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MONTH {
		first := MonthFirst(start)
		last := uint64(end.Unix())

		for i := first; i <= last; {
			curr := Time.Unix(int64(i), 0)
			valBegin := MonthFirst(curr)
			valEnd := MonthLast(curr.AddDate(0, int(unit-1), 0))

			nextTime := curr.AddDate(0, int(unit), 0)
			step := uint64(nextTime.Unix() - int64(i)) // 到下个月当前时刻的秒数

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			if i == first {
				valBegin = uint64(start.Unix())
			}

			i = i + step

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == YEAR {
		first := YearFirst(start)
		last := uint64(end.Unix())

		for i := first; i <= last; {
			curr := Time.Unix(int64(i), 0)
			valBegin := YearFirst(curr)
			valEnd := YearLast(curr.AddDate(int(unit-1), 0, 0))

			nextTime := curr.AddDate(int(unit), 0, 0)
			step := uint64(nextTime.Unix() - int64(i)) // 到明年当前时刻的秒数

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			if i == first {
				valBegin = uint64(start.Unix())
			}

			i = i + step

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	return klists.ToKSlice[*TimeSlice](list)
}

// 按指定时间间隔 进行非自然周期的切割
//
// 参数
//   - start: 开始时间
//   - end: 结束时间
//   - unit: 时间间隔单位
//   - duration: 时间间隔周期
//
// 返回值
//   - []*TimeSlice: 切割后的时间片段
func SplitUnnaturalDuration(start, end Time.Time, unit uint, duration Duration) []*TimeSlice {
	if unit == 0 {
		unit = 1
	}

	list := klists.New[*TimeSlice]()
	if duration == NONE {
		list.PushBack(&TimeSlice{Start: uint64(start.Unix()), End: uint64(end.Unix()), Tag: nil}) //  Records: klists.New[*RecordItem]()

		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == SECOND {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(unit)
		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(i), uint64(i+step-1)
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MINUTE {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(60 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(Time.Unix(int64(i), 0).Unix()), uint64(Time.Unix(int64(i+step-1), 0).Unix())
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == HOUR {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(60 * 60 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(Time.Unix(int64(i), 0).Unix()), uint64(Time.Unix(int64(i+step-1), 0).Unix())
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == DAY {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(Time.Unix(int64(i), 0).Unix()), uint64(Time.Unix(int64(i+step-1), 0).Unix())
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == WEEK {
		first := uint64(start.Unix())
		last := uint64(end.Unix())
		step := uint64(60 * 60 * 24 * 7 * unit)

		for i := first; i <= last; i = i + step {
			valBegin, valEnd := uint64(Time.Unix(int64(i), 0).Unix()), uint64(Time.Unix(int64(i+step-1), 0).Unix())

			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == MONTH {
		first := uint64(start.Unix())
		last := uint64(end.Unix())

		for i := first; i <= last; {
			curr := Time.Unix(int64(i), 0)
			nextTime := curr.AddDate(0, int(unit), 0)
			valBegin, valEnd := uint64(curr.Unix()), uint64(nextTime.Unix()-1)
			step := uint64(nextTime.Unix() - int64(i)) // 到下个月当前时刻的秒数
			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			i = i + step
			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	if duration == YEAR {
		first := uint64(start.Unix())
		last := uint64(end.Unix())

		for i := first; i <= last; {
			curr := Time.Unix(int64(i), 0)
			nextTime := curr.AddDate(int(unit), 0, 0)
			valBegin, valEnd := uint64(curr.Unix()), uint64(nextTime.Unix()-1)
			step := uint64(nextTime.Unix() - int64(i)) // 到下一年当前时刻的秒数

			if i == first {
				valBegin = uint64(start.Unix())
			}

			if i+step > last {
				valEnd = uint64(end.Unix())
			}

			i = i + step

			list.PushBack(&TimeSlice{Start: valBegin, End: valEnd, Tag: nil})
		}
		return klists.ToKSlice[*TimeSlice](list)
	}

	return klists.ToKSlice[*TimeSlice](list)
}
