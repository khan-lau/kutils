package kredis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/khan-lau/kutils/container/kcontext"
	kslices "github.com/khan-lau/kutils/container/kslices"
	"github.com/khan-lau/kutils/klogger"

	redisHd "github.com/redis/go-redis/v9"
)

type KRedisCluster struct {
	Client *redisHd.ClusterClient
	ctx    *kcontext.ContextNode
}

func redisClusterOnConnect(ctx context.Context, cn *redisHd.Conn) error {
	return nil
}

func NewKRedisCluster(ctx *kcontext.ContextNode, addrs []string, user string, password string, dbNum int) *KRedisCluster {
	client := redisHd.NewClusterClient(&redisHd.ClusterOptions{
		Addrs:     addrs,
		Username:  user,
		Password:  password,
		OnConnect: redisClusterOnConnect, // 指定连接成功后的钩子函数
		// Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 	return net.DialTimeout(network, addr, 10*time.Second)
		// },
		Protocol:              2,                       // 协议版本, 2 代表 RESP2，3 代表 RESP3, RESP3 是 Redis 6.0 之后引入的高性能新型协议
		NewClient:             nil,                     // 当集群客户端发现需要连接某个具体的单机节点时，会调用这个函数
		ClusterSlots:          nil,                     // 提供分片信息，如果没有提供，会执行一次INFO SHARDS获取分片信息
		ReadOnly:              false,                   // 从节点是否只读, 除非读压力巨大且能容忍短暂脏读，否则不建议开只读
		RouteByLatency:        false,                   // 是否根据网络延迟自动路由读请求
		MaxRedirects:          8,                       // 遇到重定向（MOVED 或 ASK 错误）时的最大重试次数
		MaxRetries:            3,                       // 最大重试次数
		MinRetryBackoff:       8 * time.Millisecond,    // 重试间隔时间下限
		MaxRetryBackoff:       512 * time.Millisecond,  // 重试间隔时间上限
		DialTimeout:           10 * time.Second,        // 连接超时时间
		ReadTimeout:           4 * time.Second,         // socket 读取超时时间
		WriteTimeout:          4 * time.Second,         // socket 写入超时时间
		PoolFIFO:              true,                    // 空闲连接池队列是否采用先进先出方式
		PoolSize:              10,                      // 连接池大小
		MaxActiveConns:        20,                      // 每个 Redis 节点在同一时刻能够分配的最大激活连接数（包含连接池内的空闲连接和正在使用的连接）
		PoolTimeout:           4 * time.Second,         // 连接池等待超时时间，如果获取连接超过这个时间就会失败
		MinIdleConns:          2,                       // 连接池最小空闲连接数
		MaxIdleConns:          4,                       // 连接池最大空闲连接数
		ConnMaxIdleTime:       5 * time.Minute,         // 【优化】长连接空闲 5 分钟自动回收，防止占着连接
		ConnMaxLifetime:       30 * time.Minute,        // 【优化】强制每 30 分钟轮换长连接，防止隐性底层网络老化,
		ContextTimeoutEnabled: true,                    // 【优化】推荐开启！允许外部传入的 ctx.WithTimeout 强行中断请求
		ClientName:            "kredis_cluster_client", // 客户端标识
		DisableIndentity:      false,                   // 连接建立时是否设置客户端标识，默认是设置
		TLSConfig:             nil,                     // TLS 配置
	})

	subCtx := ctx.NewChild("kredis_cluster_client")
	return &KRedisCluster{Client: client, ctx: subCtx}
}

