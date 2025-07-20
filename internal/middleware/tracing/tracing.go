package tracing

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
)

var TracingSet = wire.NewSet(wire.Struct(new(Tracing), "*"))

type Tracing struct {
	Logger       *logrus.Entry
	TracerConfig *jConfig.Configuration
}

// Trace 通过header传递traceID给下游服务来实现嵌套链路
func (t *Tracing) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成tracer
		tracer, closer, err := t.TracerConfig.NewTracer(
			jConfig.Logger(jaeger.NullLogger),
		) // 不打印log 没什么用
		defer func(closer io.Closer) {
			_ = closer.Close()
		}(closer)

		if err != nil {
			t.AbortWithError(c, err)

			return
		}
		// 创建span
		var span opentracing.Span

		parent := c.Request.Header.Get(config.ParentId)
		if parent != "" {
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

		if parent != "" {
			c.Request.Header.Set(config.ParentId, c.Request.Header.Get(config.SpanId))
		}

		c.Request.Header.Set(config.TraceId, traceID)
		c.Request.Header.Set(config.SpanId, spanID)
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
	t.Logger = t.Logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"client_ip":  config.GetClientIP(c),
		"user_agent": config.GetUserAgent(c),
		"token":      c.GetHeader(config.Authorization),
		"trace_id":   c.GetHeader(config.TraceId),
		"span_id":    c.GetHeader(config.SpanId),
		"parent_id":  c.GetHeader(config.ParentId),
		"app_id":     c.GetHeader(config.AppId),
		"request_id": c.GetHeader(config.RequestId),
		"user_id":    c.GetHeader(config.UserId),
	})
	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusBadRequest
		msg.Code = ext.MissingParameters
		msg.Desc = "tracing错误"
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(ext.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		t.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	t.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, ext.GenErrMsg(formattedErr))
}
