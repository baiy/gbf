package rabbitmq

import (
	"time"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util"
)

type Queue struct {
	rabbitmq *Rabbitmq
	Logger   log.Logger
}

// Send 发送消息
// topic 主题  格式: name1.name2.name3....
func (q *Queue) Send(topic string, data []byte, options ...util.WithOption[Setting]) {
	op := Setting{}
	util.ApplyOption(&op, options...)
	q.rabbitmq.Publisher(data, append([]string{topic}, op.topics...), op.ttl)
}

// Listen 监听消息
// group 组名
// topic: 主题
func (q *Queue) Listen(group string, topic string, handler ConsumerHandler, options ...util.WithOption[Setting]) {
	op := Setting{}
	util.ApplyOption(&op, options...)
	q.rabbitmq.Consumer(group, append([]string{topic}, op.topics...), handler, op.concurrency)
}

type Setting struct {
	// 更多主题
	// 格式: name1.name2.name3....  支持通配符:* 例如: name1.name2.*
	topics []string

	// 消息过期时间
	ttl time.Duration

	// 消费并发数
	concurrency int
}

func WithTopics(topic ...string) util.WithOption[Setting] {
	return func(o *Setting) {
		o.topics = topic
	}
}

func WithTTL(ttl time.Duration) util.WithOption[Setting] {
	return func(o *Setting) {
		o.ttl = ttl
	}
}

func WithConcurrency(concurrency int) util.WithOption[Setting] {
	return func(o *Setting) {
		o.concurrency = concurrency
	}
}
