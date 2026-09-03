package cache

import (
	"fmt"

	"github.com/baiy/gbf/constant"
	"github.com/baiy/gbf/manager"
)

func create(mo manager.Option[Option]) (*Cache, error) {
	var local *Local
	if mo.Value.driver.LocalIsAvailable() {
		local = newLocal()
	}
	mo.Value.driver.SetLogger(mo.Logger)
	return NewCache(
		fmt.Sprintf("%s:cache:%s:%s", constant.AppName, mo.Key, mo.Value.prefix),
		"",
		mo.Value.driver,
		mo.Logger,
		local,
		mo.Value.localKeys,
	), nil
}

var Manager = manager.New[*Cache, Option](
	"cache",
	create,
	nil,
	func() Option {
		return Option{
			prefix: "cache",
			driver: NewMemory(),
		}
	},
)
