package manager

import (
	"fmt"

	"github.com/baiy/gbf/internal/instance"
	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/starter"
	"github.com/baiy/gbf/util"
)

type Create[T any, O any] = func(mo Option[O]) (T, error)
type Start[T any] = func(T) error

type Manager[T any, O any] struct {
	instances *instance.Instance[*Instance[T, O]]

	// 创建实例方法
	create Create[T, O]

	// 实例启动方法
	start Start[T]

	// 默认配置
	defaultOptionGetter func() O
}

// Get 获取默认实例
func (m *Manager[T, O]) Get() T {
	return m.GetByKey("")
}

// GetByKey 获取指定实例
func (m *Manager[T, O]) GetByKey(key string) T {
	return m.instances.Get(key).Value
}

// New 创建实例
// 基于 key 的单例
func (m *Manager[T, O]) New(key string, options ...util.WithOption[O]) (T, error) {
	v, err := m.instances.Add(key, func(k string) (*Instance[T, O], error) {
		// 配置处理
		var option *O
		if m.defaultOptionGetter == nil {
			option = new(O)
		} else {
			option = new(m.defaultOptionGetter())
		}
		util.ApplyOption(option, options...)
		value, err := m.create(Option[O]{
			Key:    k,
			Value:  *option,
			Logger: log.With(fmt.Sprintf("%s.%s", m.instances.Name, k)),
		})
		if err != nil {
			return nil, err
		}
		return newInstance(k, value, option), nil
	})

	if err != nil {
		var zero T
		return zero, err
	}

	return v.Value, nil
}

// Starter 创建全局启动器 注入默认实例
func (m *Manager[T, O]) Starter(options ...util.WithOption[O]) starter.Starter {
	return m.StarterByKey("", options...)
}

// StarterByKey 创建全局启动器
func (m *Manager[T, O]) StarterByKey(key string, options ...util.WithOption[O]) starter.Starter {
	return starter.Wrapper(
		func() error {
			_, err := m.New(key, options...)
			return err
		},
		func() error {
			if m.start == nil {
				return nil
			}
			return m.start(m.GetByKey(key))
		},
	)
}

func New[T any, O any](
	name string,
	create Create[T, O],
	start Start[T],
	defaultOptionGetter func() O,
) *Manager[T, O] {
	return &Manager[T, O]{
		instances:           instance.New[*Instance[T, O]](name),
		create:              create,
		start:               start,
		defaultOptionGetter: defaultOptionGetter,
	}
}
