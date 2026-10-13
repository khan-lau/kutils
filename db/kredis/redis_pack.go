package kredis

import "time"

// RedisMessageRecvFunc 接收消息的回调函数
//
//  参数:
//
//  @param err 错误信息
//  @param topic 主题
//  @param recvTime 接收时间
//  @param payload 消息内容
type RedisMessageRecvFunc func(err error, topic string, recvTime int64, payload any)

type RedisRecord struct {
	Key      string
	DataType string
	PTtl     time.Duration
	Data     string
}

type RedisMessage struct {
	Topic    string // 主题
	RecvTime int64  // 接收时间
	Message  string // 消息内容
}

type RedisMessagePackConstraints interface {
	RedisRecord | RedisMessage
	// constraints.Integer | constraints.Float
}

type KRedisClient interface {
	KRedis | KRedisCluster
}
