package cache

import (
	"fmt"

	"github.com/baiy/gbf/cache"
	"github.com/baiy/gbf/captcha"
)

type Store struct {
	captcha.Base
	cacheKey string
	_client  *cache.Cache
}

func (s *Store) Set(id string, value string) error {
	s.client().Set(s.id(id), value, s.Ttl)
	return nil
}

func (s *Store) Get(id string, clear bool) string {
	defer func() {
		if clear {
			s.client().Del(s.id(id))
		}
	}()
	v, ok := s.client().Get[string](s.id(id))
	if !ok {
		return ""
	}
	return v
}

func (s *Store) Verify(id, answer string, clear bool) bool {
	return answer != "" && s.Get(id, clear) == answer
}

func (s *Store) id(id string) string {
	return fmt.Sprintf("captcha:%s:%s", s.Key, id)
}

func (s *Store) client() *cache.Cache {
	if s._client == nil {
		s._client = cache.Manager.GetByKey(s.cacheKey)
	}
	return s._client
}

func New(keys ...string) captcha.Store {
	if len(keys) < 1 {
		keys = []string{""}
	}
	return &Store{cacheKey: keys[0]}
}
