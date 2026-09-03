package worker

import "github.com/baiy/gbf/log"

type Logger struct {
	logger log.Logger
}

func (l Logger) Printf(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l Logger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l Logger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}
