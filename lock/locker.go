package lock

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/baiy/gbf/closer"
	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/worker"
	"github.com/go-redsync/redsync/v4"
)

var rs *redsync.Redsync

type Locker struct {
	name string
	mu   *redsync.Mutex

	is bool

	extendCtx    context.Context
	extendCancel context.CancelFunc
}

func (l *Locker) Lock() {
	l.is = false
	if err := l.mu.LockContext(closer.CloseCtx); err == nil {
		logger().Info("get lock")
		// 锁获取成功 启动自动续期
		l.autoExtend()
		l.is = true
		return
	}
}

func (l *Locker) TryLock() bool {
	l.is = false
	if l.mu.TryLockContext(closer.CloseCtx) == nil {
		logger().Info("get lock")
		// 锁获取成功 启动自动续期
		l.autoExtend()
		l.is = true
		return true
	}
	return false
}

func (l *Locker) Unlock() {
	logger().Info("release lock")
	if l.extendCancel != nil {
		l.extendCancel()
	}
	if l.is {
		_, err := l.mu.Unlock()
		if err != nil {
			logger().Errorf("redsync unlock %s error:%s", l.name, err)
		}
	}
	l.is = false
}

func (l *Locker) autoExtend() {
	l.extendCtx, l.extendCancel = context.WithCancel(context.Background())
	worker.Async(func() {
		// 每隔4s续期进行续期
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-l.extendCtx.Done():
				return
			case <-ticker.C:
				if _, err := l.mu.Extend(); err != nil {
					logger().Errorf("redsync extend  %s error %s", l.name, err)
					return
				}
			}
		}
	})
}

func newLocker(name string) *Locker {
	if rs == nil {
		panic("redsync not init")
	}
	if name == "" {
		panic("name is empty")
	}
	name = fmt.Sprintf("%s:%s:%s", constant.BaseName, "lock", name)
	return &Locker{
		name: name,
		mu: rs.NewMutex(
			name,
			redsync.WithExpiry(8*time.Second),
			redsync.WithTries(math.MaxInt), // 设置一个超级大的值 用于模拟阻塞锁
		),
	}
}
