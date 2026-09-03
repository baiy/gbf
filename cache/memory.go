package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Memory struct {
	client *ristretto.Cache[string, string]
	Base
}

func (d *Memory) Get(name string) (string, bool) {
	val, ok := d.client.Get(name)
	if !ok {
		return "", false
	}
	return val, true
}

func (d *Memory) Exist(name string) bool {
	_, ok := d.client.Get(name)
	return ok
}

func (d *Memory) Set(name string, value string, ttl time.Duration) {
	if ttl > 0 {
		d.client.SetWithTTL(name, value, 1, ttl)
	} else {
		d.client.Set(name, value, 1)
	}
	d.client.Wait()
}

func (d *Memory) Del(name string) {
	d.client.Del(name)
}

func (d *Memory) TTL(name string) time.Duration {
	if !d.Exist(name) {
		return time.Duration(-2)
	}
	val, _ := d.client.GetTTL(name)
	if val > 0 {
		return val
	}
	return time.Duration(-1)
}

func (d *Memory) LocalIsAvailable() bool {
	return false
}

func NewMemory() Driver {
	client, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,
		MaxCost:     10000 * 10000, // 最大缓存数量
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return &Memory{client: client}
}
