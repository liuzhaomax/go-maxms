//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template-me/internal/api"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
	"github.com/liuzhaomax/go-maxms-template-me/src/dataAPI/handler"
)

func InitInjector() (*Injector, error) {
	wire.Build(
		core.InitGinLogger,
		core.InitGinEngine,
		api.APISet,
		handler.HandlerSet,
		core.ResponseSet,
		InjectorSet,
	)
	return new(Injector), nil
}
