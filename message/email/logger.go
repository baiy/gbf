package email

import (
	"log/slog"

	bootlog "github.com/baiy/gbf/log"
	"github.com/wneessen/go-mail/log"
)

type Logger struct {
	logger bootlog.Logger
}

func (l Logger) Debugf(msg log.Log) {
	l.logger.Logf(1, slog.LevelDebug, msg.Format, msg.Messages)
}

func (l Logger) Infof(msg log.Log) {
	l.logger.Logf(1, slog.LevelInfo, msg.Format, msg.Messages)

}

func (l Logger) Warnf(msg log.Log) {
	l.logger.Logf(1, slog.LevelWarn, msg.Format, msg.Messages)
}

func (l Logger) Errorf(msg log.Log) {
	l.logger.Logf(1, slog.LevelError, msg.Format, msg.Messages)
}
