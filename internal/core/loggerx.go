package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"strings"
)

func LogSuccess(desc string) {
	cfg.App.Logger.WithField(SUCCESS, GetCallerName(3)).Debug(FormatCaller(true, GetCallerFileAndLine(3)))
	cfg.App.Logger.Info(FormatInfo(desc))
}

func LogFailure(code Code, desc string, err error) {
	cfg.App.Logger.WithField(FAILURE, GetCallerName(3)).Debug(FormatCaller(false, GetCallerFileAndLine(3)))
	cfg.App.Logger.Error(FormatError(code, desc, err))
}

var LoggerSet = wire.NewSet(wire.Struct(new(Logger), "*"), wire.Bind(new(ILogger), new(*Logger)))

// ILogger 主要用于中间件和handler
type ILogger interface {
	Succeed(string)
	Fail(Code, string, error)
	SucceedWithField(*gin.Context, string)
	FailWithField(*gin.Context, Code, string, error)
	SucceedWithFieldForRPC(context.Context, string)
	FailWithFieldForRPC(context.Context, Code, string, error)
}

const callerLevel = 4 // 调用堆栈第n层

type Logger struct {
	Logger *logrus.Logger
}

func (l *Logger) Succeed(desc string) {
	l.Logger.Debug(FormatCaller(true, GetCallerFileAndLine(callerLevel)))
	l.Logger.Info(FormatInfo(desc))
}

func (l *Logger) Fail(code Code, desc string, err error) {
	l.Logger.Info(FormatCaller(false, GetCallerFileAndLine(callerLevel)))
	l.Logger.Info(FormatError(code, desc, err))
}

func (l *Logger) SucceedWithField(c *gin.Context, desc string) {
	l.Logger.WithField(strings.ToLower(TraceId), c.Request.Header.Get(TraceId)).
		WithField(strings.ToLower(SpanId), c.Request.Header.Get(SpanId)).
		Debug(FormatCaller(true, GetCallerFileAndLine(callerLevel)))
	l.Logger.WithField(strings.ToLower(TraceId), c.Request.Header.Get(TraceId)).
		WithField(strings.ToLower(SpanId), c.Request.Header.Get(SpanId)).
		Info(FormatInfo(desc))
}

func (l *Logger) FailWithField(c *gin.Context, code Code, desc string, err error) {
	l.Logger.WithField(strings.ToLower(TraceId), c.Request.Header.Get(TraceId)).
		WithField(strings.ToLower(SpanId), c.Request.Header.Get(SpanId)).
		Debug(FormatCaller(false, GetCallerFileAndLine(callerLevel)))
	l.Logger.WithField(strings.ToLower(TraceId), c.Request.Header.Get(TraceId)).
		WithField(strings.ToLower(SpanId), c.Request.Header.Get(SpanId)).
		Info(FormatError(code, desc, err))
}

func (l *Logger) SucceedWithFieldForRPC(ctx context.Context, desc string) {
	md, _ := metadata.FromIncomingContext(ctx)
	l.Logger.WithField(strings.ToLower(TraceId), SelectFromMetadata(md, TraceId)).
		WithField(strings.ToLower(SpanId), SelectFromMetadata(md, SpanId)).
		Debug(FormatCaller(true, GetCallerFileAndLine(callerLevel)))
	l.Logger.WithField(strings.ToLower(TraceId), SelectFromMetadata(md, TraceId)).
		WithField(strings.ToLower(SpanId), SelectFromMetadata(md, SpanId)).
		Info(FormatInfo(desc))
}

func (l *Logger) FailWithFieldForRPC(ctx context.Context, code Code, desc string, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	l.Logger.WithField(strings.ToLower(TraceId), SelectFromMetadata(md, TraceId)).
		WithField(strings.ToLower(SpanId), SelectFromMetadata(md, SpanId)).
		Debug(FormatCaller(false, GetCallerFileAndLine(callerLevel)))
	l.Logger.WithField(strings.ToLower(TraceId), SelectFromMetadata(md, TraceId)).
		WithField(strings.ToLower(SpanId), SelectFromMetadata(md, SpanId)).
		Info(FormatError(code, desc, err))
}
