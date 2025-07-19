package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
)

var TracingSet = wire.NewSet(wire.Struct(new(Tracing), "*"))

type Tracing struct {
	Logger       *logrus.Logger
	TracerConfig *jConfig.Configuration
}

// Trace 通过header传递traceID给下游服务来实现嵌套链路
func (t *Tracing) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成tracer
		tracer, closer, err := t.TracerConfig.NewTracer(jConfig.Logger(jaeger.NullLogger)) // 不打印log 没什么用
		defer func(closer io.Closer) {
			_ = closer.Close()
		}(closer)
		if err != nil {
			t.AbortWithError(c, err)
			return
		}
		// 创建span
		var span opentracing.Span
		parent := c.Request.Header.Get(core.ParentId)
		if parent != core.EmptyString {
			// 有父级span的时候提取父级span的traceID
			carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
			ctx, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
			if err != nil {
				t.AbortWithError(c, err)
				return
			}
			span = tracer.StartSpan(c.Request.URL.Path, opentracing.ChildOf(ctx))
		} else {
			span = tracer.StartSpan(c.Request.URL.Path)
		}
		defer span.Finish()
		// 从 Span 上下文中获取 Trace ID 和 Span ID，并设置到header中
		spanContext := span.Context()
		traceID := spanContext.(jaeger.SpanContext).TraceID().String()
		spanID := spanContext.(jaeger.SpanContext).SpanID().String()
		if parent != core.EmptyString {
			c.Request.Header.Set(core.ParentId, c.Request.Header.Get(core.SpanId))
		}
		c.Request.Header.Set(core.TraceId, traceID)
		c.Request.Header.Set(core.SpanId, spanID)
		// 生成carrier
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			t.AbortWithError(c, err)
			return
		}
		c.Next()
	}
}

func (t *Tracing) AbortWithError(c *gin.Context, args ...any) {
	msg := &core.MiddlewareMessage{
		StatusCode: 500,
		Code:       core.InternalServerError,
		Desc:       core.EmptyString,
		Err:        nil,
	}
	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusBadRequest
		msg.Code = core.MissingParameters
		msg.Desc = "tracing错误"
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(core.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		t.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := core.FormatError(msg.Code, msg.Desc, msg.Err)
	t.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, core.GenErrMsg(formattedErr))
}
