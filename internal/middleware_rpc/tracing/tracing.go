package tracing

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

var TracingRPCSet = wire.NewSet(wire.Struct(new(TracingRPC), "*"))

type TracingRPC struct {
	Logger       core.ILogger
	TracerConfig *jConfig.Configuration
}

func (t *TracingRPC) Trace(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = t.GenErrMsg(ctx, "元数据解析错误", err)
		return
	}
	if len(md[core.ParentId]) == 0 || len(md[core.TraceId]) == 0 || len(md[core.RequestURI]) == 0 {
		err = t.GenErrMsg(ctx, "元数据信息缺失", err)
		return
	}
	// 生成tracer
	tracer, closer, err := t.TracerConfig.NewTracer(jConfig.Logger(jaeger.StdLogger))
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(closer)
	if err != nil {
		err = t.GenErrMsg(ctx, "tracer生成失败", err)
		return
	}
	// 创建span
	var span opentracing.Span
	parent := md[core.ParentId][0]
	if parent != core.EmptyString {
		// 有父级span的时候提取父级span的traceID
		carrier := opentracing.TextMapCarrier{}
		for key, valSlice := range md {
			carrier.Set(key, valSlice[0])
		}
		// 注意inject会加入uber-trace-id到headers中，而md不会，没有这个不会产生ctx，
		// 测试单个服务需要手动加入，多个服务不会走这块代码
		ctxTracer, errTracer := tracer.Extract(opentracing.TextMap, carrier)
		if errTracer != nil {
			err = t.GenErrMsg(ctx, "tracer提取失败", errTracer)
			return
		}
		span = tracer.StartSpan(md[core.RequestURI][0], opentracing.ChildOf(ctxTracer))
	} else {
		span = tracer.StartSpan(md[core.RequestURI][0])
	}
	defer span.Finish()
	// 从 Span 上下文中获取 Trace ID 和 Span ID，并设置到header中
	spanContext := span.Context()
	traceID := spanContext.(jaeger.SpanContext).TraceID().String()
	spanID := spanContext.(jaeger.SpanContext).SpanID().String()
	if parent != core.EmptyString {
		md[core.ParentId] = md[core.SpanId]
	}
	md[core.TraceId] = []string{traceID}
	md[core.SpanId] = []string{spanID}
	// 生成carrier
	carrier := opentracing.HTTPHeadersCarrier(md)
	err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		err = t.GenErrMsg(ctx, "tracer注入失败", err)
		return
	}
	md[core.UberTraceId] = md["Uber-Trace-Id"]
	newCtx := metadata.NewIncomingContext(ctx, md) // 由服务内部不同拦截器读取不可以用Outgoing
	fmt.Println(md[core.UberTraceId][0])
	resp, err = handler(newCtx, req)
	return
}

func (t *TracingRPC) GenOkMsg(ctx context.Context, desc string) string {
	t.Logger.SucceedWithFieldForRPC(ctx, desc)
	return core.FormatInfo(desc)
}

func (t *TracingRPC) GenErrMsg(ctx context.Context, desc string, err error) error {
	t.Logger.FailWithFieldForRPC(ctx, core.Unknown, desc, err)
	return core.FormatError(core.Unknown, desc, err)
}
