package kcontext

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/khan-lau/kutils/container/kmaps"
	"github.com/khan-lau/kutils/container/kstrings"
	"github.com/khan-lau/kutils/kuuid"
)

type MessageCallback func(source *ContextNode, msg interface{})

/////////////////////////////////////////////////////////////////////////

// ContextManager 管理整个上下文树，提供统一入口
type ContextTree struct {
	root *ContextNode
}

// NewContextManager 创建新的上下文管理器
// 相比原生context，提供显式树结构、资源清理和名称/ID查询功能
func NewContextTree(rootName string) *ContextTree {
	rootCtx, rootCancel := context.WithCancel(context.Background())
	uuid, _ := kuuid.NewV1()
	return &ContextTree{
		root: &ContextNode{
			name:      rootName,
			id:        uuid.ShortString(),
			tag:       nil, // tag初始化为nil
			ctx:       rootCtx,
			cancel:    rootCancel,
			parent:    nil, // 根节点没有父节点, 为nil
			children:  make([]*ContextNode, 0),
			listeners: make(map[string]MessageCallback, 0),
		},
	}
}

// Close 关闭整个上下文树，清理所有监听器和通道
func (that *ContextTree) Close() {
	that.root.closeWithChan()
}

// GetRoot 返回根节点
func (that *ContextTree) GetRoot() *ContextNode {
	return that.root
}

/////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////

// Node 代表树状结构中的一个节点，封装context、cancel函数、名称和唯一ID
// 提供显式的父子关系管理，防止内存泄露，支持并发安全操作和重名场景查询
type ContextNode struct {
	name       string                     // 节点名称，可重复
	id         string                     // 唯一标识符，区分重名节点
	tag        interface{}                // 节点标签，存储任意类型元数据
	ctx        context.Context            // 节点的context
	cancel     context.CancelFunc         // 节点的取消函数
	parent     *ContextNode               // 父节点引用
	children   []*ContextNode             // 子节点列表
	nodeMu     sync.RWMutex               // 保护节点操作的读写锁
	listeners  map[string]MessageCallback // 监听器
	listenerMu sync.RWMutex               // 保护监听器列表的读写锁
}

// AddChild 为父节点添加子节点
// 自动创建子context并维护父子关系，生成唯一ID，线程安全
func (that *ContextNode) NewChild(name string) *ContextNode {
	return that.NewChildWithTag(name, nil)
}

// AddChild 为父节点添加子节点
// 自动创建子context并维护父子关系，生成唯一ID，线程安全
func (that *ContextNode) NewChildWithTag(name string, tag interface{}) *ContextNode {
	uuid, _ := kuuid.NewV1()
	that.nodeMu.Lock()
	defer that.nodeMu.Unlock()

	// 创建子context，绑定到父context
	childCtx, childCancel := context.WithCancel(that.ctx)
	child := &ContextNode{
		name:      name,
		id:        uuid.ShortString(),
		tag:       tag,
		ctx:       childCtx,
		cancel:    childCancel,
		parent:    that,
		children:  make([]*ContextNode, 0),
		listeners: make(map[string]MessageCallback, 0),
	}

	that.children = append(that.children, child)
	return child
}

// Cancel 取消当前节点及其子节点的context
// 与原生context的cancel等价，子context自动收到Done信号
func (that *ContextNode) Cancel() {
	that.nodeMu.Lock()
	defer that.nodeMu.Unlock()

	that.cancel()
}

// AddListener 添加监听器，返回唯一监听器ID
func (that *ContextNode) AddListener(listener MessageCallback) (listenerID string) {
	if listener == nil {
		return ""
	}
	uuid, _ := kuuid.NewV1()
	listenerID = uuid.ShortString()

	that.listenerMu.Lock()
	defer that.listenerMu.Unlock()
	that.listeners[listenerID] = listener
	return listenerID
}

// RemoveListener 移除指定监听器
func (that *ContextNode) RemoveListener(listenerID string) {
	that.listenerMu.Lock()
	defer that.listenerMu.Unlock()
	delete(that.listeners, listenerID)
}

