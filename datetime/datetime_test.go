package datetime

import (
	"fmt"
	"testing"
	Time "time"
)

func Test_Minute(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastMinute(now)

	fmt.Printf("start Minute: %s\n", Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("end Minute: %s\n", Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")
}

func Test_Hour(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastHour(now)

	fmt.Printf("start Minute: %s\n", Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("end Minute: %s\n", Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")
}

func Test_Day(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastDay(now)

	fmt.Printf("start Minute: %s\n", Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("end Minute: %s\n", Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")
}

func Test_Week(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastWeek(now)

	fmt.Printf("datetime: %s, start Minute: %s, end Minute: %s\n", now.Format(DATETIME_FORMATTER), Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")

	strStart := "2024-01-14 12:14:12 +0800"

	startTime, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	start, end = FirstAndLastWeek(startTime)
	fmt.Printf("datetime: %s, start Minute: %s, end Minute: %s\n", startTime.Format(DATETIME_FORMATTER), Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))

	fmt.Printf("\n")
}

func Test_Month(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastMonth(now)

	fmt.Printf("start Minute: %s\n", Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("end Minute: %s\n", Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")
}

func Test_Year(t *testing.T) {
	now := Time.Now()
	start, end := FirstAndLastYear(now)

	fmt.Printf("start Minute: %s\n", Time.Unix(int64(start), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("end Minute: %s\n", Time.Unix(int64(end), 0).Format(DATETIME_FORMATTER))
	fmt.Printf("\n")
}

func Test_Datetime(t *testing.T) {
	fmt.Printf("local now: %d\n", Time.Now().Local().Unix())
	fmt.Printf("now: %d\n", Time.Now().Unix())

	current_time := Time.Now()
	timezone_name, timezone_offset := current_time.Zone()
	fmt.Printf("当前时区为%s，时间偏移量为%d秒\n", timezone_name, timezone_offset)
	fmt.Printf("\n")
}

func Test_SplitDurationByNone(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, NONE)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}

	fmt.Printf("\n\n")
}

func Test_SplitDurationBySecond(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, SECOND)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}

	fmt.Printf("\n\n")
}

func Test_SplitDurationByMinute(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 12:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, MINUTE)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}

	fmt.Printf("\n\n")
}

func Test_SplitDurationByHour(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-26 19:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, HOUR)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}
	fmt.Printf("\n\n")
}

func Test_SplitDurationByDay(t *testing.T) {
	strStart := "2024-01-26 12:14:12 +0800"
	strEnd := "2024-01-29 19:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, DAY)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}
	fmt.Printf("\n\n")
}

func Test_SplitDurationByWeek(t *testing.T) {
	strStart := "2024-01-14 12:14:12 +0800"
	strEnd := "2024-01-15 19:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, WEEK)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}
	fmt.Printf("\n\n")
}

func Test_SplitDurationByMonth(t *testing.T) {
	strStart := "2023-01-16 12:14:12 +0800"
	strEnd := "2024-02-23 19:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, MONTH)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}
	fmt.Printf("\n\n")
}

func Test_SplitDurationByYear(t *testing.T) {
	strStart := "2020-01-16 12:14:12 +0800"
	strEnd := "2024-01-01 19:20:10 +0800"

	start, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strStart, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}
	end, err := Time.ParseInLocation(DATETIME_TIMEZONE_FORMATTER, strEnd, Time.Local)
	if nil != err {
		t.Errorf("err: %s", err)
		return
	}

	fmt.Printf("test duration start: %s, end: %s\n\n", start.Format(DATETIME_FORMATTER), end.Format(DATETIME_FORMATTER))

	result := SplitDuration(start, end, YEAR)
	for _, item := range result {
		fmt.Printf(" start: %s, end: %s\n",
			Time.Unix(int64(item.Start), 0).Format(DATETIME_FORMATTER), Time.Unix(int64(item.End), 0).Format(DATETIME_FORMATTER))
	}
	fmt.Printf("\n\n")
}
