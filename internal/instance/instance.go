package instance

import (
	"cmp"
	"fmt"
	"sync"
)

type Instance[T any] struct {
	Name   string
	values sync.Map
}

func New[T any](n string) *Instance[T] {
	return &Instance[T]{
		values: sync.Map{},
		Name:   n,
	}
}

func (i *Instance[T]) Get(key string) T {
	v, ok := i.values.Load(i.Key(key))
	if !ok || v == nil {
		panic(fmt.Sprintf("instance %s not found: %s", i.Name, key))
	}
	return v.(T)
}

func (i *Instance[T]) Add(k string, value func(k string) (T, error)) (T, error) {
	key := i.Key(k)
	if i.Has(key) {
		return i.Get(key), nil
	}
	v, err := value(key)
	if err != nil {
		return v, err
	}
	i.values.Store(key, v)
	return v, nil
}

func (i *Instance[T]) Has(key string) bool {
	_, ok := i.values.Load(i.Key(key))
	return ok
}

func (i *Instance[T]) Delete(key string) {
	i.values.Delete(i.Key(key))
}

func (i *Instance[T]) Key(key string) string {
	return cmp.Or(key, "default")
}

func (i *Instance[T]) Values() []T {
	var values []T
	i.values.Range(func(_, value any) bool {
		values = append(values, value.(T))
		return true
	})
	return values
}
