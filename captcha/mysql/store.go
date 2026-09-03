package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/baiy/gbf/captcha"
	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/db"
	"github.com/baiy/gbf/db/mysql"
)

type Captcha struct {
	db.Primary
	db.Created
	Key   string
	Value string
}

type Store struct {
	captcha.Base
	mysqlKey string
	_client  *db.Db
}

func (s *Store) Set(id string, value string) error {
	return s.client().Insert[Captcha](context.Background(), &Captcha{
		Key:   s.id(id),
		Value: value,
	})
}

func (s *Store) Get(id string, clear bool) string {
	defer func() {
		if clear {
			s.client().Model[Captcha]().Where("key = ?", s.id(id)).Delete()
		}
	}()
	v := s.client().Model[Captcha]().Where("key = ? and created_at > ?", s.id(id), time.Now().Add(-s.Ttl).Format(time.DateTime)).Get()
	if v == nil {
		return ""
	}
	return v.Value
}

func (s *Store) Verify(id, answer string, clear bool) bool {
	return answer != "" && s.Get(id, clear) == answer
}

func (s *Store) id(id string) string {
	return fmt.Sprintf("%s:%s:%s", constant.AppName, s.Key, id)
}

func (s *Store) client() *db.Db {
	if s._client == nil {
		s._client = mysql.Manager.GetByKey(s.mysqlKey)
	}
	return s._client
}

func New(keys ...string) captcha.Store {
	if len(keys) < 1 {
		keys = []string{""}
	}
	return &Store{mysqlKey: keys[0]}
}
