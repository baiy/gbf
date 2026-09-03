package redis

import (
	"context"
	"errors"
	"time"

	"github.com/baiy/gbf/cache"
	"github.com/baiy/gbf/redis"
	goRedis "github.com/redis/go-redis/v9"
)

type Redis struct {
	_client  goRedis.UniversalClient
	redisKey string
	cache.Base
}

func (r *Redis) Get(name string) (string, bool) {
	v, err := r.client().Get(context.Background(), name).Result()
	if err != nil {
		if errors.Is(err, goRedis.Nil) {
			return "", false
		}
		panic(err)
	}
	return v, true
}

func (r *Redis) Exist(name string) bool {
	v, err := r.client().Exists(context.Background(), name).Result()
	if err != nil {
		panic(err)
	}
	return v > 0
}

func (r *Redis) Set(name string, value string, ttl time.Duration) {
	err := r.client().Set(context.Background(), name, value, ttl).Err()
	if err != nil {
		panic(err)
	}
	return
}

func (r *Redis) Del(name string) {
	err := r.client().Del(context.Background(), name).Err()
	if err != nil {
		panic(err)
	}
}

func (r *Redis) TTL(name string) time.Duration {
	expire, err := r.client().TTL(context.Background(), name).Result()
	if err != nil {
		panic(err)
	}
	return expire
}

func (r *Redis) LocalIsAvailable() bool {
	return true
}

func (r *Redis) client() goRedis.UniversalClient {
	if r._client == nil {
		r._client = redis.Manager.GetByKey(r.redisKey)
	}
	return r._client
}

func New(keys ...string) cache.Driver {
	if len(keys) < 1 {
		keys = []string{""}
	}
	return &Redis{redisKey: keys[0]}
}
