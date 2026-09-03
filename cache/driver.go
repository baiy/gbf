package cache

import (
	"time"

	"github.com/baiy/gbf/log"
)

type Driver interface {
	Get(name string) (string, bool)
	Exist(name string) bool
	Set(name string, value string, ttl time.Duration)
	Del(name string)
	TTL(name string) time.Duration // -1 永久有效, -2 数据不存在

	// LocalIsAvailable 本地二级缓存是否可用 分布式缓存一般可用 单机缓存一般不可用
	LocalIsAvailable() bool

	SetLogger(log.Logger)
}

type Base struct {
	Logger log.Logger
}

func (b *Base) SetLogger(l log.Logger) {
	b.Logger = l
}
