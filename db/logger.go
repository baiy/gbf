package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baiy/gbf/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormLogger struct {
	logger log.Logger
}

func (l *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}
func (l *gormLogger) Info(_ context.Context, s string, args ...any) {
	l.logger.Infof(s, args...)
}
func (l *gormLogger) Warn(_ context.Context, s string, args ...any) {
	l.logger.Warnf(s, args...)
}
func (l *gormLogger) Error(_ context.Context, s string, args ...any) {
	l.logger.Errorf(s, args...)
}
func (l *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	args := []any{
		"rowsAffected", rowsAffected,
		"elapsed", fmt.Sprintf("%.3fs", elapsed.Seconds()),
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		args = append(args, "error", err.Error())
	}
	l.logger.Info(sql, args...)
}
