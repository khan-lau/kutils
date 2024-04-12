package ktest

import (
	"testing"
	Time "time"

	"github.com/khan-lau/kutils/container/kstrings"
	"github.com/khan-lau/kutils/datetime"
)

func Test_Minute(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastMinute(now)

	kstrings.Print("start Minute: {}", Time.Unix(int64(start), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("end Minute: {}", Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")
}

func Test_Hour(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastHour(now)

	kstrings.Print("start Hour: {}", Time.Unix(int64(start), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("end Hour: {}", Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")
}

func Test_Day(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastDay(now)

	kstrings.Print("start Day: {}", Time.Unix(int64(start), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("end Day: {}", Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")
}

func Test_Week(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastWeek(now)
	kstrings.Print("datetime: {}, start Week: {}, end Week: {}", now.Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(start), 0).
		Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")

	start, end = datetime.WeekFirst(now), datetime.WeekLast(now)
	kstrings.Print("datetime: {}, start Week: {}, end Week: {}", now.Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(start), 0).
		Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")

	strStart := "2024-01-14 12:14:12 +0800"

	startTime, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	start, end = datetime.FirstAndLastWeek(startTime)
	kstrings.Print("datetime: {}, start Week: {}, end Week: {}", startTime.Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(start), 0).
		Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")

	start, end = datetime.WeekFirst(startTime), datetime.WeekLast(startTime)
	kstrings.Print("datetime: {}, start Week: {}, end Week: {}", startTime.Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(start), 0).
		Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")

	kstrings.Print("")
}

func Test_Month(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastMonth(now)

	kstrings.Print("start Month: {}", Time.Unix(int64(start), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("end Month: {}", Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")
}

func Test_Year(t *testing.T) {
	now := Time.Now()
	start, end := datetime.FirstAndLastYear(now)

	kstrings.Print("start Year: {}", Time.Unix(int64(start), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("end Year: {}", Time.Unix(int64(end), 0).Format(datetime.DATETIME_FORMATTER))
	kstrings.Print("")
}

func Test_Datetime(t *testing.T) {
	kstrings.Print("local now: {}", Time.Now().Local().Unix())
	kstrings.Print("now: {}", Time.Now().Unix())

	current_time := Time.Now()
	timezone_name, timezone_offset := current_time.Zone()
	kstrings.Print("当前时区为{}, 时间偏移量为{}秒\n", timezone_name, timezone_offset)
	kstrings.Print("")
}

///////////////////////////////////////////////////////////////////////

func Test_SplitDurationByNone(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.NONE)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("")
}

func Test_SplitDurationBySecond(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.SECOND)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("end - start = {}", end.Unix()-start.Unix()+1)
	kstrings.Print("duration slice size: {}", len(result))

	kstrings.Print("\n")
}

func Test_SplitDurationByMinute(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.MINUTE)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitDurationByHour(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.HOUR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitDurationByDay(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-29 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.DAY)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitDurationByWeek(t *testing.T) {
	strStart := "2024-01-14 12:14:12 +0800"
	strEnd := "2024-01-15 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.WEEK)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitDurationByMonth(t *testing.T) {
	strStart := "2023-01-16 12:14:12 +0800"
	strEnd := "2024-02-23 19:20:10 +0800"
	// strEnd := "2023-02-23 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.MONTH)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Printf("\n")
}

func Test_SplitDurationByYear(t *testing.T) {
	strStart := "2020-01-16 12:14:12 +0800"
	strEnd := "2024-01-01 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitDuration(start, end, datetime.YEAR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

///////////////////////////////////////////////////////////////////////

func Test_SplitNaturalDurationNone(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitNaturalDuration(start, end, 5, datetime.NONE) // unit 只能为1, 为其他值没有任何意义
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("")
}

func Test_SplitNaturalDurationBySecond(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"
	// strEnd := "2024-01-26 12:14:13 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(5)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.SECOND)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("end - start = {}", end.Unix()-start.Unix()+1)
	kstrings.Print("duration slice size: {}", len(result))

	kstrings.Print("\n")
}

func Test_SplitNaturalDurationByMinute(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:47:10 +0800"
	// strEnd := "2024-01-26 12:15:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(5)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.MINUTE)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitNaturalDurationByHour(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.HOUR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitNaturalDurationByDay(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-31 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.DAY)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitNaturalDurationByWeek(t *testing.T) {
	strStart := "2024-01-14 12:14:12 +0800"
	strEnd := "2024-01-22 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.WEEK)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitNaturalDurationByMonth(t *testing.T) {
	strStart := "2023-01-16 12:14:12 +0800"
	strEnd := "2024-02-23 19:20:10 +0800"
	// strEnd := "2023-03-21 10:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.MONTH)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Printf("\n")
}

func Test_SplitNaturalDurationByYear(t *testing.T) {
	strStart := "2020-01-16 12:14:12 +0800"
	strEnd := "2024-01-01 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitNaturalDuration(start, end, unit, datetime.YEAR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

///////////////////////////////////////////////////////////////////////

func Test_SplitUnNaturalDurationNone(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Print("test duration start: {}, end: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER))

	result := datetime.SplitUnnaturalDuration(start, end, 5, datetime.NONE) // unit 只能为1, 为其他值没有任何意义
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("")
}

func Test_SplitUnNaturalDurationBySecond(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"
	// strEnd := "2024-01-26 12:14:13 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(5)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.SECOND)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("end - start = {}", end.Unix()-start.Unix()+1)
	kstrings.Print("duration slice size: {}", len(result))

	kstrings.Print("\n")
}

func Test_SplitUnNaturalDurationByMinute(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:47:10 +0800"
	// strEnd := "2024-01-26 12:15:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(5)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.MINUTE)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitUnNaturalDurationByHour(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.HOUR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitUnNaturalDurationByDay(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-31 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.DAY)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitUnNaturalDurationByWeek(t *testing.T) {
	strStart := "2024-01-14 12:14:12 +0800"
	strEnd := "2024-02-22 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.WEEK)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

func Test_SplitUnNaturalDurationByMonth(t *testing.T) {
	strStart := "2023-01-16 12:14:12 +0800"
	strEnd := "2024-02-23 19:20:10 +0800"
	// strEnd := "2023-03-21 10:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.MONTH)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Printf("\n")
}

func Test_SplitUnNaturalDurationByYear(t *testing.T) {
	strStart := "2020-01-16 12:14:12 +0800"
	strEnd := "2024-01-01 19:20:10 +0800"

	start, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	unit := uint(2)
	kstrings.Print("test duration start: {}, end: {}, unit: {}", start.Format(datetime.DATETIME_FORMATTER), end.Format(datetime.DATETIME_FORMATTER), int(unit))

	result := datetime.SplitUnnaturalDuration(start, end, unit, datetime.YEAR)
	for _, item := range result {
		kstrings.Print(" start: {}, end: {}",
			Time.Unix(int64(item.Start), 0).Format(datetime.DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(datetime.DATETIME_FORMATTER))
	}

	kstrings.Print("duration slice size: {}", len(result))
	kstrings.Print("\n")
}

///////////////////////////////////////////////////////////////////////

func Test_Duration(t *testing.T) {
	startStr1 := "2020-01-26 12:14:12 +0800"
	endStr1 := "2021-01-26 12:20:10 +0800"

	startStr2 := "2021-01-26 12:14:12 +0800"
	endStr2 := "2022-01-26 12:20:10 +0800"

	startStr3 := "2022-01-26 12:14:12 +0800"
	endStr3 := "2023-01-26 12:20:10 +0800"

	startStr4 := "2023-01-26 12:14:12 +0800"
	endStr4 := "2024-01-26 12:20:10 +0800"

	star1, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, startStr1, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	end1, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, endStr1, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	star2, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, startStr2, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	end2, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, endStr2, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	star3, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, startStr3, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	end3, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, endStr3, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	star4, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, startStr4, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	end4, err := Time.ParseInLocation(datetime.DATETIME_TIMEZONE_FORMATTER, endStr4, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	kstrings.Debug("2020 duration: {}, days:{}", end1.Unix()-star1.Unix(), (end1.Unix()-star1.Unix())/24/60/60)
	kstrings.Debug("2021 duration: {}, days:{}", end2.Unix()-star2.Unix(), (end2.Unix()-star2.Unix())/24/60/60)
	kstrings.Debug("2022 duration: {}, days:{}", end3.Unix()-star3.Unix(), (end3.Unix()-star3.Unix())/24/60/60)
	kstrings.Debug("2023 duration: {}, days:{}", end4.Unix()-star4.Unix(), (end4.Unix()-star4.Unix())/24/60/60)
}
