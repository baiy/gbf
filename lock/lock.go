package lock

import (
	"github.com/baiy/gbf/redis"
	"github.com/baiy/gbf/starter"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

// TryLock 尝试加锁执行
// 无法获取锁立即返回 相关逻辑不执行
func TryLock(name string, fn func()) bool {
	locker := newLocker(name)
	defer locker.Unlock()
	if !locker.TryLock() {
		logger().Warnf("try lock [%s] fail", name)
		return false
	}
	fn()
	return true
}

// Lock 加锁阻塞执行
// 未获取到锁一直等待 常用于常驻脚本多副本单例运行
func Lock(name string, fn func()) {
	locker := newLocker(name)
	defer locker.Unlock()
	locker.Lock()
	fn()
}

// Starter 创建启动器
func Starter(redisKey ...string) starter.Starter {
	return starter.Wrapper(func() error {
		if len(redisKey) == 0 {
			redisKey = []string{""}
		}
		rs = redsync.New(goredis.NewPool(redis.Manager.GetByKey(redisKey[0])))
		return nil
	})
}
