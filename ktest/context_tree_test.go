package ktest

import (
	"sync"
	"testing"
	"time"

	"github.com/khan-lau/kutils/container/kcontext"
	"github.com/khan-lau/kutils/container/kstrings"
)

func Test_ContextTree(t *testing.T) {
	// 创建上下文树
	tree := kcontext.NewContextTree("根节点")
	root := tree.GetRoot()

	// 设置根节点tag
	root.SetTag("根标签")

	// 创建树结构，包含重名节点
	child1 := root.NewChild("子节点1")
	child2 := root.NewChildWithTag("子节点2", make(chan string, 1))
	grandchild1 := child1.NewChildWithTag("子节点1->任务1", make(chan int, 1)) // 重名任务1
	grandchild2 := child2.NewChild("子节点2->子节点2")                          // 重名任务1
	grandchild3 := child2.NewChild("子节点2->子节点3")                          // 重名子节点1

	// 添加监听器
	listenerID1 := child1.AddListener(func(source *kcontext.ContextNode, msg interface{}) {
		kstrings.Debugf("---  {} 收到消息: `{}`, 来源: {}\n", child1.Name(), msg, source.Name())
	})
	listenerID2 := grandchild1.AddListener(func(source *kcontext.ContextNode, msg interface{}) {
		if ch, ok := source.Tag().(chan int); ok {
			ch <- msg.(int)
		}
	})

	// 并发添加子节点和监听器
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node := child2.NewChild(kstrings.Sprintf("子任务%d", i))
			node.AddListener(func(source *kcontext.ContextNode, msg interface{}) {
				kstrings.Debugf("{} 收到消息: {}, 来源: {}\n", node.Name(), msg, source.Name())
			})
		}(i)
	}
	wg.Wait()

	// 使用chan类型的tag
	if ch, ok := child2.Tag().(chan string); ok {
		go func() { ch <- "子节点2消息" }()
		select {
		case msg := <-ch:
			kstrings.Debugf("从子节点2的tag通道接收: %s\n", msg)
		case <-time.After(time.Second):
			kstrings.Debugf("子节点2的tag通道超时\n")
		}
	}

	// 通知消息
	child1.Notify(grandchild1, "测试消息")
	grandchild1.Notify(child1, 42)

	// 检查grandchild1的tag通道
	if ch, ok := grandchild1.Tag().(chan int); ok {
		select {
		case val := <-ch:
			kstrings.Debugf("从{}的tag通道接收: {}\n", grandchild1.Name(), val)
		case <-time.After(time.Millisecond * 100):
			kstrings.Debugf("{}的tag通道超时\n", grandchild1.Name())
		}
	}

	// 输出JSON
	kstrings.Debugf("\n树JSON表示: {}", root.String())

	// 查询重名节点
	kstrings.Debug("\n查找所有名称为'任务1'的节点：")
	results := root.FindAllNodesByName("任务1")
	for i, result := range results {
		node := result.Node
		parent := node.Parent()
		tag := node.Tag()
		kstrings.Debugf("匹配 %d: %s (ID: %s), 路径: %v, 父节点: %s, 标签: %v\n",
			i+1, node.Name(), node.ID(), result.Path, parent.Name(), tag)
	}

	// 移除监听器
	child1.RemoveListener(listenerID1)
	grandchild1.RemoveListener(listenerID2)

	// 取消和移除
	kstrings.Debug("\n取消 {}", child1.Name())
	child1.Cancel()

	kstrings.Println("移除 {}...", child2.Name())
	child2.Remove()

	grandchild2.Cancel()
	kstrings.Debug("取消{}...", grandchild2.Name())
	grandchild2.Remove()

	grandchild3.Cancel()
	kstrings.Debug("取消{}...", grandchild3.Name())
	grandchild3.Remove()

	kstrings.Debug("关闭整个树...")
	tree.Close()

	grandchild2.Cancel()
	kstrings.Debug("取消{}...", grandchild2.Name())
	grandchild2.Remove()

	// 输出JSON
	kstrings.Debugf("\n树JSON表示: {}", root.String())
}
