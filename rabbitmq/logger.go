package rabbitmq

import (
	"fmt"
	"log/slog"

	"github.com/baiy/gbf/log"
)

type Logger struct {
	Logger log.Logger
}

func (l Logger) log(level slog.Level, s string, i ...any) {
	l.Logger.Log(2, level, fmt.Sprintf(s, i...))
}

func (l Logger) Fatalf(s string, i ...any) {
	l.log(slog.LevelError, s, i...)
}

func (l Logger) Errorf(s string, i ...any) {
	l.log(slog.LevelError, s, i...)
}

func (l Logger) Warnf(s string, i ...any) {
	l.log(slog.LevelWarn, s, i...)
}

func (l Logger) Infof(s string, i ...any) {
	l.log(slog.LevelInfo, s, i...)
}

func (l Logger) Debugf(s string, i ...any) {
	l.log(slog.LevelDebug, s, i...)
}

func (l Logger) With(key string) Logger {
	return Logger{
		Logger: l.Logger.With(key),
	}
}
