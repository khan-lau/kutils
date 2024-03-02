package kredis

import "time"

type RedisRecord struct {
	Key      string
	DataType string
	PTtl     time.Duration
	Data     string
}

type RedisMessage struct {
	Topic   string
	Message string
}

type RedisMessagePackConstraints interface {
	RedisRecord | RedisMessage
	// constraints.Integer | constraints.Float
}
