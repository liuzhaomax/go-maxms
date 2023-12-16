//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/api"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
	"github.com/liuzhaomax/go-maxms-template-me/src/handler"
)

func InitInjector() (*Injector, error) {
	wire.Build(
		core.InitGinEngine,
		api.APISet,
		handler.HandlerSet,
		InjectorSet,
	)
	return new(Injector), nil
}
