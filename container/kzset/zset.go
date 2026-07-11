package kzset

import (
	"math/rand"
	"sync"
	"time"
)

const (
	maxLevel = 16   // 定义了跳表能够允许达到的最高指针层数。16层索引在概率上足以支撑 50万 至 100万 级别的白名单数据保持 O(log N) 的极速查找。32层索引可以支撑 3000万条数据, redis实现是 32 层
	p        = 0.25 // 代表跳表节点在创建时，向更高一层晋升的概率。 0.25 (1/4) 是经典的空间与时间折中值（Redis 亦采用此值），平均每个节点占用 1.33 个指针。
)

// skipListNode 代表跳表内部的单链表节点实体。
// 该结构体为内部私有结构，不对外部业务线暴露。
type skipListNode struct {
	member string          // 唯一标识元素（如：白名单测点ID）
	score  int64           // 排序权重分数（如：绝对过期Unix时间戳）
	next   []*skipListNode // 向前指针数组。next[i] 表示当前节点在第 i 层索引指向的下一个节点实体
}

// skipList 代表无锁保护的原生跳跃表结构。
// 该结构体不具备并发安全性，所有高频写操作必须由上层的 GoZSet 锁逻辑进行编排。
type skipList struct {
	header *skipListNode // 跳表的头哨兵节点，专用于各层指针检索的起点
	level  int           // 当前跳表中已经建立起索引的实际最高层数
}

// GoZSet 是一个高性能、线程安全的有序集合（ZSet）实现。
//
// 使用场景：
// 适用于高频、大吞吐量的物联网/微服务动态白名单控制、动态生命周期监控等场景。
//
// 设计原理：
// 内部通过“哈希表（Map）+ 跳跃表（SkipList）”双轨驱动。
// Map 提供 O(1) 的精准反查与去重，SkipList 提供 O(log N) 的时间线自动排序与高效切片。
//
// 并发特性：
// 内部封装了 sync.RWMutex 读写分离锁。
// 支持全量高频多协程并发无锁读取（RLock），写/删/更新操作彼此互斥（Lock）。
//
// example:
//
//	zset := NewGoZSet()
//	zset.ZAdd("point_A", 1000)
//	zset.ZAdd("point_B", 2000)
//	fmt.Println("🟢 方案 3 运行成功！已过期数据:", zset.ZRangeByScore(1500))
type GoZSet struct {
	mu         sync.RWMutex     // 读写分离保护锁，协调底层 Map 与 SkipList 的级联并发安全
	dict       map[string]int64 // 哈希字典：Member -> Score，用于 O(1) 的精准反查与去重
	sl         *skipList        // 跳跃表：用于维持所有成员严格按 Score 从小到大有序排列
	randSource *rand.Rand       // 对象私有的独立随机数发生器，彻底绕过标准库全局锁

	// 🚀 方案 3 杀手锏：长在结构体身上的写操作公共指针缓存区。
	// 物理上属于连续内存（数组），利用写锁互斥的特性，实现 0-Alloc 级复用。
	updateBuf [maxLevel]*skipListNode
}

