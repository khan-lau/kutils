package ktest_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/khan-lau/kutils/container/queue"
)

// ============================================================
// 辅助函数
// ============================================================

func newTestQueue(tb testing.TB, name string, maxBytes int64) queue.Interface {
	tb.Helper()
	dataPath := filepath.Join(os.TempDir(), "diskq-test", fmt.Sprintf("%s-%d", name, time.Now().UnixNano()))
	q, err := queue.New(name, dataPath, maxBytes, 1, 1024*1024, 2500, 2*time.Second, nil)
	if err != nil {
		tb.Fatalf("New() failed: %v", err)
	}
	tb.Cleanup(func() {
		q.Close()
		os.RemoveAll(dataPath)
	})
	return q
}

func putMessage(tb testing.TB, q queue.Interface, msg []byte) {
	tb.Helper()
	if err := q.Put(msg); err != nil {
		tb.Fatalf("Put() failed: %v", err)
	}
}

// ============================================================
// 功能测试
// ============================================================

func TestPutAndRead(t *testing.T) {
	q := newTestQueue(t, "TestPutAndRead", 1024*1024)
	msg := []byte("hello diskqueue")

	putMessage(t, q, msg)

	select {
	case data := <-q.ReadChan():
		if string(data) != string(msg) {
			t.Fatalf("ReadChan() = %q, want %q", data, msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ReadChan() timed out waiting for message")
	}
}

func TestPutAndReadMultiple(t *testing.T) {
	q := newTestQueue(t, "TestPutAndReadMultiple", 1024*1024)
	n := 100
	msgs := make([][]byte, n)
	for i := range n {
		msgs[i] = []byte(fmt.Sprintf("msg-%d", i))
		putMessage(t, q, msgs[i])
	}

	for i := range n {
		select {
		case data := <-q.ReadChan():
			if string(data) != string(msgs[i]) {
				t.Fatalf("ReadChan()[%d] = %q, want %q", i, data, msgs[i])
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("ReadChan()[%d] timed out", i)
		}
	}

	// 确认队列已空
	if depth := q.Depth(); depth != 0 {
		t.Fatalf("Depth() = %d, want 0", depth)
	}
}

func TestPeekChan(t *testing.T) {
	q := newTestQueue(t, "TestPeekChan", 1024*1024)
	msg := []byte("peek test")

	putMessage(t, q, msg)

	// Peek 两次，应该都能读到相同数据
	for i := range 2 {
		select {
		case data := <-q.PeekChan():
			if string(data) != string(msg) {
				t.Fatalf("PeekChan()[%d] = %q, want %q", i, data, msg)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("PeekChan()[%d] timed out", i)
		}
	}

	// ReadChan 消费后队列深度应为 0
	select {
	case data := <-q.ReadChan():
		if string(data) != string(msg) {
			t.Fatalf("ReadChan() = %q, want %q", data, msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ReadChan() timed out")
	}

	if depth := q.Depth(); depth != 0 {
		t.Fatalf("Depth() = %d, want 0", depth)
	}
}

func TestEmpty(t *testing.T) {
	q := newTestQueue(t, "TestEmpty", 1024*1024)

	for i := range 50 {
		putMessage(t, q, []byte(fmt.Sprintf("msg-%d", i)))
	}

	if err := q.Empty(); err != nil {
		t.Fatalf("Empty() failed: %v", err)
	}

	if depth := q.Depth(); depth != 0 {
		t.Fatalf("Depth() after Empty = %d, want 0", depth)
	}

	// 确认 Empty 后写入仍正常
	putMessage(t, q, []byte("after empty"))
	select {
	case data := <-q.ReadChan():
		if string(data) != "after empty" {
			t.Fatalf("ReadChan() after Empty = %q, want %q", data, "after empty")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ReadChan() after Empty timed out")
	}
}

func TestCloseAndReopen(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "diskq-test", fmt.Sprintf("TestCloseAndReopen-%d", time.Now().UnixNano()))
	t.Cleanup(func() { os.RemoveAll(dataPath) })

	// 创建队列并写入数据
	q1, err := queue.New("persist", dataPath, 1024*1024, 1, 1024*1024, 2500, 2*time.Second, nil)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	msgs := []string{"alpha", "beta", "gamma"}
	for _, m := range msgs {
		if err = q1.Put([]byte(m)); err != nil {
			t.Fatalf("Put(%q) failed: %v", m, err)
		}
	}

	// 先消费一条然后关闭，验证状态持久化
	select {
	case data := <-q1.ReadChan():
		if string(data) != msgs[0] {
			t.Fatalf("ReadChan() = %q, want %q", data, msgs[0])
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ReadChan() timed out")
	}

	if err = q1.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// 重新打开，验证剩余数据
	q2, err := queue.New("persist", dataPath, 1024*1024, 1, 1024*1024, 2500, 2*time.Second, nil)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer q2.Close()

	for i := 1; i < len(msgs); i++ {
		select {
		case data := <-q2.ReadChan():
			if string(data) != msgs[i] {
				t.Fatalf("ReadChan() after reopen = %q, want %q", data, msgs[i])
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("ReadChan() after reopen timed out for msg %q", msgs[i])
		}
	}
}

func TestDelete(t *testing.T) {
	dataPath := filepath.Join(os.TempDir(), "diskq-test", fmt.Sprintf("TestDelete-%d", time.Now().UnixNano()))
	t.Cleanup(func() { os.RemoveAll(dataPath) })

	q, err := queue.New("delete", dataPath, 1024*1024, 1, 1024*1024, 2500, 2*time.Second, nil)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	putMessage(t, q, []byte("to be deleted"))

	if err := q.Delete(); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// 确认数据目录已被清理或不包含 meta 文件
	metaFile := filepath.Join(dataPath, "delete.diskqueue.meta.dat")
	if _, err := os.Stat(metaFile); !os.IsNotExist(err) {
		t.Fatalf("meta file still exists after Delete(): %v", err)
	}
}

func TestFileRolling(t *testing.T) {
	// maxBytesPerFile=4096 强制文件滚动
	q := newTestQueue(t, "TestFileRolling", 4096)
	msg := make([]byte, 512)

	for range 50 {
		putMessage(t, q, msg)
	}

	// 消费所有消息，验证顺序一致
	for i := range 50 {
		select {
		case data := <-q.ReadChan():
			if len(data) != 512 {
				t.Fatalf("ReadChan()[%d] len = %d, want 512", i, len(data))
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("ReadChan()[%d] timed out", i)
		}
	}
}

func TestMaxMsgSizeBoundary(t *testing.T) {
	q := newTestQueue(t, "TestMaxMsgSizeBoundary", 1024*1024)

	// 空消息（minMsgSize=1，空消息应失败）
	if err := q.Put([]byte{}); err == nil {
		t.Fatal("Put(empty) should have failed")
	}

	// 合法消息
	if err := q.Put([]byte("x")); err != nil {
		t.Fatalf("Put(min size) failed: %v", err)
	}

	select {
	case data := <-q.ReadChan():
		if string(data) != "x" {
			t.Fatalf("ReadChan() = %q, want %q", data, "x")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ReadChan() timed out")
	}
}

func TestConcurrentPut(t *testing.T) {
	q := newTestQueue(t, "TestConcurrentPut", 1024*1024)

	var wg sync.WaitGroup
	n := 5
	msgsPerG := 200

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range msgsPerG {
				if err := q.Put(fmt.Appendf(nil, "concurrent-%d", i)); err != nil {
					t.Errorf("Put() failed: %v", err)
				}
			}
		}()
	}
	wg.Wait()

	total := n * msgsPerG
	for range total {
		select {
		case <-q.ReadChan():
		case <-time.After(5 * time.Second):
			t.Fatalf("ReadChan() timed out, consumed %d/%d", total, total)
		}
	}

	if depth := q.Depth(); depth != 0 {
		t.Fatalf("Depth() = %d, want 0", depth)
	}
}

func TestDepth(t *testing.T) {
	q := newTestQueue(t, "TestDepth", 1024*1024)

	if depth := q.Depth(); depth != 0 {
		t.Fatalf("Depth() = %d, want 0 (empty queue)", depth)
	}

	putMessage(t, q, []byte("a"))
	putMessage(t, q, []byte("b"))
	putMessage(t, q, []byte("c"))

	if depth := q.Depth(); depth != 3 {
		t.Fatalf("Depth() = %d, want 3", depth)
	}

	<-q.ReadChan()
	// Depth 查询可能有短暂延迟，等一会
	time.Sleep(50 * time.Millisecond)
	depth := q.Depth()
	if depth != 2 {
		t.Fatalf("Depth() after consume = %d, want 2", depth)
	}
}

// ============================================================
// 性能测试
// ============================================================

func BenchmarkPut(b *testing.B) {
	for _, msgSize := range []int{128, 1024, 4096, 65536} {
		b.Run(fmt.Sprintf("msgSize=%d", msgSize), func(b *testing.B) {
			q := newTestQueue(b, "BenchPut", 64*1024*1024)
			msg := make([]byte, msgSize)

			b.ResetTimer()
			b.SetBytes(int64(msgSize))
			for range b.N {
				if err := q.Put(msg); err != nil {
					b.Fatalf("Put() failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkPutAndRead(b *testing.B) {
	for _, msgSize := range []int{128, 1024, 4096, 65536} {
		b.Run(fmt.Sprintf("msgSize=%d", msgSize), func(b *testing.B) {
			q := newTestQueue(b, "BenchPutRead", 64*1024*1024)
			msg := make([]byte, msgSize)

			done := make(chan struct{})
			go func() {
				for range b.N {
					<-q.ReadChan()
				}
				close(done)
			}()

			b.ResetTimer()
			b.SetBytes(int64(msgSize))
			for range b.N {
				if err := q.Put(msg); err != nil {
					b.Fatalf("Put() failed: %v", err)
				}
			}
			<-done
		})
	}
}

func BenchmarkPutConsumerBehind(b *testing.B) {
	// 模拟生产远快于消费的场景：先批量写入，再消费
	b.StopTimer()
	q := newTestQueue(b, "BenchPutBehind", 64*1024*1024)
	msg := make([]byte, 1024)

	b.StartTimer()
	b.SetBytes(1024)

	// 先写入 b.N 条
	for range b.N {
		if err := q.Put(msg); err != nil {
			b.Fatalf("Put() failed: %v", err)
		}
	}

	b.StopTimer()
	// 消费完
	for range b.N {
		<-q.ReadChan()
	}
}

func BenchmarkReadChan(b *testing.B) {
	q := newTestQueue(b, "BenchRead", 64*1024*1024)
	msg := make([]byte, 1024)

	// 预填充
	for range b.N {
		if err := q.Put(msg); err != nil {
			b.Fatalf("Put() failed: %v", err)
		}
	}

	b.ResetTimer()
	b.SetBytes(1024)
	for range b.N {
		<-q.ReadChan()
	}
}

// ============================================================
// 稳定性测试
// ============================================================

func TestLargeMessagePutAndRead(t *testing.T) {
	q := newTestQueue(t, "TestLargeMsg", 64*1024*1024)
	msg := make([]byte, 512*1024) // 512KB
	for i := range msg {
		msg[i] = byte(i % 256)
	}

	putMessage(t, q, msg)

	select {
	case data := <-q.ReadChan():
		if len(data) != len(msg) {
			t.Fatalf("ReadChan() len = %d, want %d", len(data), len(msg))
		}
		for i := range msg {
			if data[i] != msg[i] {
				t.Fatalf("data mismatch at byte %d: got %d, want %d", i, data[i], msg[i])
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ReadChan() timed out for large message")
	}
}

func TestWriteAfterClose(t *testing.T) {
	q := newTestQueue(t, "TestWriteAfterClose", 1024*1024)
	if err := q.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	if err := q.Put([]byte("after close")); err == nil {
		t.Fatal("Put() after Close should have failed")
	}
}