// Notify 异步通知所有监听器
func (that *ContextNode) Notify(source *ContextNode, msg interface{}) {
	that.listenerMu.RLock()
	listeners := make(map[string]MessageCallback, len(that.listeners))
	for id, listener := range that.listeners {
		listeners[id] = listener
	}
	that.listenerMu.RUnlock()

	// 异步调用监听器
	var wg sync.WaitGroup
	for _, listener := range listeners {
		if listener != nil {
			wg.Add(1)
			go func(callback MessageCallback) {
				defer wg.Done()
				callback(source, msg)
			}(listener)
		}
	}
	wg.Wait()
}

// Context 返回节点的context
func (that *ContextNode) Context() context.Context {
	return that.ctx
}

// Name 返回节点的名称
func (that *ContextNode) Name() string {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()
	return that.name
}

// ID 返回节点的唯一ID
func (that *ContextNode) ID() string {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()
	return that.id
}

// SetTag 设置节点的tag值
// 支持任意类型，线程安全
func (that *ContextNode) SetTag(tag interface{}) {
	that.nodeMu.Lock()
	defer that.nodeMu.Unlock()
	that.tag = tag
}

// GetTag 获取节点的tag值
// 返回nil如果未设置，线程安全
func (that *ContextNode) Tag() interface{} {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()
	return that.tag
}

// GetParent 返回当前节点的父节点
// 如果是根节点，返回nil，线程安全
func (that *ContextNode) Parent() *ContextNode {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()
	return that.parent
}

// GetRoot 返回树的根节点
// 通过递归向上遍历，线程安全
func (that *ContextNode) Root() *ContextNode {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	if that.parent == nil {
		return that
	}
	return that.parent.Root()
}

// Remove 从父节点中移除当前节点，清理引用以防止内存泄露
func (that *ContextNode) Remove() {
	that.nodeMu.Lock()
	defer that.nodeMu.Unlock()

	that.cancel()

	// 从父节点的子节点列表中移除当前节点
	if that.parent != nil {
		that.parent.nodeMu.Lock()
		for i, child := range that.parent.children {
			if child == that {
				that.parent.children = append(that.parent.children[:i], that.parent.children[i+1:]...)
				break
			}
		}
		that.parent.nodeMu.Unlock()
		that.parent = nil
	}
}

// GetChildren 返回当前节点的子节点列表（只读）
func (that *ContextNode) Children() []*ContextNode {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()
	return append([]*ContextNode{}, that.children...)
}

// GetPath 返回节点从根到自身的路径（名称列表）
func (that *ContextNode) GetPath() []string {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	path := []string{that.name}
	for curr := that.parent; curr != nil; curr = curr.parent {
		curr.nodeMu.RLock()
		path = append([]string{curr.name}, path...)
		curr.nodeMu.RUnlock()
	}
	return path
}

// GetPathWithIDs 返回节点从根到自身的路径（ID列表）
func (that *ContextNode) GetPathWithIDs() []string {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	path := []string{that.id}
	for curr := that.parent; curr != nil; curr = curr.parent {
		curr.nodeMu.RLock()
		path = append([]string{curr.id}, path...)
		curr.nodeMu.RUnlock()
	}
	return path
}

// FindNodeByName 使用DFS查找第一个匹配名称的节点
// 可选parentID约束查找范围，返回节点及其路径，若未找到返回nil
func (that *ContextNode) FindNodeByName(name string, parentID string) (*ContextNode, []string) {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	// 如果指定parentID，检查当前节点是否在正确分支
	if parentID != "" && that.id != parentID {
		isDescendant := false
		for curr := that.parent; curr != nil; curr = curr.parent {
			curr.nodeMu.RLock()
			if curr.id == parentID {
				isDescendant = true
			}
			curr.nodeMu.RUnlock()
			if isDescendant {
				break
			}
		}
		if !isDescendant && that.parent != nil {
			return nil, nil
		}
	}

	// 检查当前节点
	if that.name == name {
		return that, that.GetPath()
	}

	// 递归搜索子节点
	children := append([]*ContextNode{}, that.children...)
	for _, child := range children {
		if found, path := child.FindNodeByName(name, parentID); found != nil {
			return found, path
		}
	}
	return nil, nil
}

