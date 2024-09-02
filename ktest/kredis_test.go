package ktest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/khan-lau/kutils/db/kredis"
	"github.com/khan-lau/kutils/logger"
)

var (
	redisHd *kredis.KRedis
	klog    *logger.Logger
)

func init() {
	klog = nil
	ctx := context.Background()
	redisHd = kredis.NewKRedis(ctx, "127.0.0.1", 16379, "", "WuTz@DtXyTeCh.com", 0)
	fmt.Println("init redis success!")
}

func Test(t *testing.T) {
	Test_RedisScan(t)
	t.Log("")
	Test_RedisScanMatch(t)
}

func Test_MatchArray(t *testing.T) {
	ret := kredis.MatchFilter([]string{"windRtEvent*", "*NSFC"}, "windRtEvent:DTNXJK:HSBFC:Q1:W009")
	if ret {
		t.Log("windRtEvent:DTNXJK:HSBFC:Q1:W009", "match success")
	} else {
		t.Log("windRtEvent:DTNXJK:HSBFC:Q1:W009", "match fail")
	}

	ret = kredis.MatchFilter([]string{"windRtEvent*", "*HSBFC"}, "windRtEvent:DTNXJK:HSBFC")
	if ret {
		t.Log("windRtEvent:DTNXJK:HSBFC", "match success")
	} else {
		t.Log("windRtEvent:DTNXJK:HSBFC", "match fail")
	}
}

func Test_RedisScan(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisScan begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		ret, err := redisHd.Scan(500, []string{"rejson", "rejson-rl", "hash", "list"}, []string{"JSON$DTHYJK"}, []string{"windRtEvent", "NSFC"}, false, nil)
		endTime := time.Now()
		t.Log("redisScan end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisScan cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("key count:", len(ret))
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_RedisScanMatch(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisScanMatch begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		ret, err := redisHd.ScanMatch(500, []string{"rejson", "rejson-rl", "hash", "list"}, []string{"JSON$DTHYJK"}, []string{"windRtEvent*", "*NSFC"}, false, LogFunc)
		endTime := time.Now()
		t.Log("redisScanMatch end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisScanMatch cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("key count:", len(ret))
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func LogFunc(lvl logger.Level, f string, args ...interface{}) {
	if nil == klog {
		klog = logger.LoggerInstanceOnlyConsole(int8(logger.DebugLevel))
		klog.Warrn("Not init logger")
	}
	skip := 1

	switch lvl {
	case logger.DebugLevel:
		klog.KDebug(skip, f, args...)
	case logger.InfoLevel:
		klog.KInfo(skip, f, args...)
	case logger.WarnLevel:
		klog.KWarrn(skip, f, args...)
	case logger.ErrorLevel:
		klog.KError(skip, f, args...)
	case logger.DPanicLevel:
		klog.KDPanic(skip, f, args...)
	case logger.FatalLevel:
		klog.KFatal(skip, f, args...)
	default:
		klog.KInfo(skip, fmt.Sprintf(lvl.String()+": "+f), args)
	}
}
