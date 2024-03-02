package kredis

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	kslices "github.com/khan-lau/kutils/container/kslices"
	"github.com/khan-lau/kutils/logger"

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

// 设置某个key的值, 并指定ttl
func (mr *KRedis) Set(key string, value interface{}, duration time.Duration) (bool, error) {
	err := mr.Client.Set(mr.ctx, key, value, duration).Err()
	if err != nil {
		return false, err
	}
	return true, nil
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

func (mr *KRedis) Scan(limit int, aboutTypes []string, ignoreKeys []string, includeKeys []string, needDel bool, logf logger.AppLogFunc) ([]*RedisRecord, error) {
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