// FindAllNodesByName 使用DFS查找所有匹配名称的节点
// 返回节点及其路径的列表，若无匹配返回空列表
func (that *ContextNode) FindAllNodesByName(name string) []struct {
	Node *ContextNode
	Path []string
} {
	var results []struct {
		Node *ContextNode
		Path []string
	}

	that.nodeMu.RLock()
	if that.name == name {
		results = append(results, struct {
			Node *ContextNode
			Path []string
		}{Node: that, Path: that.GetPath()})
	}
	that.nodeMu.RUnlock()

	// 递归搜索子节点
	that.nodeMu.RLock()
	children := append([]*ContextNode{}, that.children...)
	that.nodeMu.RUnlock()

	for _, child := range children {
		results = append(results, child.FindAllNodesByName(name)...)
	}
	return results
}

// FindNodeByID 使用DFS查找匹配ID的节点
// ID唯一，直接返回节点及其路径，若未找到返回nil
func (that *ContextNode) FindNodeByID(id string) (*ContextNode, []string) {
	that.nodeMu.RLock()
	if that.id == id {
		that.nodeMu.RUnlock()
		return that, that.GetPath()
	}
	that.nodeMu.RUnlock()

	that.nodeMu.RLock()
	children := append([]*ContextNode{}, that.children...)
	that.nodeMu.RUnlock()

	for _, child := range children {
		if found, path := child.FindNodeByID(id); found != nil {
			return found, path
		}
	}
	return nil, nil
}

// String 返回以当前节点为根的子树的JSON表示
// 仅包含id、name和children字段，格式化输出，线程安全
func (that *ContextNode) String() string {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	// 构建JSON节点
	jn := &jsonNode{
		ID:       that.id,
		Name:     that.name,
		Children: make([]*jsonNode, 0, len(that.children)),
	}

	// 递归添加子节点
	for _, child := range that.children {
		childJSON := child.toJSONNode()
		jn.Children = append(jn.Children, childJSON)
	}

	// 序列化为格式化JSON
	data, err := json.MarshalIndent(jn, "", "  ")
	if err != nil {
		return kstrings.Sprintf("{\"error\": \"failed to marshal JSON: {}\"}", err)
	}
	return string(data)
}

/////////////////////////////////////////////////////////////////////////

// closeWithChan 递归关闭节点及其子节点的tag中的通道和监听器
func (that *ContextNode) closeWithChan() {
	that.nodeMu.Lock()
	defer that.nodeMu.Unlock()

	that.cancel()
	that.tryCloseChan()

	// 清空监听器
	that.listenerMu.Lock()
	kmaps.Clear(that.listeners)
	that.listenerMu.Unlock()

	for _, child := range that.children {
		child.closeWithChan()
	}
}

// tryCloseChan 尝试关闭chan类型的tag
func (that *ContextNode) tryCloseChan() {
	if that.tag == nil {
		return
	}

	v := reflect.ValueOf(that.tag)
	if v.Kind() == reflect.Chan && !v.IsNil() {
		v.Close()
	}
}

// jsonNode 辅助结构体，用于JSON序列化，仅包含id、name和children
type jsonNode struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Children []*jsonNode `json:"children"`
}

// toJSONNode 辅助方法，将Node转换为jsonNode
func (that *ContextNode) toJSONNode() *jsonNode {
	that.nodeMu.RLock()
	defer that.nodeMu.RUnlock()

	jn := &jsonNode{
		ID:       that.id,
		Name:     that.name,
		Children: make([]*jsonNode, 0, len(that.children)),
	}

	for _, child := range that.children {
		jn.Children = append(jn.Children, child.toJSONNode())
	}
	return jn
}
