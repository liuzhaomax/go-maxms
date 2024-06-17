package api

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/middleware_rpc"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/business"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/pb"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net/http"
)

var APIRPCSet = wire.NewSet(wire.Struct(new(HandlerRPC), "*"), wire.Bind(new(APIRPC), new(*HandlerRPC)))

type APIRPC interface {
	Register() *grpc.Server
}

type HandlerRPC struct {
	PrometheusRegistry *prometheus.Registry
	MiddlewareRPC      *middleware_rpc.MiddlewareRPC
	BusinessRPC        *business.BusinessUser
}

func (h *HandlerRPC) Register() *grpc.Server {
	interceptorsBasicChoice := []grpc.UnaryServerInterceptor{
		h.MiddlewareRPC.TracingRPC.Trace,
		core.LoggerForRPC,
		h.MiddlewareRPC.ValidatorRPC.ValidateMetadata,
		h.MiddlewareRPC.AuthRPC.ValidateSignature,
	}
	interceptorMap := map[string][]grpc.UnaryServerInterceptor{
		"/StatsService/GetStatsArticleMain": interceptorsBasicChoice,
	}

	// 连接多个中间件
	serverOpts := []grpc.ServerOption{}
	serverOpts = append(serverOpts, grpc.UnaryInterceptor(middleware_rpc.ChainUnaryInterceptors(interceptorMap)))
	// 创建gRPC服务
	server := grpc.NewServer(serverOpts...)
	// 注册接口
	pb.RegisterUserServiceServer(server, h.BusinessRPC)

	// 健康检查
	healthCheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthCheck)

	// prometheus TODO 默认是9090提供http服务来提供metrics，需要使用gin来提供此服务，同时注册到consul，
	// TODO 也就是一个rpc会启动两个服务，一个提供rpc接口，一个提供监控http接口，但这样需要再开一个docker端口映射，后面再说吧
	// grpc_prometheus.Register(server)
	// http.Handle("/metrics", http.HandlerFunc(h.MetricsHandler))

	return server
}

func (h *HandlerRPC) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(h.PrometheusRegistry, promhttp.HandlerOpts{
		Registry: h.PrometheusRegistry,
	}).ServeHTTP(w, r)
}
