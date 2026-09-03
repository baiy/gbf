package cache

import (
	"github.com/baiy/gbf/util"
)

type Option struct {
	driver Driver
	prefix string

	// 启用本地二级缓存的key
	localKeys []string
}

func WithDriver(d Driver) util.WithOption[Option] {
	return func(o *Option) {
		o.driver = d
	}
}

func WithPrefix(prefix string) util.WithOption[Option] {
	return func(o *Option) {
		o.prefix = prefix
	}
}

func WithLocalKeys(keys ...string) util.WithOption[Option] {
	return func(o *Option) {
		o.localKeys = keys
	}
}
