package mysql

import (
	"context"
	"time"

	"github.com/baiy/gbf/cache"
	"github.com/baiy/gbf/db"
	"github.com/baiy/gbf/db/mysql"
	"github.com/baiy/gbf/db/mysql/field"
	"github.com/baiy/gbf/errors"
)

type Cache struct {
	db.Base
	Name   string
	Value  string
	Expire *field.DateTime
}

type Mysql struct {
	mysqlKey string
	_client  *db.Db
	cache.Base
}

func (m *Mysql) Get(name string) (string, bool) {
	v := m.client().Model[Cache]().
		Where("name = ? and (expire >= ? or expire is null)", name, time.Now().Format(time.DateTime)).
		Get()
	if v == nil {
		return "", false
	}
	return v.Value, true
}

func (m *Mysql) Exist(name string) bool {
	return m.client().Model[Cache]().
		Where("name = ? and (expire >= ? or expire is null)", name, time.Now().Format(time.DateTime)).
		Exist()
}

func (m *Mysql) Set(name string, value string, ttl time.Duration) {
	defer func() {
		err := errors.PanicToError(func() error {
			m.clear()
			return nil
		})
		if err != nil {
			m.Logger.Errorf("clear cache error: %v", err)
			return
		}
	}()
	err := m.client().Transaction(func(tx *db.Db) error {
		v := m.client().Model[Cache]().Omit("value").
			Where("name = ?", name).
			Get()
		if v == nil {
			v = &Cache{
				Name: name,
			}
		}
		v.Value = value
		v.Expire = nil
		if ttl > 0 {
			v.Expire = new(field.NewDateTime(time.Now().Add(ttl)))
		}
		return m.client().Save(context.Background(), v)
	})
	if err != nil {
		panic(err)
	}
}

func (m *Mysql) Del(name string) {
	m.client().Model[Cache]().Where("name = ?", name).Delete()
}

func (m *Mysql) TTL(name string) time.Duration {
	v := m.client().Model[Cache]().Select("id,expire").
		Where("name = ? and (expire >= ? or expire is null)", name, time.Now().Format(time.DateTime)).
		Get()
	if v == nil {
		return time.Duration(-2)
	}
	if v.Expire == nil {
		return time.Duration(-1)
	}
	return v.Expire.Time().Sub(time.Now())
}

func (m *Mysql) LocalIsAvailable() bool {
	return true
}

func (m *Mysql) clear() {
	m.Logger.Info("clearing cache")
	// 删除30分钟之前的过期数据
	m.client().Model[Cache]().Where("expire <= ?", time.Now().Add(-30*time.Minute).Format(time.DateTime)).Delete()
}

func (m *Mysql) client() *db.Db {
	if m._client == nil {
		m._client = mysql.Manager.GetByKey(m.mysqlKey)
	}
	return m._client
}

func New(keys ...string) cache.Driver {
	if len(keys) < 1 {
		keys = []string{""}
	}
	return &Mysql{mysqlKey: keys[0]}
}
