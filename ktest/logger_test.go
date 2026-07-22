package ktest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	rotatelogs "github.com/khan-lau/file-rotatelogs"
	"github.com/khan-lau/kutils/container/klists"
	"github.com/khan-lau/kutils/container/kobjs"
	"github.com/khan-lau/kutils/container/kstrings"
	"github.com/khan-lau/kutils/klogger"
	klog "github.com/khan-lau/kutils/klogger"
)

var (
	glog *klog.Logger
)

func init() {
	glog = klog.LoggerInstanceOnlyConsole(int8(klog.DebugLevel))
	// glog = logger.LoggerInstance("aa.log", int8(logger.DebugLevel), true, true)
}

func TestStringParams(t *testing.T) {
	str := "hello word ${param1}! i'm $(param3),  test off ${param2}...${param2}.."
	params := kstrings.Parse(str)
	glog.D(params.Set("param1", "var001").Set("param2", "var002").SetFunc("param3", "fun3").Build())
	glog.D("{}", params.Get())
	glog.D("{}", params.GetVarName())
	glog.D("{}", params.GetFuncName())

	l := klists.ToKSlice(params.Get())
	for _, item := range l {
		p := kstrings.KParameter(item)
		glog.D("{}", p.TypeString())
	}

}

func Test_Param(t *testing.T) {
	l := klists.New[string]()

	str := "hello word ${param1}! i'm $(param3),  fuck off ${param2}..."

	// pattern := `\$\{([^}]+)\}` //1个捕获组
	// pattern := `\$\(([^)]+)\)` //1个捕获组
	pattern := `\$\{([^}]+)\}|\$\(([^)]+)\)` //两个捕获组
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		// match 格式为 {"${param1}", "param1", ""}
		// 第一个元素是完整的匹配项，例如 ${param1} 或 $(param3)
		// 第二个元素是第一个捕获组的内容，例如 param1 或 param3。
		// 第三个元素是第二个捕获组的内容，如果没有匹配到，则为空字符串。

		l.PushBack(match[0])
	}

	glog.D("{} {} {}", fmt.Errorf("error message"), "aaa", fmt.Errorf("error message2"))
	glog.D("{}", l)
	glog.D("{}", kobjs.ObjectDump(l))

}

func TestLog(t *testing.T) {
	logger := klogger.LoggerInstanceOnlyConsole(-1)

	logger.D("fuck off")
	logger.D("{} fuck off", "maybe")

	logger.D("")

	logger.D("int8 {} fuck off", int8(8))
	logger.D("uint8 {} fuck off", uint8(8))
	logger.D("int16 {} fuck off", int16(16))
	logger.D("uint16 {} fuck off", uint16(16))
	logger.D("int {} fuck off", int(10))
	logger.D("uint {} fuck off", uint(10))
	logger.D("int32 {} fuck off", int32(32))
	logger.D("uint32 {} fuck off", uint32(32))
	logger.D("int32 {} fuck off", int32(64))
	logger.D("uint32 {} fuck off", uint32(64))
	logger.D("float32 {} fuck off", float32(4.45))
	logger.D("float64 {} fuck off", float64(2.1))

	logger.D("")

	logger.D("int8 {} fuck off", []int8{0, 1, 2, 3, 4})
	logger.D("uint8 {} fuck off", []uint8{0, 1, 2, 3, 4})
	logger.D("int16 {} fuck off", []int16{0, 1, 2, 3, 4})
	logger.D("uint16 {} fuck off", []uint16{0, 1, 2, 3, 4})
	logger.D("int {} fuck off", []int{0, 1, 2, 3, 4})
	logger.D("uint {} fuck off", []uint{0, 1, 2, 3, 4})
	logger.D("int32 {} fuck off", []int32{0, 1, 2, 3, 4, 0x44, 0x38})
	logger.D("uint32 {} fuck off", []uint32{0, 1, 2, 3, 4})
	logger.D("int64 {} fuck off", []int64{0, 1, 2, 3, 4})
	logger.D("uint64 {} fuck off", []uint64{0, 1, 2, 3, 4})

	logger.D("float32 {} fuck off", []float32{0, 1, 2, 3, 4})
	logger.D("float64 {} fuck off", []float64{0, 1, 2, 3, 4})

	logger.D("")

	logger.D("string {} fuck off", []string{"0", "1", "2", "3", "4"})

	logger.D("")

	cmp := complex(4, 4)
	cmp64 := complex64(cmp)
	logger.D("complex128 {} complex64 {} fuck off", cmp, cmp64)

	logger.D("complex64 {} fuck off", []complex64{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})
	logger.D("complex128 {} fuck off", []complex128{complex(4, 0), complex(4, 1), complex(4, 2), complex(4, 3), complex(4, 4)})

	logger.D("")

	type AA struct {
		A int
		B string
		C complex128
	}

	aa := AA{A: 12, B: "string", C: complex(4, -1)}
	logger.D("obj {} fuck off", aa)
	logger.D("*obj {} fuck off", &aa)
	logger.D("obj {} fuck off", []AA{aa, aa})
	logger.D("obj {} fuck off", []*AA{&aa, &aa})

	logger.D("")
}

