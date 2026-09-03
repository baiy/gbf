package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

type Local struct {
	client *ristretto.Cache[string, LocalItem]
}

type LocalItem struct {
	version string
	value   any
}

func (l *Local) Get[T any](key, version string) (T, bool) {
	item, ok := l.client.Get(key)
	if !ok || version == "" || item.version != version {
		var zero T
		return zero, false
	}
	return item.value.(T), true
}

func (l *Local) Set(key string, value any, version string) {
	l.client.SetWithTTL(key, LocalItem{
		version: version,
		value:   value,
	}, 1, time.Hour*24)
	l.client.Wait()
}

func (l *Local) Del(key string) {
	l.client.Del(key)
}

func newLocal() *Local {
	client, err := ristretto.NewCache(&ristretto.Config[string, LocalItem]{
		NumCounters: 1e7,
		MaxCost:     10000 * 10000,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return &Local{client: client}
}
