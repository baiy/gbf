package redis

import (
	"cmp"
	"fmt"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/manager"
	"github.com/redis/go-redis/v9"
)

func create(mo manager.Option[Option]) (redis.UniversalClient, error) {
	option := mo.Value
	o := &redis.UniversalOptions{
		Addrs:      option.Addrs,
		Username:   cmp.Or(option.Username, ""),
		Password:   cmp.Or(option.Password, ""),
		MasterName: option.MasterName,
		DB:         option.Db,
	}
	if option.Pool != 0 {
		o.PoolSize = option.Pool
	}
	r := redis.NewUniversalClient(o)
	mo.Logger.Infof("redis connect addrs:%v db:%d", mo.Value.Addrs, mo.Value.Db)
	closer.Register(fmt.Sprintf("redis:%s:%s", mo.Key, mo.Value.Addrs), r)
	return r, nil
}

var Manager = manager.New("redis", create, nil, func() Option {
	return Option{Addrs: []string{":6379"}}
})
