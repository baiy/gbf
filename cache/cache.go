package cache

import (
	"fmt"
	"slices"
	"time"
	"uuid"

	"github.com/baiy/gbf/log"
	"github.com/baiy/gbf/util/json"
)

type Cache struct {
	prefix string

	// 子缓存前缀
	subPrefix string

	driver Driver
	logger log.Logger

	localStore *Local
	localKeys  []string
}

// Get 获取缓存
func (c *Cache) Get[T any](name string) (T, bool) {
	// 缓存key
	key := c.Key(name)
	c.logger.Infof("cache get %s", key)
	// 二级缓存key
	localKey := c.localKey(name)
	writeLocalVersion := ""
	// 优化获取本地二级缓存
	if localKey != "" {
		version, exist := c.driver.Get(localKey)
		if !exist {
			// 不存在版本号，删除本地缓存
			c.localStore.Del(key)
		} else {
			v, ok := c.localStore.Get[T](key, version)
			if ok {
				c.logger.Infof("local cache hit")
				// 本地缓存可用 返回
				return v, true
			}
			// 本地缓存不可用 标记写入版本号
			writeLocalVersion = version
		}
	}

	value, exist := c.driver.Get(key)
	if !exist {
		var zero T
		return zero, false
	}
	v := json.Unserialize[T](value)

	// 写入本地二级缓存
	if writeLocalVersion != "" {
		c.logger.Infof("local cache get set %s=>%s", key, writeLocalVersion)
		c.localStore.Set(key, v, writeLocalVersion)
	}
	return v, true
}

func (c *Cache) GetDel[T any](name string) (T, bool) {
	defer func() {
		c.Del(name)
	}()
	return c.Get[T](name)
}

// Set 设置缓存
func (c *Cache) Set(name string, value any, _ttl ...time.Duration) {
	key := c.Key(name)
	c.logger.Infof("cache set %s", key)
	var ttl time.Duration
	if len(_ttl) > 0 {
		ttl = _ttl[0]
	}
	// 二级缓存key
	localKey := c.localKey(name)
	if localKey != "" {
		version := uuid.NewV4().String()
		c.logger.Infof("local cache set %s=>%s", key, version)
		c.driver.Set(localKey, version, ttl)
	}
	c.driver.Set(key, json.Serialize(value), ttl)

}

// Del 删除缓存
func (c *Cache) Del(name string) {
	key := c.Key(name)
	c.logger.Infof("cache del %s", key)
	c.driver.Del(key)
	// 二级缓存key
	localKey := c.localKey(name)
	if localKey != "" {
		c.driver.Del(localKey)
	}
}

// Exist 缓存是否存在
func (c *Cache) Exist(name string) bool {
	return c.driver.Exist(c.Key(name))
}

// TTL 获取过期时间
func (c *Cache) TTL(name string) time.Duration {
	return c.driver.TTL(c.Key(name))
}

// Key 获取缓存key
func (c *Cache) Key(name string) string {
	return fmt.Sprintf("%s:%s%s", c.prefix, c.subPrefix, name)
}

// localKey 获取本地二级缓存版本key
// 不启用本地二级缓存 返回 空
func (c *Cache) localKey(name string) string {
	if c.localStore == nil || len(c.localKeys) == 0 || !slices.Contains(c.localKeys, fmt.Sprintf("%s%s", c.subPrefix, name)) {
		return ""
	}
	return c.Key(fmt.Sprintf("%s:%s", name, "version"))
}

func (c *Cache) With(subPrefix string) *Cache {
	logger := c.logger.With(subPrefix)
	return NewCache(
		c.prefix,
		fmt.Sprintf("%s%s", c.subPrefix, subPrefix),
		c.driver,
		logger,
		c.localStore,
		c.localKeys,
	)
}

func NewCache(prefix, subPrefix string, driver Driver, logger log.Logger, local *Local, localKey []string) *Cache {
	return &Cache{
		driver: driver,
		logger: logger,

		prefix:    prefix,
		subPrefix: subPrefix,

		localStore: local,
		localKeys:  localKey,
	}
}