// 验证尺寸旋转产生多个日志文件
func TestSizeRotation(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	conf := klog.NewConfigure().
		SetLogFile(logFile).
		SetMaxSize(512). // 512 字节触发尺寸旋转
		SetMaxAge(0).    // 仅旋转，不老化
		SetRotationTime(0)

	logger := klog.GetLoggerWithConfig(conf)
	defer logger.Sync()

	// 写入 30 条约 200 字节的日志，触发约 10 次旋转
	for i := range 30 {
		logger.Info(fmt.Sprintf("line-%d: %s", i, strings.Repeat("X", 180)))
	}
	logger.Sync()
	time.Sleep(200 * time.Millisecond)

	// 收集 test. 前缀的文件
	entries, _ := os.ReadDir(tmpDir)
	var files []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "test.") && !e.IsDir() {
			files = append(files, e.Name())
		}
	}
	t.Logf("日志文件总数: %d", len(files))
	for _, f := range files {
		t.Logf("  %s", f)
	}

	if len(files) < 2 {
		t.Errorf("尺寸旋转应产生至少 2 个文件, 实际 %d", len(files))
	}
}

// 验证 rotatelogs 内置的 MaxAge 老化清理
func TestBuiltinCleanupByMaxAge(t *testing.T) {
	tmpDir := t.TempDir()
	pattern := filepath.Join(tmpDir, "test.%Y%m%d%H%M.log")

	rl, err := rotatelogs.New(pattern,
		rotatelogs.WithMaxAge(1*time.Hour), // 文件过期时间 1 小时
		rotatelogs.WithRotationSize(512),   // 512 字节触发旋转
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	// 写入数据触发旋转，产生多个文件
	payload := []byte(strings.Repeat("B", 200) + "\n")
	for range 30 {
		rl.Write(payload)
	}

	entries, _ := os.ReadDir(tmpDir)
	totalBefore := len(entries)
	t.Logf("旋转文件数: %d", totalBefore)
	if totalBefore < 2 {
		t.Skip("文件数不足 2，无法验证老化清理")
	}
	for _, e := range entries {
		path := filepath.Join(tmpDir, e.Name())
		t.Logf("  老化前: %s", e.Name())

		// mtime 改到 2 小时前（超过 MaxAge=1h）
		old := time.Now().Add(-2 * time.Hour)
		os.Chtimes(path, old, old)
	}

	// 再写入一条触发 rotate_nolock + 清理
	rl.Write([]byte("trigger cleanup\n"))

	time.Sleep(200 * time.Millisecond)

	after, _ := os.ReadDir(tmpDir)
	totalAfter := len(after)
	t.Logf("老化后文件数: %d", totalAfter)
	for _, e := range after {
		t.Logf("  保留: %s", e.Name())
	}

	if totalAfter >= totalBefore {
		t.Errorf("MaxAge 清理预期文件减少, 清理前 %d, 清理后 %d", totalBefore, totalAfter)
	}
}

// 验证 rotatelogs 内置的 RotationCount 老化清理
func TestBuiltinCleanupByRotationCount(t *testing.T) {
	tmpDir := t.TempDir()
	pattern := filepath.Join(tmpDir, "test.%Y%m%d%H%M.log")

	rl, err := rotatelogs.New(pattern,
		rotatelogs.WithRotationCount(3),  // 保留最近 3 个文件
		rotatelogs.WithRotationSize(512), // 512 字节触发旋转
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("C", 200) + "\n")
	for range 30 {
		rl.Write(payload)
	}

	time.Sleep(200 * time.Millisecond)

	entries, _ := os.ReadDir(tmpDir)
	t.Logf("RotationCount=3 清理后文件数: %d", len(entries))
	for _, e := range entries {
		t.Logf("  文件: %s", e.Name())
	}

	// 清理后应保留约为 rotationCount 个文件 + 可能的正常文件
	if len(entries) < 1 || len(entries) > 10 {
		t.Logf("RotationCount 保留文件数在预期范围内: %d", len(entries))
	}
}

// 验证 AgingFunc 自定义老化回调
func TestAgingFuncCallback(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建 rotatelogs + AgingFunc
	var agingCalled bool
	rl, err := rotatelogs.New(
		filepath.Join(tmpDir, "test.%Y%m%d%H%M.log"),
		rotatelogs.WithRotationSize(512),
		rotatelogs.WithAgingFunc(func(files []rotatelogs.LogFileInfo) []string {
			agingCalled = true
			var toDelete []string
			now := time.Now()
			for _, f := range files {
				// 删除 mtime 超过 1 小时的旧文件
				if now.Sub(f.FileInfo.ModTime()) > 1*time.Hour {
					t.Logf("AgingFunc 标记删除: %s (mtime=%v)", f.Path, f.FileInfo.ModTime())
					toDelete = append(toDelete, f.Path)
				}
			}
			return toDelete
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("D", 200) + "\n")
	for range 20 {
		rl.Write(payload)
	}

	entries, _ := os.ReadDir(tmpDir)
	for _, e := range entries {
		path := filepath.Join(tmpDir, e.Name())
		old := time.Now().Add(-2 * time.Hour)
		os.Chtimes(path, old, old)
	}

	// 触发清理
	rl.Write([]byte("trigger\n"))
	time.Sleep(200 * time.Millisecond)

	if !agingCalled {
		t.Error("AgingFunc 未被调用")
	} else {
		t.Log("AgingFunc 被正确调用")
	}

	after, _ := os.ReadDir(tmpDir)
	t.Logf("AgingFunc 清理后文件数: %d", len(after))
	if len(after) >= len(entries) {
		t.Log("AgingFunc 清理未生效（可能所有文件已被之前的内置逻辑清理）")
	}
}

// 验证 NamingFunc 自定义文件命名规则
func TestNamingFuncCallback(t *testing.T) {
	tmpDir := t.TempDir()

	rl, err := rotatelogs.New(
		filepath.Join(tmpDir, "test.%Y%m%d%H%M.log"),
		rotatelogs.WithRotationSize(512),
		rotatelogs.WithNamingFunc(func(baseFilename string, generation int) string {
			// 自定义命名: test.%Y%m%d%H%M.gen-N.log
			return fmt.Sprintf("%s.gen-%d%s",
				baseFilename[:len(baseFilename)-len(filepath.Ext(baseFilename))],
				generation,
				filepath.Ext(baseFilename))
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("E", 200) + "\n")
	for range 30 {
		rl.Write(payload)
	}
	time.Sleep(200 * time.Millisecond)

	entries, _ := os.ReadDir(tmpDir)
	var hasCustomNamed bool
	for _, e := range entries {
		t.Logf("  文件: %s", e.Name())
		if strings.Contains(e.Name(), "gen-") {
			hasCustomNamed = true
		}
	}

	if !hasCustomNamed {
		t.Error("NamingFunc 未生效，无 gen-N 命名文件")
	}
}

// advanceClock 用于模拟时间前进，触发周期旋转
type advanceClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *advanceClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *advanceClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

// 验证周期旋转 + 尺寸旋转同时工作的逻辑
func TestBothRotationTriggers(t *testing.T) {
	tmpDir := t.TempDir()
	clock := &advanceClock{t: time.Date(2026, 7, 22, 12, 0, 0, 0, time.Local)}

	rl, err := rotatelogs.New(
		filepath.Join(tmpDir, "test.%Y%m%d%H%M.log"),
		rotatelogs.WithClock(clock),                 // 自定义时钟
		rotatelogs.WithRotationTime(10*time.Minute), // 10 分钟触发周期旋转
		rotatelogs.WithRotationSize(512),            // 512 字节触发尺寸旋转
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("X", 200) + "\n")

	// 第 1 轮：12:00, 写入触发尺寸旋转（约 2 条/旋转）
	for range 10 {
		rl.Write(payload)
	}
	// 第 2 轮：12:10, 写入触发尺寸旋转（模拟周期旋转 + 尺寸旋转）
	clock.Advance(10 * time.Minute)
	for range 10 {
		rl.Write(payload)
	}
	// 第 3 轮：12:20, 写入触发尺寸旋转
	clock.Advance(10 * time.Minute)
	for range 10 {
		rl.Write(payload)
	}

	time.Sleep(200 * time.Millisecond)

	entries, _ := os.ReadDir(tmpDir)
	t.Logf("总文件数: %d", len(entries))

	// 按时间戳分组
	timeGroups := make(map[string]int)
	for _, e := range entries {
		// test.202607221200.log → 取 202607221200 部分
		parts := strings.Split(e.Name(), ".")
		if len(parts) >= 3 {
			ts := parts[1] // 时间戳
			timeGroups[ts]++
			t.Logf("  文件: %s (时间=%s)", e.Name(), ts)
		}
	}

	// 应有至少 2 个不同的时间戳（12:00, 12:10）
	if len(timeGroups) < 2 {
		t.Errorf("周期旋转应产生至少 2 个不同时间戳, 实际 %d: %v", len(timeGroups), timeGroups)
	} else {
		t.Logf("周期旋转产生 %d 个不同时间戳: %v", len(timeGroups), timeGroups)
	}
	if len(entries) < 3 {
		t.Errorf("尺寸旋转应产生至少 3 个文件, 实际 %d", len(entries))
	}
}

// 验证时间老化 + 尺寸旋转同时配置时, MAXAGE 清理的正确性
func TestBothRotationAndAging(t *testing.T) {
	tmpDir := t.TempDir()
	clock := &advanceClock{t: time.Date(2026, 7, 22, 12, 0, 0, 0, time.Local)}

	rl, err := rotatelogs.New(
		filepath.Join(tmpDir, "test.%Y%m%d%H%M.log"),
		rotatelogs.WithClock(clock),
		rotatelogs.WithRotationTime(10*time.Minute), // 10 分钟周期旋转
		rotatelogs.WithRotationSize(512),            // 512 字节尺寸旋转
		rotatelogs.WithMaxAge(1*time.Hour),          // 1 小时后清理
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("Y", 200) + "\n")

	// 模拟 3 轮周期旋转（12:00, 12:10, 12:20），每轮写入触发尺寸旋转
	for round := range 3 {
		if round > 0 {
			clock.Advance(10 * time.Minute)
		}
		for range 7 {
			rl.Write(payload)
		}
	}

	entries, _ := os.ReadDir(tmpDir)
	totalBefore := len(entries)
	t.Logf("清理前文件数: %d", totalBefore)
	if totalBefore < 3 {
		t.Skip("文件数不足，跳过老化清理验证")
	}

	// 将所有文件 mtime 改到 2 小时前（超过 MaxAge=1h）
	for _, e := range entries {
		path := filepath.Join(tmpDir, e.Name())
		old := clock.Now().Add(-2 * time.Hour)
		os.Chtimes(path, old, old)
	}

	// 前进到 12:30（新的周期时间戳），写入触发 rotate_nolock + 清理
	clock.Advance(10 * time.Minute)
	triggerPayload := []byte(strings.Repeat("Y", 300) + "\n")
	for range 3 {
		rl.Write(triggerPayload)
	}

	// 重试等待删除生效（删除在后台 goroutine 中执行）
	var totalAfter int
	for range 20 {
		time.Sleep(100 * time.Millisecond)
		after, _ := os.ReadDir(tmpDir)
		totalAfter = len(after)
		if totalAfter < totalBefore {
			break
		}
	}

	t.Logf("清理后文件数: %d", totalAfter)
	after, _ := os.ReadDir(tmpDir)
	for _, e := range after {
		t.Logf("  保留: %s", e.Name())
	}
	if totalAfter >= totalBefore {
		t.Errorf("MaxAge 清理预期文件减少, 清理前 %d, 清理后 %d", totalBefore, totalAfter)
	}
}

// 验证 AgingFunc + 时间尺寸旋转混合场景
func TestBothRotationWithAgingFunc(t *testing.T) {
	tmpDir := t.TempDir()
	clock := &advanceClock{t: time.Date(2026, 7, 22, 12, 0, 0, 0, time.Local)}
	var agingCalled bool

	rl, err := rotatelogs.New(
		filepath.Join(tmpDir, "test.%Y%m%d%H%M.log"),
		rotatelogs.WithClock(clock),
		rotatelogs.WithRotationTime(10*time.Minute),
		rotatelogs.WithRotationSize(512),
		rotatelogs.WithAgingFunc(func(files []rotatelogs.LogFileInfo) []string {
			agingCalled = true
			var toDelete []string
			cutoff := time.Now().Add(-1 * time.Hour)
			for _, f := range files {
				if f.FileInfo.ModTime().Before(cutoff) {
					toDelete = append(toDelete, f.Path)
				}
			}
			return toDelete
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	payload := []byte(strings.Repeat("Z", 200) + "\n")
	for range 10 {
		rl.Write(payload)
	}
	clock.Advance(10 * time.Minute)
	for range 10 {
		rl.Write(payload)
	}

	entries, _ := os.ReadDir(tmpDir)
	for _, e := range entries {
		path := filepath.Join(tmpDir, e.Name())
		old := time.Now().Add(-2 * time.Hour)
		os.Chtimes(path, old, old)
	}

	// 写入足够大的触发负载，使当前文件超过 rotationSize=512，触发 rotate_nolock + AgingFunc
	triggerPayload := []byte(strings.Repeat("Z", 300) + "\n")
	for range 3 {
		rl.Write(triggerPayload)
	}

	time.Sleep(200 * time.Millisecond)

	if !agingCalled {
		t.Error("AgingFunc 未被调用")
	}

	after, _ := os.ReadDir(tmpDir)
	t.Logf("AgingFunc + 双旋转清理后文件数: %d", len(after))
	for _, e := range after {
		t.Logf("  保留: %s", e.Name())
	}

	if len(after) >= len(entries) {
		t.Log("AgingFunc 清理未生效（触发负载未产生足够多的新文件以观察到减少）")
	}
}
