package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
	"io"
	"net/http"
)

var TracingSet = wire.NewSet(wire.Struct(new(Tracing), "*"))

type Tracing struct {
	Logger       core.ILogger
	TracerConfig *jConfig.Configuration
}

// Trace 通过header传递traceID给下游服务来实现嵌套链路
func (t *Tracing) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成tracer
		tracer, closer, err := t.TracerConfig.NewTracer(jConfig.Logger(jaeger.StdLogger))
		defer func(closer io.Closer) {
			_ = closer.Close()
		}(closer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, t.GenErrMsg(c, "tracer生成失败", err))
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
				c.AbortWithStatusJSON(http.StatusUnauthorized, t.GenErrMsg(c, "tracer提取失败", err))
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
		} else {
			c.Request.Header.Set(core.ParentId, spanID)
		}
		c.Request.Header.Set(core.TraceId, traceID)
		c.Request.Header.Set(core.SpanId, spanID)
		// 生成carrier
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, t.GenErrMsg(c, "tracer注入失败", err))
			return
		}
		c.Next()
	}
}

func (t *Tracing) GenOkMsg(c *gin.Context, desc string) string {
	t.Logger.SucceedWithField(c, desc)
	return core.FormatInfo(desc)
}

func (t *Tracing) GenErrMsg(c *gin.Context, desc string, err error) error {
	t.Logger.FailWithField(c, core.Unknown, desc, err)
	return core.FormatError(core.Unknown, desc, err)
}
