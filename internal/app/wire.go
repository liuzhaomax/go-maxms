//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/api"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/liuzhaomax/go-maxms/internal/core/pool"
	"github.com/liuzhaomax/go-maxms/internal/core/pool/ws"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/internal/middleware_rpc"
	"github.com/liuzhaomax/go-maxms/src/set"
)

func InitInjector() (*Injector, func(), error) {
	wire.Build(
		config.InitLogrus,
		config.InitGinEngine,
		config.InitDB,
		config.InitRedis,
		config.InitWebSocket,
		config.InitTracer,
		config.InitPrometheusRegistry,
		pool.InitPool,
		ws.InitWsPool,
		api.APISet,
		api.APIWSSet,
		api.APIRPCSet,
		set.HandlerSet,
		set.ModelSet,
		ext.TransactionSet,
		config.RocketMQSet,
		middleware.MwsSet,
		middleware.MiddlewareSet,
		middleware_rpc.MwsRPCSet,
		middleware_rpc.MiddlewareRPCSet,
		InjectorHTTPSet,
		InjectorWSSet,
		InjectorRPCSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
