package env

import (
	"github.com/baiy/gbf/starter"
)

func Starter(v string) starter.Starter {
	return starter.Wrapper(func() error {
		Set(v)
		return nil
	})
}