// 执行指令
func (that *KRedisCluster) Do(args ...any) (any, error) {
	val, err := that.Client.Do(that.ctx.Context(), args...).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 获取一个key的值
func (that *KRedisCluster) Get(key string) (any, error) {
	val, err := that.Client.Do(that.ctx.Context(), "GET", key).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 设置某个key的值, 并指定ttl
func (that *KRedisCluster) Set(key string, value any, duration time.Duration) (bool, error) {
	err := that.Client.Set(that.ctx.Context(), key, value, duration).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 判断某个key是否存在
func (that *KRedisCluster) Exist(key string) (bool, error) {
	_, err := that.Client.Get(that.ctx.Context(), key).Result()
	if err == redisHd.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// 获取一个key的hash字段的值
func (that *KRedisCluster) HGet(key string, field string) (any, error) {
	val, err := that.Client.HGet(that.ctx.Context(), key, field).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 设置一个key的hash字段的值
func (that *KRedisCluster) HSet(key string, field string, value any) error {
	err := that.Client.HSet(that.ctx.Context(), key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取一个key的hash字段的值列表
func (that *KRedisCluster) HGetAll(key string) (map[string]string, error) {
	valMap, err := that.Client.HGetAll(that.ctx.Context(), key).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return valMap, nil
}

// 设置一个key的hash字段的值列表
func (that *KRedisCluster) HSetAll(key string, fields map[string]any) error {
	err := that.Client.HSet(that.ctx.Context(), key, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

// 判断一个key的hash字段是否存在
func (that *KRedisCluster) HExists(key string, field string) (bool, error) {
	isExists, err := that.Client.HExists(that.ctx.Context(), key, field).Result()
	if nil != err {
		return false, err
	}
	return isExists, nil
}

// 获取一个key的hash字段的数量
func (that *KRedisCluster) HLen(key string) (int64, error) {
	val, err := that.Client.HLen(that.ctx.Context(), key).Result()
	if nil != err {
		return 0, err
	}
	return val, nil
}

// 获取一个key的hash字段的所有Key
func (that *KRedisCluster) HKeys(ctx context.Context, key string) ([]string, error) {
	array, err := that.Client.HKeys(ctx, key).Result()
	if nil != err {
		return nil, err
	}
	return array, nil
}

// 获取一个key的hash字段的所有值
func (that *KRedisCluster) HVals(ctx context.Context, key string) ([]string, error) {
	array, err := that.Client.HVals(ctx, key).Result()
	if nil != err {
		return nil, err
	}
	return array, nil
}

// 设置一个key的hash字段的值列表, 如果不存在则创建
func (that *KRedisCluster) HSetNX(ctx context.Context, key, field string, value any) (bool, error) {
	isExists, err := that.Client.HSetNX(ctx, key, field, value).Result()
	if nil != err {
		return false, err
	}
	return isExists, nil
}

// 删除一个key的hash字段的值列表
func (that *KRedisCluster) HDel(key string, fields ...string) error {
	_, err := that.Client.HDel(that.ctx.Context(), key, fields...).Result()
	if nil != err {
		return err
	}
	return nil
}

// 获取一个key的hash字段的值列表
func (that *KRedisCluster) HMGet(key string, fields ...string) ([]any, error) {
	valMap, err := that.Client.HMGet(that.ctx.Context(), key, fields...).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if nil != err {
		return nil, err
	}
	return valMap, nil
}

// 设置一个key的hash字段的值列表, 如果不存在则创建
func (that *KRedisCluster) HMSet(key string, fields map[string]any) error {
	err := that.Client.HMSet(that.ctx.Context(), key, fields).Err()
	if nil != err {
		return err
	}
	return nil
}

// 从列表左边插入数据
func (that *KRedisCluster) LPush(key string, values ...any) (int64, error) {
	return that.Client.LPush(that.ctx.Context(), key, values...).Result()
}

// 从列表左边插入数据, 如果不存在则不插入数据
func (that *KRedisCluster) LPushX(key string, values ...any) (int64, error) {
	return that.Client.LPushX(that.ctx.Context(), key, values...).Result()
}

// 从列表右边插入数据
func (that *KRedisCluster) RPush(key string, values ...any) (int64, error) {
	return that.Client.RPush(that.ctx.Context(), key, values...).Result()
}

// 从列表右边插入数据, 如果不存在则不插入数据
func (that *KRedisCluster) RPushX(key string, values ...any) (int64, error) {
	return that.Client.RPushX(that.ctx.Context(), key, values...).Result()
}

// 从列表左边弹出数据
func (that *KRedisCluster) LPop(key string) (string, error) {
	return that.Client.LPop(that.ctx.Context(), key).Result()
}

// 从列表右边弹出数据
func (that *KRedisCluster) RPop(key string) (string, error) {
	return that.Client.RPop(that.ctx.Context(), key).Result()
}

// 返回列表的一个范围内的数据，也可以返回全部数据
func (that *KRedisCluster) LRange(key string, start int64, stop int64) ([]string, error) {
	return that.Client.LRange(that.ctx.Context(), key, start, stop).Result()
}

// 返回列表的大小
func (that *KRedisCluster) LLen(key string) (int64, error) {
	return that.Client.LLen(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) LTrim(key string, start int64, stop int64) error {
	return that.Client.LTrim(that.ctx.Context(), key, start, stop).Err()
}

func (that *KRedisCluster) LSet(key string, index int64, value any) error {
	return that.Client.LSet(that.ctx.Context(), key, index, value).Err()
}

// 删除列表中的数据
func (that *KRedisCluster) LRem(key string, count int64, value any) (int64, error) {
	return that.Client.LRem(that.ctx.Context(), key, count, value).Result()
}

// 根据索引坐标，查询列表中的数据
func (that *KRedisCluster) LIndex(key string, index int64) (string, error) {
	return that.Client.LIndex(that.ctx.Context(), key, index).Result()
}

// 在指定位置插入数据，在头部插入用"before"，尾部插入用"after"
func (that *KRedisCluster) LInsert(key string, position string, pivot any, value any) (int64, error) {
	return that.Client.LInsert(that.ctx.Context(), key, position, pivot, value).Result()
}

func (that *KRedisCluster) SAdd(key string, members ...any) (int64, error) {
	return that.Client.SAdd(that.ctx.Context(), key, members...).Result()
}

func (that *KRedisCluster) SMembers(key string) ([]string, error) {
	return that.Client.SMembers(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) SRem(key string, members ...any) (int64, error) {
	return that.Client.SRem(that.ctx.Context(), key, members...).Result()
}

func (that *KRedisCluster) SIsMember(key string, member any) (bool, error) {
	return that.Client.SIsMember(that.ctx.Context(), key, member).Result()
}

func (that *KRedisCluster) SCard(key string) (int64, error) {
	return that.Client.SCard(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) SPop(key string) (string, error) {
	return that.Client.SPop(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) SPopN(key string, count int64) ([]string, error) {
	return that.Client.SPopN(that.ctx.Context(), key, count).Result()
}

func (that *KRedisCluster) SUnion(keys ...string) ([]string, error) {
	return that.Client.SUnion(that.ctx.Context(), keys...).Result()
}

func (that *KRedisCluster) SUnionStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SUnionStore(that.ctx.Context(), destKey, keys...).Result()
}

func (that *KRedisCluster) SInter(keys ...string) ([]string, error) {
	return that.Client.SInter(that.ctx.Context(), keys...).Result()
}

func (that *KRedisCluster) SInterStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SInterStore(that.ctx.Context(), destKey, keys...).Result()
}

func (that *KRedisCluster) SDiff(keys ...string) ([]string, error) {
	return that.Client.SDiff(that.ctx.Context(), keys...).Result()
}

func (that *KRedisCluster) SDiffStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SDiffStore(that.ctx.Context(), destKey, keys...).Result()
}

func (that *KRedisCluster) SMove(source, destination string, member any) (bool, error) {
	return that.Client.SMove(that.ctx.Context(), source, destination, member).Result()
}

func (that *KRedisCluster) SRandMember(key string) (string, error) {
	return that.Client.SRandMember(that.ctx.Context(), key).Result()
}

// 获取一个key的数据类型, 数据类型全小写
func (that *KRedisCluster) Type(key string) (string, error) {
	dataType, err := that.Client.Type(that.ctx.Context(), key).Result()
	if nil != err {
		return "", err
	}
	return strings.ToLower(dataType), nil
}

// 返回一个Key的过期时间, 单位为毫秒
func (that *KRedisCluster) PTTL(key string) (time.Duration, error) {
	return that.Client.PTTL(that.ctx.Context(), key).Result()
}

// 返回一个Key的过期时间, 单位为秒
func (that *KRedisCluster) TTL(key string) (time.Duration, error) {
	return that.Client.TTL(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) Expire(key string, expiration time.Duration) bool {
	return that.Client.Expire(that.ctx.Context(), key, expiration).Val()
}

func (that *KRedisCluster) ExpireAt(key string, tm time.Time) bool {
	return that.Client.ExpireAt(that.ctx.Context(), key, tm).Val()
}

// JsonGet 封装了 Redis JSON.GET 命令，并直接返回原始的 JSON 字符串。
// 调用方将负责对返回的字符串进行反序列化。
// key: Redis 中的键名。
// paths: 可选的 JSON Path 参数。如果没有提供，则获取整个 JSON 文档。
// 返回值：JSON 字符串。如果键或路径不存在，或结果为空，则返回空字符串和 nil 错误。
func (that *KRedisCluster) JsonGet(key string, paths ...string) (string, error) {
	args := make([]any, 0, 2+len(paths))
	args = append(args, "JSON.GET", key)
	for _, path := range paths {
		args = append(args, path)
	}
	// 执行 Redis 命令
	cmd := that.Client.Do(that.ctx.Context(), args...)
	// 检查 Redis 命令执行是否出错
	if err := cmd.Err(); err != nil {
		// 如果错误是 redis.Nil (表示键或路径不存在)，则返回空字符串和 nil 错误
		if err == redisHd.Nil {
			return "", nil // 数据未找到，非致命错误，返回空字符串
		}
		// 其他错误则包装并返回
		return "", err
	}
	// 获取 Redis 返回的原始 JSON 字符串
	// JSON.GET 返回的是一个 JSON 格式的字符串
	jsonString, err := cmd.Text() // 使用 Text() 方法直接获取字符串
	if err != nil {
		// 如果从命令结果中获取字符串失败
		return "", err
	}

	return jsonString, nil // 返回 JSON 字符串和 nil 错误
}

// JsonSet 封装了 Redis JSON.SET 命令。
// 它会将 Go 值自动序列化为 JSON 字符串并存储。
// key: Redis 中的键名。
// path: JSON Path。
// value: 要存储的 Go 值，将被序列化为 JSON。
// 返回值：如果操作成功则返回 nil，否则返回错误。
func (that *KRedisCluster) JsonSet(key string, path string, value string) error {
	// 构建 Redis 命令参数：JSON.SET key path jsonValueString
	// `string(jsonValue)` 将字节切片转换为字符串，go-redis 可以接受
	cmd := that.Client.Do(that.ctx.Context(), "JSON.SET", key, path, value)

	// 检查 Redis 命令执行是否出错
	if err := cmd.Err(); err != nil {
		return err
	}
	// JSON.SET 通常返回 "OK" 或 nil，这里我们只关心错误
	return nil
}

// JsonMerge 封装了 Redis JSON.MERGE 命令。
// 该命令会将 Go 值序列化为 JSON 并与现有值合并。如果路径不存在，则创建新字段。
// key: Redis 中的键名。
// path: JSON Path，指定要合并的位置。
// value: 要存储的 Go 值，将被序列化为 JSON 并与现有值合并。
// 返回值：如果操作成功则返回 nil，否则返回错误。
func (that *KRedisCluster) JsonMerge(key string, path string, value string) error {
	cmd := that.Client.Do(that.ctx.Context(), "JSON.MERGE", key, path, value)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// JsonDel 封装了 Redis JSON.DEL 命令。
// key: Redis 中的键名。
// path: 可选的 JSON Path。如果为空字符串，则删除整个 JSON 文档。
// 返回值：被删除的 JSON 值数量。如果键或路径不存在，通常返回 0。
func (that *KRedisCluster) JsonDel(key string, path string) (int64, error) {
	// 构建 Redis 命令参数
	args := make([]any, 0, 3)
	args = append(args, "JSON.DEL", key)
	if path != "" { // 如果 path 不为空，则添加到参数中
		args = append(args, path)
	}

	cmd := that.Client.Do(that.ctx.Context(), args...)

	// 检查 Redis 命令执行是否出错
	if err := cmd.Err(); err != nil {
		// JSON.DEL 在键或路径不存在时通常返回 0 而不是 redis.Nil
		// 所以如果这里有错误，通常是更底层的问题
		return 0, err
	}

	// JSON.DEL 返回被删除的路径数量，这是一个整数
	result, err := cmd.Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// JsonType 封装了 Redis JSON.TYPE 命令。
// key: Redis 中的键名。
// path: 可选的 JSON Path。如果为空字符串，则返回根路径的类型; 不支持同时指定多个路径。
// 返回值：一个包含 JSON 值类型的字符串切片。如果键或路径不存在，则返回 nil 切片和 nil 错误。
func (that *KRedisCluster) JsonType(key string, path string) ([]string, error) {
	args := make([]any, 0, 3)
	args = append(args, "JSON.TYPE", key)
	if path != "" { // 如果 path 不为空，则添加到参数中
		args = append(args, path)
	}

	cmd := that.Client.Do(that.ctx.Context(), args...)

	if err := cmd.Err(); err != nil {
		if err == redisHd.Nil {
			return nil, nil // 键或路径不存在，返回 nil 切片和 nil 错误
		}
		return nil, err
	}

	// str, err := cmd.Text()
	// fmt.Printf("%s, %v\n", str, err)

	// JSON.TYPE 返回一个字符串数组（即使只有一个结果），例如 `["string"]`
	types, err := cmd.Slice()
	if err != nil {
		return nil, err
	}
	stringTypes := make([]string, 0, len(types))
	for _, v := range types {
		switch val := v.(type) { // 使用类型断言和 switch 语句来处理不同类型的值
		case string: // 如果 v 是字符串类型，直接赋值
			stringTypes = append(stringTypes, val) // 直接将字符串值追加到切片中
		case []byte: // 如果 v 是字节切片类型，转换为字符串
			stringTypes = append(stringTypes, string(val)) // 将字节切片转换为字符串后追加到切片中
		case []string:
			stringTypes = append(stringTypes, val...) // 直接将每个元素追加到结果切片中
		case []any:
			for _, iv := range val { // 遍历 []any 中的每个元素
				stringTypes = append(stringTypes, iv.(string)) // 断言为 string 并追加到结果切片中
			}
		case nil: // 如果 v 是 nil，则忽略它（通常不会发生）
			continue
		default:
			return nil, fmt.Errorf("unexpected type: %v", val) // 如果不是预期的类型（这里是 string 或 []byte），则返回错误
		}
		// stringTypes[i] = v.(string) // 确保切片中的元素类型为 string
	}
	return stringTypes, nil
}

// JsonObjKeys 封装了 Redis JSON.OBJKEYS 命令。
// key: Redis 中的键名。
// path: 可选的 JSON Path。如果为空字符串，则返回根对象的键。
// 返回值：一个包含对象键的字符串切片。如果键、路径不存在或路径对应的不是对象，则返回 nil 切片和 nil 错误。
func (that *KRedisCluster) JsonObjKeys(key string, path string) ([]string, error) {
	args := make([]any, 0, 3)
	args = append(args, "JSON.OBJKEYS", key)
	if path != "" { // 如果 path 不为空，则添加到参数中
		args = append(args, path)
	}

	cmd := that.Client.Do(that.ctx.Context(), args...)

	if err := cmd.Err(); err != nil {
		if err == redisHd.Nil {
			return nil, nil // 键或路径不存在，返回 nil 切片和 nil 错误
		}
		return nil, err
	}

	// JSON.OBJKEYS 返回一个字符串数组（对象的键）
	keys, err := cmd.StringSlice()
	if err != nil {
		// ***** 关键修正：这里声明 redisErr 为接口类型本身，而不是指向接口的指针 *****
		var redisErr redisHd.Error // 声明一个 redisHd.Error 接口类型的变量

		// errors.As 会尝试将 err 转换为 redisHd.Error 接口类型，并赋值给 redisErr
		if errors.As(err, &redisErr) { // 传递 redisErr 变量的地址
			// 现在 redisErr 是一个 redisHd.Error 接口类型的值，
			// 你可以安全地调用它的方法，包括继承自 error 接口的 Error() 方法
			if redisErr.Error() == "ERR wrong type" {
				return nil, nil // 针对“ERR wrong type”错误，返回 nil 切片和 nil 错误
			}
			// 你也可以调用自定义的 RedisError() 方法，如果需要的话
			// redisErr.RedisError()
		}
		// 如果 err 不是 redisHd.Error 接口类型，或者 errors.As 失败，则返回原始错误
		return nil, err
	}
	return keys, nil
}

// OBJLEN 获取JSON 对象中键的数量，如果匹配的 JSON 值不是对象，则为 -1
// key: Redis 中的键名。
// path: 可选的 JSON Path。如果为空字符串，则返回根对象的键。
// 返回值：一个包含对象键的字符串切片。如果键、路径不存在或路径对应的不是对象，则返回 nil 切片和 nil 错误。
func (that *KRedisCluster) JsonObjLen(key string, path string) ([]int64, error) {
	args := make([]any, 0, 3)
	args = append(args, "JSON.OBJLEN", key)
	if path != "" { // 如果 path 不为空，则添加到参数中
		args = append(args, path)
	}

	cmd := that.Client.Do(that.ctx.Context(), args...)

	if err := cmd.Err(); err != nil {
		if err == redisHd.Nil {
			return nil, nil // 键或路径不存在, 返回 nil 切片和 nil 错误
			// return nil, err
		}
		return nil, err
	}

	result := cmd.Val()
	switch retVal := result.(type) { // 使用类型断言和 switch 语句来处理不同类型的值
	case int64:
		return []int64{retVal}, nil
	case []any:
		array := retVal
		lenArray := make([]int64, 0, len(array))
		for _, v := range array {
			switch val := v.(type) {
			case int64:
				lenArray = append(lenArray, val)
			case nil:
				lenArray = append(lenArray, -1)
			default:
				return nil, fmt.Errorf("unexpected type: %v", v)
			}
		}
		return lenArray, nil
	default:
		return nil, fmt.Errorf("unexpected type: %v", retVal)
	}

}

func (that *KRedisCluster) Pipeline() redisHd.Pipeliner {
	return that.Client.Pipeline()
}

func (that *KRedisCluster) Dump(key string) (string, error) {
	return that.Client.Dump(that.ctx.Context(), key).Result()
}

func (that *KRedisCluster) RestoreReplace(key string, ttl time.Duration, value string) (string, error) {
	return that.Client.RestoreReplace(that.ctx.Context(), key, ttl, value).Result()
}

func (that *KRedisCluster) Restore(key string, ttl time.Duration, value string) (string, error) {
	return that.Client.Restore(that.ctx.Context(), key, ttl, value).Result()
}

// 删除一批key
func (that *KRedisCluster) Del(keys ...string) (int64, error) {
	return that.Client.Del(that.ctx.Context(), keys...).Result()
}

// 探测服务是否正常
func (that *KRedisCluster) Ping() bool {
	_, err := that.Client.Ping(that.ctx.Context()).Result()
	return nil == err
}

func (that *KRedisCluster) ScanMatch(limit int, aboutTypes []string, ignoreKeys []string, includeKeys []string, needDel bool, logf klogger.AppLogFuncWithTag) ([]*RedisRecord, error) {
	cursor := uint64(0)
	allKeys := make([]string, 0, 50000)

	count := 0
	for {
		var keys []string
		err := error(nil)
		keys, cursor, err = that.Client.Scan(that.ctx.Context(), cursor, "", int64(limit)).Result()
		if nil != err {
			return nil, err
		}

		count += len(keys)
		// logf(logger.InfoLevel, "scan %d keys, limit: %d, cursor: %d", count, limit, cursor)
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			// 扫描完成
			break
		}
	}

	dataList := make([]*RedisRecord, 0, limit)
	for _, key := range allKeys {
		// 黑名单过滤, 以对应关键字的key被过滤掉, 支持(pattern*, *pattern, pattern) 三种格式的匹配规则
		if MatchFilter(ignoreKeys, key) {
			continue
		}

		// logf(logger.InfoLevel, "idx: %d, filter ignore key:%s, type:%s", idx, key)

		// 如果有白名单, 则启用白名单规则, 不在白名单的被过滤掉, 白名单优先级低于黑名单, 支持(pattern*, *pattern, pattern) 三种格式的匹配规则
		if len(includeKeys) > 0 {
			if !MatchFilter(includeKeys, key) {
				continue
			}
		}

		dataType, err := that.Type(key)
		if nil != err {
			return nil, err
		}

		// if key == "windRtEvent:DTNXJK:HSBFC:Q1:W009" {
		// 	logf(logger.InfoLevel, "key:%s, type:%s", key, dataType)
		// }

		// logf(logger.InfoLevel, "idx: %d, key:%s, type:%s", idx, key, dataType)
		if !kslices.Contains(aboutTypes, strings.ToLower(dataType)) { //过滤出需要的数据类型
			continue
		}

		ttl, err := that.PTTL(key)
		if nil != err {
			return nil, err
		}

		data, err := that.Dump(key)
		if nil != err {
			return nil, err
		}

		dataList = append(dataList, &RedisRecord{Key: key, PTtl: ttl, DataType: dataType, Data: data})
	}

	return dataList, nil
}

func (that *KRedisCluster) Scan(limit int, aboutTypes []string, ignoreKeys []string, includeKeys []string, needDel bool, logf klogger.AppLogFuncWithTag) ([]*RedisRecord, error) {
	cursor := uint64(0)
	allKeys := make([]string, 0, 50000)

	// var m runtime.MemStats

	count := 0
	for {
		var keys []string
		err := error(nil)
		keys, cursor, err = that.Client.Scan(that.ctx.Context(), cursor, "", int64(limit)).Result()
		if nil != err {
			return nil, err
		}

		count += len(keys)
		// logf(logger.DebugLevel,"scan %d keys, limit: %d, cursor: %d", count, limit, cursor)
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			// 扫描完成
			break
		}
	}
	//index := int64(0)

	// runtime.ReadMemStats(&m)
	// logf(logger.DebugLevel,"%+v, os %d\n", m, m.Sys)

	dataList := make([]*RedisRecord, 0, limit)
	for _, key := range allKeys {
		// 黑名单过滤
		if kslices.Contains(ignoreKeys, key) { //过滤掉不需要的key
			continue
		}

		// 如果有白名单, 则启用白名单规则, 不在白名单的被过滤掉, 白名单优先级低于黑名单
		if len(includeKeys) > 0 {
			if !kslices.Contains(includeKeys, key) {
				continue
			}
		}

		dataType, err := that.Type(key)
		if nil != err {
			return nil, err
		}

		//logf(logger.DebugLevel,"idx: %d, key:%s, type:%s", idx, key, dataType)
		if !kslices.Contains(aboutTypes, strings.ToLower(dataType)) { //过滤出需要的数据类型
			continue
		}

		ttl, err := that.PTTL(key)
		if nil != err {
			return nil, err
		}

		data, err := that.Dump(key)
		if nil != err {
			return nil, err
		}
		// logf(logger.DebugLevel,"key:%s, type:%s, ttl:%d", key, dataType, ttl)
		//index++ //golang中 `++` 与 `--` 运算符只能作为语句存在, 不能作为表达式, 个小垃圾
		dataList = append(dataList, &RedisRecord{Key: key, PTtl: ttl, DataType: dataType, Data: data})
	}

	// runtime.ReadMemStats(&m)
	// logf(logger.DebugLevel,"%+v, os %d\n", m, m.Sys)

	return dataList, nil
}

// 向指定topic发布消息
func (that *KRedisCluster) Publish(topic string, payload any) error {
	return that.Client.Publish(that.ctx.Context(), topic, payload).Err()
}

// 使用pipeline 向指定topic发布多条消息
func (that *KRedisCluster) PublishArray(messages []*RedisMessage) []error {
	if len(messages) == 0 {
		return nil
	}

	pipeline := that.Client.Pipeline()
	for _, msg := range messages {
		pipeline.Publish(that.ctx.Context(), msg.Topic, msg.Message)
	}

	cmders, err := pipeline.Exec(that.ctx.Context())
	if nil == err {
		return nil // 整批完全成功
	}

	errs := make([]error, 0, len(messages))
	for _, cmd := range cmders {
		if nil != cmd.Err() {
			errs = append(errs, cmd.Err())
		}
	}
	return errs
}

// 使用pipeline 向指定topic发布多条消息
func (that *KRedisCluster) PublishArrayWithCtx(ctx *kcontext.ContextNode, messages []*RedisMessage) []error {
	if len(messages) == 0 {
		return nil
	}

	pipeline := that.Client.Pipeline()
	for _, msg := range messages {
		pipeline.Publish(ctx.Context(), msg.Topic, msg.Message)
	}

	cmders, err := pipeline.Exec(ctx.Context())
	if nil == err {
		return nil // 整批完全成功
	}

	errs := make([]error, 0, len(messages))
	for _, cmd := range cmders {
		if nil != cmd.Err() {
			errs = append(errs, cmd.Err())
		}
	}
	return errs
}

// 从指定topic订阅消息, 底层API, 最好使用Subscribe替代
func (that *KRedisCluster) SyncSubscribeLow(callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
		defer pubsub.Close()

	forEnd: //这个标签
		for {
			message, err := pubsub.ReceiveMessage(that.ctx.Context())
			callback(err, message.Channel, message.Payload) // 开一个协程用于加工收到的消息

			select {
			case <-that.ctx.Context().Done():
				break forEnd
			default:
				continue
			}
		}
	}()

	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, 底层API, 最好使用Subscribe替代
func (that *KRedisCluster) SubscribeLow(callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
		defer pubsub.Close()

	forEnd: //这个标签
		for {
			message, err := pubsub.ReceiveMessage(that.ctx.Context())
			go callback(err, message.Channel, message.Payload) // 开一个协程用于加工收到的消息

			select {
			case <-that.ctx.Context().Done():
				break forEnd
			default:
				continue
			}
		}
	}()

	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, 底层API, 最好使用Subscribe替代
func (that *KRedisCluster) SyncPSubscribeLow(callback func(err error, topic string, payload any), topics ...string) {
	pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
	defer pubsub.Close()

forEnd: //这个标签
	for {
		message, err := pubsub.ReceiveMessage(that.ctx.Context())
		callback(err, message.Channel, message.Payload) // 开一个协程用于加工收到的消息

		select {
		case <-that.ctx.Context().Done():
			break forEnd
		default:
			continue
		}
	}

	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, 底层API, 最好使用Subscribe替代
func (that *KRedisCluster) PSubscribeLow(callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
		defer pubsub.Close()

	forEnd: //这个标签
		for {
			message, err := pubsub.ReceiveMessage(that.ctx.Context())
			go callback(err, message.Channel, message.Payload) // 开一个协程用于加工收到的消息

			select {
			case <-that.ctx.Context().Done():
				break forEnd
			default:
				continue
			}
		}
	}()

	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息
func (that *KRedisCluster) SyncSubscribeWithoutTimeout(callback func(err error, topic string, payload any), topics ...string) {
	pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
	ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
forEnd: //这个标签
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				goto END
			} else {
				callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
			}
		case <-that.ctx.Context().Done():
			break forEnd
		}
	}

	pubsub.Close()
	// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
	for msg := range ch {
		callback(nil, msg.Channel, msg.Payload)
	}

END:
	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息
func (that *KRedisCluster) SubscribeWithoutTimeout(callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
					goto END
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-that.ctx.Context().Done():
				break forEnd
			}
		}

		pubsub.Close()
		// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
		for msg := range ch {
			callback(nil, msg.Channel, msg.Payload)
		}

	END:
		callback(ErrUnSubscribe, "", nil)
	}()
}

// 从指定topic订阅消息, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) SyncSubscribe(timeout int, callback func(err error, topic string, payload any), topics ...string) {
	pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
	// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
	ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
forEnd: //这个标签
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				goto END
			} else {
				callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
			continue
		case <-that.ctx.Context().Done():
			break forEnd
		}
	}

	pubsub.Close()
	// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
	for msg := range ch {
		callback(nil, msg.Channel, msg.Payload)
	}

END:
	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) Subscribe(timeout int, callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.Subscribe(that.ctx.Context(), topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
					goto END
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-that.ctx.Context().Done():
				break forEnd
			}
		}

		pubsub.Close()
		// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
		for msg := range ch {
			callback(nil, msg.Channel, msg.Payload)
		}

	END:
		callback(ErrUnSubscribe, "", nil)
	}()
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; chanSize 最大允许队列大小, 如果< 1, 则为1; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) SyncPSubscribeWithChanSize(timeout int, chanSize int, callback func(err error, topic string, payload any), topics ...string) {
	pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
	// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
	if chanSize < 1 {
		chanSize = 1
	}
	ch := pubsub.Channel(redisHd.WithChannelSize(chanSize), redisHd.WithChannelHealthCheckInterval(time.Second*30))
forEnd: //这个标签
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				goto END
			} else {
				callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
			continue
		case <-that.ctx.Context().Done():
			break forEnd
		}
	}

	pubsub.Close()
	// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
	for msg := range ch {
		callback(nil, msg.Channel, msg.Payload)
	}

END:
	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; chanSize 最大允许队列大小, 如果< 1, 则为1; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) PSubscribeWithChanSize(timeout int, chanSize int, callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		if chanSize < 1 {
			chanSize = 1
		}
		ch := pubsub.Channel(redisHd.WithChannelSize(chanSize), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
					goto END
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-that.ctx.Context().Done():
				break forEnd
			}
		}

		pubsub.Close()
		// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
		for msg := range ch {
			callback(nil, msg.Channel, msg.Payload)
		}

	END:
		callback(ErrUnSubscribe, "", nil)
	}()
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) SyncPSubscribe(timeout int, callback func(err error, topic string, payload any), topics ...string) {
	pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
	// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
	ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
forEnd: //这个标签
	for {
		select {
		case message, ok := <-ch:
			if !ok {
				callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				goto END
			} else {
				callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
			continue
		case <-that.ctx.Context().Done():
			break forEnd
		}
	}
	pubsub.Close()
	// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
	for msg := range ch {
		callback(nil, msg.Channel, msg.Payload)
	}

END:
	callback(ErrUnSubscribe, "", nil)
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (that *KRedisCluster) PSubscribe(timeout int, callback func(err error, topic string, payload any), topics ...string) {
	go func() {
		pubsub := that.Client.PSubscribe(that.ctx.Context(), topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(ErrChannelClosed, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
					goto END
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-that.ctx.Context().Done():
				break forEnd
			}
		}
		pubsub.Close()
		// 此时 ch 已经被 close 了，range 会处理完 Buffer 里的数据后自动退出
		for msg := range ch {
			callback(nil, msg.Channel, msg.Payload)
		}

	END:
		callback(ErrUnSubscribe, "", nil)
	}()
}

func (that *KRedisCluster) Stop() {
	// that.CancelSubscribe()
	that.ctx.Cancel()
	that.ctx.Remove() // 移除上下文树中的节点
	that.Client.Close()
}

////////////////////////////////////
