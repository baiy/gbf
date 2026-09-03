package redis

import (
	"context"
	"fmt"

	"github.com/baiy/gbf/captcha"
	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/redis"
	goRedis "github.com/redis/go-redis/v9"
)

type Store struct {
	captcha.Base
	redisKey string
	_client  goRedis.UniversalClient
}

func (s *Store) Set(id string, value string) error {
	return s.client().Set(context.Background(), s.id(id), value, s.Ttl).Err()
}

func (s *Store) Get(id string, clear bool) string {
	defer func() {
		if clear {
			s.client().Del(context.Background(), s.id(id))
		}
	}()
	v, err := s.client().Get(context.Background(), s.id(id)).Result()
	if err != nil {
		s.Logger.Errorf("redis get [%s] error: %s", s.id(id), err)
	}

	return v
}

func (s *Store) Verify(id, answer string, clear bool) bool {
	return answer != "" && s.Get(id, clear) == answer
}

func (s *Store) id(id string) string {
	return fmt.Sprintf("%s:captcha:%s:%s", constant.AppName, s.Key, id)
}

func (s *Store) client() goRedis.UniversalClient {
	if s._client == nil {
		s._client = redis.Manager.GetByKey(s.redisKey)
	}
	return s._client
}

func New(keys ...string) captcha.Store {
	if len(keys) < 1 {
		keys = []string{""}
	}
	return &Store{redisKey: keys[0]}
}
