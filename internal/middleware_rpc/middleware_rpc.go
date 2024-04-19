package middleware_rpc

import (
	"context"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/middleware_rpc/auth"
	"github.com/liuzhaomax/go-maxms/internal/middleware_rpc/tracing"
	"github.com/liuzhaomax/go-maxms/internal/middleware_rpc/validator"
	"google.golang.org/grpc"
)

var MiddlewareRPCSet = wire.NewSet(wire.Struct(new(MiddlewareRPC), "*"))

type MiddlewareRPC struct {
	AuthRPC      *auth.AuthRPC
	ValidatorRPC *validator.ValidatorRPC
	TracingRPC   *tracing.TracingRPC
}

var MwsRPCSet = wire.NewSet(
	auth.AuthRPCSet,
	validator.ValidatorRPCSet,
	tracing.TracingRPCSet,
)

type IMiddlewareRPC interface {
	GenOkMsg(context.Context, string) string
	GenErrMsg(context.Context, string, error) error
}

var _ IMiddlewareRPC = (*auth.AuthRPC)(nil)
var _ IMiddlewareRPC = (*validator.ValidatorRPC)(nil)
var _ IMiddlewareRPC = (*tracing.TracingRPC)(nil)

// 连接多个 UnaryInterceptor
func ChainUnaryInterceptors(interceptorMap map[string][]grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 根据请求的接口信息获取对应的拦截器列表
		interceptors, ok := interceptorMap[info.FullMethod] // 例子"/StatsService/GetStatsArticleMain"在生成的pb中有具体名字，不加包名
		if !ok {
			return handler(ctx, req) // 如果没有找到对应接口的拦截器，直接执行原始处理函数
		}
		// 递归调用，将拦截器链连接起来
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = wrapUnaryInterceptor(interceptors[i], chain, info)
		}
		return chain(ctx, req)
	}
}

// 辅助函数，用于包装单个 UnaryInterceptor
func wrapUnaryInterceptor(interceptor grpc.UnaryServerInterceptor, handler grpc.UnaryHandler, info *grpc.UnaryServerInfo) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		})
	}
}
