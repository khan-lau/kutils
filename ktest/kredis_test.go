package ktest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/khan-lau/kutils/db/kredis"
	"github.com/khan-lau/kutils/klogger"
)

var (
	redisHd *kredis.KRedis
	klog    *klogger.Logger
)

func init() {
	klog = nil
	ctx := context.Background()
	redisHd = kredis.NewKRedis(ctx, "10.50.145.10", 16379, "", "WuTz@DtXyTeCh.com", 0)
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

func Test_JsonSet(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonSet begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		err := redisHd.JsonSet("testJson", ".", `{"name":"khan","age":30}`)
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonGet(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonSet begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		r, err := redisHd.JsonGet("testJson", ".")
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("json:", r)
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonType(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonType begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		r, err := redisHd.JsonType("testJson", "$.name")
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("json type:", strings.Join(r, ","))
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonObjKeys(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonObjKeys begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		r, err := redisHd.JsonObjKeys("testJson", ".")
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("json keys:", strings.Join(r, ","))
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonDel(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonDel begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		r, err := redisHd.JsonDel("testJson", "$.age")
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Logf("json del count %d keys", r)
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonMerge(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonMerge begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		// err := redisHd.JsonMerge("testJson", ".", `{"age":30}`)
		err := redisHd.JsonMerge("testJson", ".", `{"another":{}}`)
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func Test_JsonObjLen(t *testing.T) {
	if redisHd.Ping() {
		beginTime := time.Now()
		t.Log("redisJsonObjLen begin ", beginTime.Format("2006-01-02 15:04:05.000"))
		r, err := redisHd.JsonObjLen("testJson", "$.another")
		// r, err := redisHd.JsonObjLen("docJson", "$..a")
		endTime := time.Now()
		t.Log("redisJsonSet end", endTime.Format("2006-01-02 15:04:05.000"))
		t.Log("redisJsonSet cost ", endTime.Sub(beginTime).Milliseconds())
		if err != nil {
			t.Fatal(err)
		} else {
			t.Logf("json obj len %d keys", r)
		}
	} else {
		t.Fatal("redis not connect!")
	}
}

func LogFunc(lvl klogger.Level, f string, args ...interface{}) {
	if nil == klog {
		klog = klogger.LoggerInstanceOnlyConsole(int8(klogger.DebugLevel))
		klog.Warrn("Not init logger")
	}
	skip := 1

	switch lvl {
	case klogger.DebugLevel:
		klog.KDebug(skip, f, args...)
	case klogger.InfoLevel:
		klog.KInfo(skip, f, args...)
	case klogger.WarnLevel:
		klog.KWarrn(skip, f, args...)
	case klogger.ErrorLevel:
		klog.KError(skip, f, args...)
	case klogger.DPanicLevel:
		klog.KDPanic(skip, f, args...)
	case klogger.FatalLevel:
		klog.KFatal(skip, f, args...)
	default:
		klog.KInfo(skip, fmt.Sprintf(lvl.String()+": "+f), args)
	}
}
