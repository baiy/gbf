package lock

import (
	"sync"

	"github.com/baiy/gbf/log"
)

var loggerOne sync.Once
var _logger log.Logger

func logger() log.Logger {
	loggerOne.Do(func() {
		_logger = log.With("lock")
	})
	return _logger
}
