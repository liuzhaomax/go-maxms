package core

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	Config logger.Config
}

func InitGormLogger() *GormLogger {
	return &GormLogger{
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			ParameterizedQueries:      true,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return &GormLogger{}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	cfg.App.Logger.WithContext(ctx).Infof(msg, data...)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	cfg.App.Logger.WithContext(ctx).Errorf(msg, data...)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	cfg.App.Logger.WithContext(ctx).Errorf(msg, data...)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 语句和返回条数
	sql, rows := fc()
	// trace ID
	traceId := ctx.Value(TraceId)
	// 通用字段
	logFields := logrus.Fields{
		"sql":      sql,
		"start":    time.Now().Format(time.RFC3339Nano),
		"rows":     rows,
		"trace_id": traceId,
	}
	// Gorm 错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cfg.App.Logger.WithContext(ctx).WithFields(logFields).Info(FormatError(NotFound, "数据库 ErrRecordNotFound", err))
		} else {
			cfg.App.Logger.WithContext(ctx).WithFields(logFields).Error(FormatError(NotFound, "数据库 Error", err))
		}
	}
	// 慢查询日志
	if l.Config.SlowThreshold != 0 && elapsed > l.Config.SlowThreshold {
		cfg.App.Logger.WithContext(ctx).WithFields(logFields).Info(FormatInfo("数据库 Slow Log"))
	}
	// Debug模式下，且存在trace id，则记录所有 SQL 请求
	if cfg.Lib.DB.Debug && traceId != nil {
		cfg.App.Logger.WithContext(ctx).Debug(FormatCaller(true, GetCallerFileAndLine(5)))
		cfg.App.Logger.WithContext(ctx).WithFields(logFields).Info(FormatInfo("数据库 Query"))
	}
}
