//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/api"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/src/set"
)

func InitInjector() (*Injector, func(), error) {
	wire.Build(
		core.InitGinLogger,
		core.InitGinEngine,
		core.InitDB,
		api.APISet,
		set.HandlerSet,
		set.BusinessSet,
		set.ModelSet,
		core.ResponseSet,
		core.TransactionSet,
		middleware.MwsSet,
		middleware.MiddlewareSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
