package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
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

func (t *Tracing) Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer, closer, err := t.TracerConfig.NewTracer(jConfig.Logger(jaeger.StdLogger))
		defer func(closer io.Closer) {
			_ = closer.Close()
		}(closer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, t.GenErrMsg(c, "tracer生成失败", err))
		}
		span := tracer.StartSpan(c.Request.URL.Path)
		c.Set(core.Tracer, tracer)
		c.Set(core.Parent, span)
		c.Next()
		span.Finish()
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
