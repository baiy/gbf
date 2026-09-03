package rabbitmq

import (
	"fmt"
	"time"

	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/errors"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util/slice"
	"github.com/baiy/gbf/worker"
	"github.com/wagslane/go-rabbitmq"
)

var (
	normalExchange  = "normal"
	delayedExchange = "exchange"
)

type Rabbitmq struct {
	prefix string
	logger log.Logger
	conn   *rabbitmq.Conn

	// 正常消息发布器
	normalPublisher *rabbitmq.Publisher
	// 延迟消息发布器
	delayedPublisher *rabbitmq.Publisher

	// 消费者
	consumers map[string]*rabbitmq.Consumer
}

// Publisher 发布消息
func (r *Rabbitmq) Publisher(value []byte, routingKeys []string, ttl time.Duration) {
	options := []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsPersistentDelivery, // 持久化消息
		rabbitmq.WithPublishOptionsMandatory,          // 消息无法路由时退回给生产者 保证 publisher.NotifyReturn 可用
	}
	routingKeys = slice.Map(routingKeys, func(_ int, routingKey string) string {
		return r.namer(routingKey)
	})
	var err error
	if ttl > 0 {
		err = r.delayedPublisher.Publish(
			value,
			routingKeys,
			append(
				options,
				rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table{
					"x-delay": uint32(ttl.Milliseconds()),
				}),
				rabbitmq.WithPublishOptionsExchange(r.namer(delayedExchange)),
			)...,
		)
	} else {
		err = r.normalPublisher.Publish(
			value,
			routingKeys,
			append(options, rabbitmq.WithPublishOptionsExchange(r.namer(normalExchange)))...,
		)
	}
	if err != nil {
		r.logger.Errorf("publisher error: %v %v => %s", err, routingKeys, value)
		return
	}
	r.logger.Infof("publisher success: %v => %s", routingKeys, value)
}

// Consumer 消费消息
func (r *Rabbitmq) Consumer(queue string, routingKeys []string, handler ConsumerHandler, concurrency int) {
	routingKeys = slice.Map(routingKeys, func(_ int, routingKey string) string {
		return r.namer(routingKey)
	})
	queue = r.namer(queue)
	worker.Async(func() {
		// 常驻 保证中断后能重启
		for {
			time.Sleep(time.Second * 5)
			if r.conn.IsClosed() {
				r.logger.Errorf("rabbitmq connection closed")
				return
			}

			consumer, ok := r.consumers[queue]
			if ok {
				consumer.Close()
				delete(r.consumers, queue)
				r.logger.Infof("consumer closed: %s", queue)
			}

			consumer, err := r.createConsumer(queue, routingKeys, concurrency)
			if err != nil {
				r.logger.Errorf("create consumer error: %v", err)
				continue
			}
			r.logger.Infof("create consumer : %s => %v", queue, routingKeys)
			r.consumers[queue] = consumer

			err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
				defer log.ElapsedTime(fmt.Sprintf("rabbitmq:consumer:%s", queue))()
				log.TrackReset()
				r.logger.Infof("consumer data: %s", string(d.Body))

				err = errors.PanicToError(func() error {
					return handler(d.Body)
				})
				if err != nil {
					r.logger.Errorf("handler error: %v", err)
				}
				return rabbitmq.Ack
			})
			if err != nil {
				r.logger.Errorf("consumer run error: %v", err)
				continue
			}
		}
	})
}

func (r *Rabbitmq) Close() error {
	r.normalPublisher.Close()
	for _, c := range r.consumers {
		c.Close()
	}
	return r.conn.Close()
}

func (r *Rabbitmq) createPublisher(options ...func(*rabbitmq.PublisherOptions)) (*rabbitmq.Publisher, error) {
	logger := Logger{r.logger.With("publisher")}
	options = append([]func(*rabbitmq.PublisherOptions){
		rabbitmq.WithPublisherOptionsLogger(logger),
		rabbitmq.WithPublisherOptionsExchangeDeclare, // 自动声明交换机
		rabbitmq.WithPublisherOptionsExchangeDurable, // 持久化交换机
	}, options...)
	publisher, err := rabbitmq.NewPublisher(
		r.conn,
		options...,
	)
	if err != nil {
		return nil, err
	}
	publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
		if c.Ack {
			logger.Infof("confirm success, delivery tag: %d", c.DeliveryTag)
			return
		}
		logger.Errorf("confirm failure, delivery tag: %d", c.DeliveryTag)
	})
	publisher.NotifyReturn(func(mq rabbitmq.Return) {
		if mq.Exchange == r.namer(delayedExchange) {
			// 延迟消息过滤
			return
		}
		logger.Errorf("notify return: [%s] => %s, reply code: %d", mq.RoutingKey, string(mq.Body), mq.ReplyCode)
	})
	return publisher, nil
}

func (r *Rabbitmq) createConsumer(queue string, routingKeys []string, concurrency int) (*rabbitmq.Consumer, error) {
	logger := Logger{r.logger.With("consumer")}

	options := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsLogger(logger),
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsQueueArgs(rabbitmq.Table{
			// 空闲超时，自动删除队列
			"x-expires": (2 * time.Hour).Milliseconds(),
		}),
	}
	if concurrency > 1 {
		options = append(options, rabbitmq.WithConsumerOptionsConcurrency(concurrency))
		options = append(options, rabbitmq.WithConsumerOptionsQOSPrefetch(concurrency*3))
	}

	bindings := slice.Map(routingKeys, func(_ int, routingKey string) rabbitmq.Binding {
		return rabbitmq.Binding{
			NoWait:     false,
			Declare:    true,
			RoutingKey: routingKey,
		}
	})

	options = append(
		options,
		rabbitmq.WithConsumerOptionsExchangeOptions(
			rabbitmq.ExchangeOptions{
				Name:     r.namer(normalExchange),
				Declare:  false,
				Bindings: bindings,
			},
		),
		rabbitmq.WithConsumerOptionsExchangeOptions(
			rabbitmq.ExchangeOptions{
				Name:     r.namer(delayedExchange),
				Declare:  false,
				Bindings: bindings,
			},
		),
	)

	return rabbitmq.NewConsumer(
		r.conn,
		queue,
		options...,
	)
}

func newRabbitmq(conn *rabbitmq.Conn, logger log.Logger, prefix string) (*Rabbitmq, error) {
	r := &Rabbitmq{
		prefix:    prefix,
		conn:      conn,
		logger:    logger,
		consumers: make(map[string]*rabbitmq.Consumer),
	}

	// 创建正常消息发布器
	normalPublisher, err := r.createPublisher(
		rabbitmq.WithPublisherOptionsExchangeName(r.namer(normalExchange)),
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
	)
	if err != nil {
		return nil, err
	}

	// 创建延迟消息发布器
	delayedPublisher, err := r.createPublisher(
		rabbitmq.WithPublisherOptionsExchangeName(r.namer(delayedExchange)),
		rabbitmq.WithPublisherOptionsExchangeKind("x-delayed-message"),
		rabbitmq.WithPublisherOptionsExchangeArgs(rabbitmq.Table{"x-delayed-type": "topic"}),
	)
	if err != nil {
		return nil, err
	}
	r.normalPublisher = normalPublisher
	r.delayedPublisher = delayedPublisher

	return r, nil
}

func (r *Rabbitmq) namer(v string) string {
	if r.prefix == "" {
		return fmt.Sprintf("%s.%s", constant.BaseName, v)
	}
	return fmt.Sprintf("%s.%s.%s", constant.BaseName, r.prefix, v)
}
