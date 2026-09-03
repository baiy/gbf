package rabbitmq

import "github.com/baiy/gbf/util"

type ConsumerHandler func(data []byte) error

type Option struct {
	username    string
	password    string
	host        string
	port        uint
	virtualHost string

	// 相关队列/交换机名称前缀 用于区分不同项目
	prefix string
}

func WithUsername(v string) util.WithOption[Option] {
	return func(o *Option) {
		o.username = v
	}
}

func WithPassword(v string) util.WithOption[Option] {
	return func(o *Option) {
		o.password = v
	}
}

func WithHost(v string) util.WithOption[Option] {
	return func(o *Option) {
		o.host = v
	}
}

func WithPort(v uint) util.WithOption[Option] {
	return func(o *Option) {
		o.port = v
	}
}
func WithPrefix(v string) util.WithOption[Option] {
	return func(o *Option) {
		o.prefix = v
	}
}

func WithVirtualHost(v string) util.WithOption[Option] {
	return func(o *Option) {
		o.virtualHost = v
	}
}