// NewGoZSet 用于初始化并返回一个完整的 GoZSet 容器实例。
//
// 返回值：
//
//	@returns *GoZSet: 初始化完毕的有序集合指针，已内嵌预分配好的内存空间与私有随机源。
//
// 使用方法：
//
//	zset := NewGoZSet()
func NewGoZSet() *GoZSet {
	return &GoZSet{
		dict:       make(map[string]int64, 500000), // 预分配 50w 容量，防止线上扩容挂起
		sl:         &skipList{header: &skipListNode{next: make([]*skipListNode, maxLevel)}, level: 1},
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// randomLevel 是一个内部辅助方法，用于通过随机概率决定一个新插入的节点应该拥有多高的索引层数。
// 遵循幂律分布，层数越高，被选中的概率呈几何级数（1/4）递减。
//
// 返回值：
//
//	@returns int: 随机生成的层数，范围严格限制在 [1, maxLevel] 之间。
func (z *GoZSet) randomLevel() int {
	lvl := 1
	for z.randSource.Float64() < p && lvl < maxLevel {
		lvl++
	}
	return lvl
}

// ZAdd 往有序集合中添加一个新成员，或者更新一个已存在成员的分数。
//
// 参数描述：
//
//	@param member (string): 成员的唯一标识标识（例如：测点编码 "NM0001"）。保证全局唯一，不可重复。
//	@param score (int64)  : 该成员绑定的分数权重（例如：未来的绝对过期秒级时间戳 1719912000）。
//
// 行为逻辑：
//  1. 幂等/更新检查：若 member 已存在且分数相同，直接忽略并返回；若分数不同，则先从底层跳表无脑剥离旧节点，再执行新数据插入。
//  2. 排序构建：新节点在底层通过多层级联指针搜索，耗时 O(log N) 自动降落到有序铁轨的正确坑位中。
//
// 注意事项：
//
//	此函数内部会上【写锁】，多个协程同时触发 ZAdd 时会发生微秒级排队，属于正常并发保护。
//
// 使用方法：
//
//	zset.ZAdd("point_A", time.Now().Add(10 * time.Minute).Unix())
func (z *GoZSet) ZAdd(member string, score int64) {
	z.mu.Lock()
	defer z.mu.Unlock()

	// 1. 拦截更新/去重逻辑
	if oldScore, exists := z.dict[member]; exists {
		if oldScore == score {
			return
		}
		// 分数变更，必须先从跳表中切除旧节点（复用无锁的内部单例删除方法）
		z.sl.internalDelete(member, oldScore, &z.updateBuf)
	}

	// 2. 映射哈希字典
	z.dict[member] = score

	// 3. 级联插入跳表结构
	curr := z.sl.header

	// 从当前最高层向下逼近，寻找最贴近新插入分数的前驱节点集合
	for i := z.sl.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && (curr.next[i].score < score || (curr.next[i].score == score && curr.next[i].member < member)) {
			curr = curr.next[i]
		}
		// 🟢 方案 3 的威力显现：直接写入公共缓存，无任何 make 开销
		z.updateBuf[i] = curr
	}

	// 计算新节点拥有的层数高度
	lvl := z.randomLevel()
	if lvl > z.sl.level {
		for i := z.sl.level; i < lvl; i++ {
			z.updateBuf[i] = z.sl.header
		}
		z.sl.level = lvl
	}

	// 建立新节点，缝合高架桥指针
	newNode := &skipListNode{member: member, score: score, next: make([]*skipListNode, lvl)}
	for i := range lvl {
		newNode.next[i] = z.updateBuf[i].next[i]
		z.updateBuf[i].next[i] = newNode
	}
}

// ZRem 从有序集合中精准、级联删除指定的成员。
//
// 参数描述：
//
//	@param member (string): 准备剔除的成员唯一标识名字。
//
// 返回值：
//
//	@returns bool: 删除结果状态。返回 true 代表该元素存在且已成功抹去；返回 false 代表元素本身不存在，什么都没做。
//
// 行为逻辑：
//
//	通过 Map 获取该成员的分数，随后利用跳表的多层高架索引高速搜寻该节点，将其前后指针断开。
//	内存释放交给 Go 原生 GC，耗时严格控制在 O(log N)。
//
// 使用方法：
//
//	wasDeleted := zset.ZRem("point_A")
func (z *GoZSet) ZRem(member string) bool {
	z.mu.Lock()
	defer z.mu.Unlock()

	// 通过 O(1) 的 Map 确定分数，避免盲目搜索
	score, exists := z.dict[member]
	if !exists {
		return false
	}

	// 双轨联动：同步清除哈希表与跳跃表
	delete(z.dict, member)
	return z.sl.internalDelete(member, score, &z.updateBuf)
}

// ZRangeByScore 获取所有分数小于或等于指定最大分数的成员名单（用于提取已到期数据）。
//
// 参数描述：
//
//	@param maxScore (int64): 查询分数的临界上限（在白名单场景下，通常传入当前系统的时间戳 UnixTime）。
//
// 返回值：
//
//	@returns []string: 一个包含所有符合条件的成员名字的切片。若无过期数据，则返回 nil。
//
// 核心超频设计：
//
//	由于跳表底层（第 0 层）是一条严格从小到大单向排列的单链表。此函数会【直接上读锁】，
//	从最左侧（分数最小、最早过期的数据）向右线性顺序收割，一旦遇到第一个大于 maxScore 的元素，
//	说明后续几十万数据全部合法，立马执行 break 熔断。
//	当 50w 数据中只有 5 条过期时，只会循环 5 次。耗时与总数 N 彻底解耦！支持千万级高并发。
//
// 使用方法：
//
//	expiredPoints := zset.ZRangeByScore(time.Now().Unix())
func (z *GoZSet) ZRangeByScore(maxScore int64) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()

	var expired []string
	curr := z.sl.header.next[0] // 降落到跳表最底层的有序单链表起点

	for curr != nil {
		if curr.score <= maxScore {
			expired = append(expired, curr.member)
			curr = curr.next[0]
		} else {
			// ⚡ 核心熔断点：由于严格有序，后续元素必然无需扫描，完美避开无效大循环
			break
		}
	}
	return expired
}

// internalDelete 是跳表底层的内部辅助私有方法（不带锁保护，不可外部直接调用）。
// 它的职责是配合 ZAdd/ZRem，在写锁安全的上下文里执行跳表多级指针的断开与层数重修。
//
// 参数描述：
//
//	@param member (string): 待删除成员名。
//	@param score (int64)  : 待删除成员的分数。
//
// 返回值：
//
//	@returns bool: 是否在跳表中找到了该数据并成功清除。
//
// 参数优化：
//
//	通过增传 updateBuf 指针，让删除逻辑同样复用结构体自带的连续内存，完全干掉了底层的 make。
func (sl *skipList) internalDelete(member string, score int64, updateBuf *[maxLevel]*skipListNode) bool {
	curr := sl.header

	// 自上而下搜寻需要修改指针的前驱节点群，直接填入公共 Buffer
	for i := sl.level - 1; i >= 0; i-- {
		for curr.next[i] != nil && (curr.next[i].score < score || (curr.next[i].score == score && curr.next[i].member < member)) {
			curr = curr.next[i]
		}
		updateBuf[i] = curr
	}

	curr = curr.next[0]

	// 精准比对名字与分数，执行缝合手术断开链表指针
	if curr != nil && curr.score == score && curr.member == member {
		for i := 0; i < sl.level; i++ {
			if updateBuf[i].next[i] != curr {
				break
			}
			updateBuf[i].next[i] = curr.next[i]
		}
		// 对跳表实际最高层数执行降级
		for sl.level > 1 && sl.header.next[sl.level-1] == nil {
			sl.level--
		}
		return true
	}
	return false
}
