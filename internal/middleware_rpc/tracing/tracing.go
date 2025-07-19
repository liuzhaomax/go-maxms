package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var TracingRPCSet = wire.NewSet(wire.Struct(new(TracingRPC), "*"))

type TracingRPC struct {
	Logger       *logrus.Logger
	TracerConfig *jConfig.Configuration
}

func (t *TracingRPC) Trace(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		t.AbortWithError(ctx, ext.ParseIssue, "元数据解析错误", err)

		return resp, err
	}

	if len(md[config.ParentId]) == 0 || len(md[config.TraceId]) == 0 ||
		len(md[config.RequestURI]) == 0 {
		t.AbortWithError(ctx, ext.ParseIssue, "元数据信息缺失", err)

		return resp, err
	}
	// 生成tracer
	tracer, closer, err := t.TracerConfig.NewTracer(jConfig.Logger(jaeger.StdLogger))
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(closer)

	if err != nil {
		t.AbortWithError(ctx, ext.MissingParameters, "tracer生成失败", err)

		return resp, err
	}
	// 创建span
	var span opentracing.Span

	parent := md[config.ParentId][0]
	if parent != "" {
		// 有父级span的时候提取父级span的traceID
		carrier := opentracing.TextMapCarrier{}
		for key, valSlice := range md {
			carrier.Set(key, valSlice[0])
		}
		// 注意inject会加入uber-trace-id到headers中，而md不会，没有这个不会产生ctx，
		// 测试单个服务需要手动加入，多个服务不会走这块代码
		ctxTracer, errTracer := tracer.Extract(opentracing.TextMap, carrier)
		if errTracer != nil {
			t.AbortWithError(ctx, ext.MissingParameters, "tracer提取失败", err)

			return resp, err
		}

		span = tracer.StartSpan(md[config.RequestURI][0], opentracing.ChildOf(ctxTracer))
	} else {
		span = tracer.StartSpan(md[config.RequestURI][0])
	}

	defer span.Finish()
	// 从 Span 上下文中获取 Trace ID 和 Span ID，并设置到header中
	spanContext := span.Context()
	traceID := spanContext.(jaeger.SpanContext).TraceID().String()
	spanID := spanContext.(jaeger.SpanContext).SpanID().String()

	if parent != "" {
		md[config.ParentId] = md[config.SpanId]
	}

	md[config.TraceId] = []string{traceID}
	md[config.SpanId] = []string{spanID}
	// 生成carrier
	carrier := opentracing.HTTPHeadersCarrier(md)

	err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		t.AbortWithError(ctx, ext.MissingParameters, "tracer注入失败", err)

		return resp, err
	}

	md[config.UberTraceId] = md["Uber-Trace-Id"]
	newCtx := metadata.NewIncomingContext(ctx, md) // 由服务内部不同拦截器读取不可以用Outgoing
	fmt.Println(md[config.UberTraceId][0])

	resp, err = handler(newCtx, req)

	return resp, err
}

func (t *TracingRPC) AbortWithError(ctx context.Context, args ...any) {
	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Code = ext.MissingParameters
		msg.Desc = "tracing错误"
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, code, desc, err)
		msg.Code = args[0].(ext.Code)
		msg.Desc = args[1].(string)
		msg.Err = args[2].(error)
	default:
		t.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	t.Logger.Error(formattedErr)
}
