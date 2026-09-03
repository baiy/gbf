package redis

import (
	"github.com/baiy/gbf/util"
)

type Option struct {
	// 连接地址 masterName 配置查看 https://redis.uptrace.dev/zh/guide/universal.html
	Addrs []string `yaml:"addrs"`
	// 哨兵模式
	MasterName string `yaml:"masterName"`

	// 用户名
	Username string `yaml:"username"`
	// 密码
	Password string `yaml:"password"`
	// 数据库
	Db int `yaml:"db"`

	// 是否启用SSL
	Ssl bool `yaml:"ssl"`

	// https://redis.uptrace.dev/zh/guide/go-redis-debugging.html#%E8%BF%9E%E6%8E%A5%E6%B1%A0%E5%A4%A7%E5%B0%8F
	// 连接池 一般不用配置内部有默认连接池
	Pool int `yaml:"pool"`
}

func WithAddr(addrs ...string) util.WithOption[Option] {
	return func(o *Option) {
		o.Addrs = addrs
	}
}

func WithMasterName(name string) util.WithOption[Option] {
	return func(o *Option) {
		o.MasterName = name
	}
}

func WithUsername(username string) util.WithOption[Option] {
	return func(o *Option) {
		o.Username = username
	}
}

func WithPassword(password string) util.WithOption[Option] {
	return func(o *Option) {
		o.Password = password
	}
}

func WithDb(db int) util.WithOption[Option] {
	return func(o *Option) {
		o.Db = db
	}
}

func WithPoolSize(size int) util.WithOption[Option] {
	return func(o *Option) {
		o.Pool = size
	}
}

func WithSsl() util.WithOption[Option] {
	return func(o *Option) {
		o.Ssl = true
	}
}
