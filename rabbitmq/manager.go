package rabbitmq

import (
	"fmt"
	"net/url"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/manager"
	"github.com/wagslane/go-rabbitmq"
)

func create(mo manager.Option[Option]) (*Queue, error) {
	encodedUser := url.QueryEscape(mo.Value.username)
	encodedPass := url.QueryEscape(mo.Value.password)
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:%d/%s", encodedUser, encodedPass, mo.Value.host, mo.Value.port, mo.Value.virtualHost),
		rabbitmq.WithConnectionOptionsLogger(Logger{mo.Logger}),
		rabbitmq.WithConnectionOptionsBaseReconnectInterval(time.Second*5), // 自动重连
	)
	if err != nil {
		return nil, err
	}
	rmq, err := newRabbitmq(conn, mo.Logger.With("rabbitmq"), mo.Value.prefix)
	if err != nil {
		return nil, err
	}
	r := &Queue{
		Logger:   mo.Logger,
		rabbitmq: rmq,
	}
	mo.Logger.Infof("rabbitmq connect %s@%s:%d/%s", mo.Value.username, mo.Value.host, mo.Value.port, mo.Value.virtualHost)
	closer.Register(fmt.Sprintf("rabbitmq:%s", mo.Key), rmq)
	return r, nil
}

var Manager = manager.New[*Queue, Option](
	"rabbitmq",
	create,
	nil,
	func() Option {
		return Option{
			host:     "127.0.0.1",
			username: "admin",
			password: "admin",
			port:     5672,
		}
	},
)
