package kredis

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	kslices "github.com/khan-lau/kutils/container/kslices"
	"github.com/khan-lau/kutils/klogger"

	redisHd "github.com/redis/go-redis/v9"
)

type Empty struct{}

type KRedis struct {
	Client *redisHd.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func redisOnConnect(ctx context.Context, cn *redisHd.Conn) error {
	return nil
}

func NewKRedis(ctx context.Context, host string, port int, user string, password string, dbNum int) *KRedis {
	client := redisHd.NewClient(&redisHd.Options{
		Addr:            host + ":" + strconv.Itoa(port),
		Username:        user,     // redis 6.0以上版本
		Password:        password, // 没有密码，默认值
		DB:              dbNum,    // 默认DB 0
		MaxRetries:      3,        // 自动重连3次, 失败后报错
		DialTimeout:     10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: 30 * time.Second, // 链路最大空闲时间
		OnConnect:       redisOnConnect,
	})

	subCtx, subCancel := context.WithCancel(ctx)

	return &KRedis{Client: client, ctx: subCtx, cancel: subCancel}
}

// 执行指令
func (mr *KRedis) Do(args ...interface{}) (interface{}, error) {
	val, err := mr.Client.Do(mr.ctx, args...).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 获取一个key的值
func (mr *KRedis) Get(key string) (interface{}, error) {
	val, err := mr.Client.Do(mr.ctx, "GET", key).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 设置某个key的值, 并指定ttl
func (mr *KRedis) Set(key string, value interface{}, duration time.Duration) (bool, error) {
	err := mr.Client.Set(mr.ctx, key, value, duration).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 判断某个key是否存在
func (mr *KRedis) Exist(key string) (bool, error) {
	_, err := mr.Client.Get(mr.ctx, key).Result()
	if err == redisHd.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// 获取一个key的hash字段的值
func (that *KRedis) HGet(key string, field string) (interface{}, error) {
	val, err := that.Client.HGet(that.ctx, key, field).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// 设置一个key的hash字段的值
func (that *KRedis) HSet(key string, field string, value interface{}) error {
	err := that.Client.HSet(that.ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取一个key的hash字段的值列表
func (that *KRedis) HGetAll(key string) (map[string]string, error) {
	valMap, err := that.Client.HGetAll(that.ctx, key).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return valMap, nil
}

// 设置一个key的hash字段的值列表
func (that *KRedis) HSetAll(key string, fields map[string]interface{}) error {
	err := that.Client.HSet(that.ctx, key, fields).Err()
	if err != nil {
		return err
	}
	return nil
}

// 判断一个key的hash字段是否存在
func (that *KRedis) HExists(key string, field string) (bool, error) {
	isExists, err := that.Client.HExists(that.ctx, key, field).Result()
	if nil != err {
		return false, err
	}
	return isExists, nil
}

// 获取一个key的hash字段的数量
func (that *KRedis) HLen(key string) (int64, error) {
	val, err := that.Client.HLen(that.ctx, key).Result()
	if nil != err {
		return 0, err
	}
	return val, nil
}

// 获取一个key的hash字段的所有Key
func (that *KRedis) HKeys(ctx context.Context, key string) ([]string, error) {
	array, err := that.Client.HKeys(ctx, key).Result()
	if nil != err {
		return nil, err
	}
	return array, nil
}

// 获取一个key的hash字段的所有值
func (that *KRedis) HVals(ctx context.Context, key string) ([]string, error) {
	array, err := that.Client.HVals(ctx, key).Result()
	if nil != err {
		return nil, err
	}
	return array, nil
}

// 设置一个key的hash字段的值列表, 如果不存在则创建
func (that *KRedis) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	isExists, err := that.Client.HSetNX(ctx, key, field, value).Result()
	if nil != err {
		return false, err
	}
	return isExists, nil
}

// 删除一个key的hash字段的值列表
func (that *KRedis) HDel(key string, fields []string) error {
	_, err := that.Client.HDel(that.ctx, key, fields...).Result()
	if nil != err {
		return err
	}
	return nil
}

// 获取一个key的hash字段的值列表
func (that *KRedis) HMGet(key string, fields []string) ([]interface{}, error) {
	valMap, err := that.Client.HMGet(that.ctx, key, fields...).Result()
	if err == redisHd.Nil {
		return nil, nil
	} else if nil != err {
		return nil, err
	}
	return valMap, nil
}

// 设置一个key的hash字段的值列表, 如果不存在则创建
func (that *KRedis) HMSet(key string, fields map[string]interface{}) error {
	err := that.Client.HMSet(that.ctx, key, fields).Err()
	if nil != err {
		return err
	}
	return nil
}

// 从列表左边插入数据
func (that *KRedis) LPush(key string, values ...interface{}) (int64, error) {
	return that.Client.LPush(that.ctx, key, values...).Result()
}

// 从列表左边插入数据, 如果不存在则不插入数据
func (that *KRedis) LPushX(key string, values ...interface{}) (int64, error) {
	return that.Client.LPushX(that.ctx, key, values...).Result()
}

// 从列表右边插入数据
func (that *KRedis) RPush(key string, values ...interface{}) (int64, error) {
	return that.Client.RPush(that.ctx, key, values...).Result()
}

// 从列表右边插入数据, 如果不存在则不插入数据
func (that *KRedis) RPushX(key string, values ...interface{}) (int64, error) {
	return that.Client.RPushX(that.ctx, key, values...).Result()
}

// 从列表左边弹出数据
func (that *KRedis) LPop(key string) (string, error) {
	return that.Client.LPop(that.ctx, key).Result()
}

// 从列表右边弹出数据
func (that *KRedis) RPop(key string) (string, error) {
	return that.Client.RPop(that.ctx, key).Result()
}

// 返回列表的一个范围内的数据，也可以返回全部数据
func (that *KRedis) LRange(key string, start int64, stop int64) ([]string, error) {
	return that.Client.LRange(that.ctx, key, start, stop).Result()
}

// 返回列表的大小
func (that *KRedis) LLen(key string) (int64, error) {
	return that.Client.LLen(that.ctx, key).Result()
}

func (that *KRedis) LTrim(key string, start int64, stop int64) error {
	return that.Client.LTrim(that.ctx, key, start, stop).Err()
}

func (that *KRedis) LSet(key string, index int64, value interface{}) error {
	return that.Client.LSet(that.ctx, key, index, value).Err()
}

// 删除列表中的数据
func (that *KRedis) LRem(key string, count int64, value interface{}) (int64, error) {
	return that.Client.LRem(that.ctx, key, count, value).Result()
}

// 根据索引坐标，查询列表中的数据
func (that *KRedis) LIndex(key string, index int64) (string, error) {
	return that.Client.LIndex(that.ctx, key, index).Result()
}

// 在指定位置插入数据，在头部插入用"before"，尾部插入用"after"
func (that *KRedis) LInsert(key string, position string, pivot interface{}, value interface{}) (int64, error) {
	return that.Client.LInsert(that.ctx, key, position, pivot, value).Result()
}

func (that *KRedis) SAdd(key string, members ...interface{}) (int64, error) {
	return that.Client.SAdd(that.ctx, key, members...).Result()
}

func (that *KRedis) SMembers(key string) ([]string, error) {
	return that.Client.SMembers(that.ctx, key).Result()
}

func (that *KRedis) SRem(key string, members ...interface{}) (int64, error) {
	return that.Client.SRem(that.ctx, key, members...).Result()
}

func (that *KRedis) SIsMember(key string, member interface{}) (bool, error) {
	return that.Client.SIsMember(that.ctx, key, member).Result()
}

func (that *KRedis) SCard(key string) (int64, error) {
	return that.Client.SCard(that.ctx, key).Result()
}

func (that *KRedis) SPop(key string) (string, error) {
	return that.Client.SPop(that.ctx, key).Result()
}

func (that *KRedis) SPopN(key string, count int64) ([]string, error) {
	return that.Client.SPopN(that.ctx, key, count).Result()
}

func (that *KRedis) SUnion(keys ...string) ([]string, error) {
	return that.Client.SUnion(that.ctx, keys...).Result()
}

func (that *KRedis) SUnionStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SUnionStore(that.ctx, destKey, keys...).Result()
}

func (that *KRedis) SInter(keys ...string) ([]string, error) {
	return that.Client.SInter(that.ctx, keys...).Result()
}

func (that *KRedis) SInterStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SInterStore(that.ctx, destKey, keys...).Result()
}

func (that *KRedis) SDiff(keys ...string) ([]string, error) {
	return that.Client.SDiff(that.ctx, keys...).Result()
}

func (that *KRedis) SDiffStore(destKey string, keys ...string) (int64, error) {
	return that.Client.SDiffStore(that.ctx, destKey, keys...).Result()
}

func (that *KRedis) SMove(source, destination string, member interface{}) (bool, error) {
	return that.Client.SMove(that.ctx, source, destination, member).Result()
}

func (that *KRedis) SRandMember(key string) (string, error) {
	return that.Client.SRandMember(that.ctx, key).Result()
}

// 获取一个key的数据类型, 数据类型全小写
func (mr *KRedis) Type(key string) (string, error) {
	dataType, err := mr.Client.Type(mr.ctx, key).Result()
	if nil != err {
		return "", err
	}
	return strings.ToLower(dataType), nil
}

// 返回一个Key的过期时间, 单位为毫秒
func (mr *KRedis) PTTL(key string) (time.Duration, error) {
	return mr.Client.PTTL(mr.ctx, key).Result()
}

// 返回一个Key的过期时间, 单位为秒
func (mr *KRedis) TTL(key string) (time.Duration, error) {
	return mr.Client.TTL(mr.ctx, key).Result()
}

func (mr *KRedis) Pipeline() redisHd.Pipeliner {
	return mr.Client.Pipeline()
}

func (mr *KRedis) Dump(key string) (string, error) {
	return mr.Client.Dump(mr.ctx, key).Result()
}

func (mr *KRedis) RestoreReplace(key string, ttl time.Duration, value string) (string, error) {
	return mr.Client.RestoreReplace(mr.ctx, key, ttl, value).Result()
}

func (mr *KRedis) Restore(key string, ttl time.Duration, value string) (string, error) {
	return mr.Client.Restore(mr.ctx, key, ttl, value).Result()
}

// 删除一批key
func (mr *KRedis) Del(keys ...string) (int64, error) {
	return mr.Client.Del(mr.ctx, keys...).Result()
}

// 探测服务是否正常
func (mr *KRedis) Ping() bool {
	_, err := mr.Client.Ping(mr.ctx).Result()
	return nil == err
}

func (mr *KRedis) ScanMatch(limit int, aboutTypes []string, ignoreKeys []string, includeKeys []string, needDel bool, logf klogger.AppLogFunc) ([]*RedisRecord, error) {
	cursor := uint64(0)
	allKeys := make([]string, 0, 50000)

	count := 0
	for {
		var keys []string
		err := error(nil)
		keys, cursor, err = mr.Client.Scan(mr.ctx, cursor, "", int64(limit)).Result()
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

		dataType, err := mr.Type(key)
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

		ttl, err := mr.PTTL(key)
		if nil != err {
			return nil, err
		}

		data, err := mr.Dump(key)
		if nil != err {
			return nil, err
		}

		dataList = append(dataList, &RedisRecord{Key: key, PTtl: ttl, DataType: dataType, Data: data})
	}

	return dataList, nil
}

func (mr *KRedis) Scan(limit int, aboutTypes []string, ignoreKeys []string, includeKeys []string, needDel bool, logf klogger.AppLogFunc) ([]*RedisRecord, error) {
	cursor := uint64(0)
	allKeys := make([]string, 0, 50000)

	// var m runtime.MemStats

	count := 0
	for {
		var keys []string
		err := error(nil)
		keys, cursor, err = mr.Client.Scan(mr.ctx, cursor, "", int64(limit)).Result()
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
			if !kslices.Contains[string](includeKeys, key) {
				continue
			}
		}

		dataType, err := mr.Type(key)
		if nil != err {
			return nil, err
		}

		//logf(logger.DebugLevel,"idx: %d, key:%s, type:%s", idx, key, dataType)
		if !kslices.Contains(aboutTypes, strings.ToLower(dataType)) { //过滤出需要的数据类型
			continue
		}

		ttl, err := mr.PTTL(key)
		if nil != err {
			return nil, err
		}

		data, err := mr.Dump(key)
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
func (mr *KRedis) Publish(topic string, payload interface{}) error {
	return mr.Client.Publish(mr.ctx, topic, payload).Err()
}

// 从指定topic订阅消息, 底层API, 最好使用Subscribe替代
func (mr *KRedis) SubscribeLow(callback func(err error, topic string, payload interface{}), topics ...string) {
	go func() {
		pubsub := mr.Client.Subscribe(mr.ctx, topics...)
		defer pubsub.Close()

	forEnd: //这个标签
		for {
			message, err := pubsub.ReceiveMessage(mr.ctx)
			go callback(err, message.Channel, message.Payload) // 开一个协程用于加工收到的消息

			select {
			case <-mr.ctx.Done():
				break forEnd
			default:
				continue
			}
		}
	}()

	callback(errors.New("func UnSubscribe be called"), "", nil)
}

// 从指定topic订阅消息
func (mr *KRedis) SubscribeWithoutTimeout(callback func(err error, topic string, payload interface{}), topics ...string) {
	go func() {
		pubsub := mr.Client.Subscribe(mr.ctx, topics...)
		defer pubsub.Close()

		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(errors.New("channel be closed"), message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-mr.ctx.Done():
				break forEnd
			}
		}

		callback(errors.New("func UnSubscribe be called"), "", nil)
	}()
}

// 从指定topic订阅消息, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (mr *KRedis) Subscribe(timeout int, callback func(err error, topic string, payload interface{}), topics ...string) {
	go func() {
		pubsub := mr.Client.Subscribe(mr.ctx, topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		defer pubsub.Close()

		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(errors.New("channel be closed"), message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-mr.ctx.Done():
				break forEnd

			}
		}

		callback(errors.New("func UnSubscribe be called"), "", nil)
	}()
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; chanSize 最大允许队列大小, 如果< 100, 则为100; callback为接收消息的回调函数; topics为需要订阅的topic
func (mr *KRedis) PSubscribeWithChanSize(timeout int, chanSize int, callback func(err error, topic string, payload interface{}), topics ...string) {
	go func() {
		pubsub := mr.Client.PSubscribe(mr.ctx, topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		defer pubsub.Close()
		if chanSize < 100 {
			chanSize = 100
		}
		ch := pubsub.Channel(redisHd.WithChannelSize(chanSize), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(errors.New("channel be closed"), message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-mr.ctx.Done():
				break forEnd

			}
		}

		callback(errors.New("func UnSubscribe be called"), "", nil)
	}()
}

// 从指定topic订阅消息, topic支持通配符, timeout 设置轮询超时时间, 单位ms; callback为接收消息的回调函数; topics为需要订阅的topic
func (mr *KRedis) PSubscribe(timeout int, callback func(err error, topic string, payload interface{}), topics ...string) {
	go func() {
		pubsub := mr.Client.PSubscribe(mr.ctx, topics...)
		// pubsub.Unsubscribe(mr.ctx, "xxx") //不关闭订阅的情况下取消订阅
		defer pubsub.Close()

		ch := pubsub.Channel(redisHd.WithChannelSize(100), redisHd.WithChannelHealthCheckInterval(time.Second*30))
	forEnd: //这个标签
		for {
			select {
			case message, ok := <-ch:
				if !ok {
					go callback(errors.New("channel be closed"), message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				} else {
					go callback(nil, message.Channel, message.Payload) // 开一个协程用于加工收到的消息
				}
			case <-time.After(time.Duration(timeout) * time.Millisecond): //上面的ch如果一直没数据会阻塞，那么select也会检测其他case条件，检测到后timeout指定毫秒超时
				continue
			case <-mr.ctx.Done():
				break forEnd

			}
		}

		callback(errors.New("func UnSubscribe be called"), "", nil)
	}()
}

func (mr *KRedis) Stop() {
	// mr.CancelSubscribe()
	mr.cancel()
	mr.Client.Close()
}

////////////////////////////////////

func MatchFilter(patterns []string, key string) bool {
	blacklistIndex := kslices.IndexFunc(patterns, func(item string) bool {
		length := len(item)
		if length == 0 || len(key) == 0 {
			return false
		}
		if item[:1] == "*" {
			return strings.HasSuffix(key, item[1:length])
		} else if item[length-1:] == "*" {
			return strings.HasPrefix(key, item[0:length-1])
		} else {
			return key == item
		}
	})
	return blacklistIndex >= 0
}
