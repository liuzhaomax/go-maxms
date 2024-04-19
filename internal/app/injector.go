package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/api"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	InjectorHTTP
	InjectorRPC
}

var InjectorHTTPSet = wire.NewSet(wire.Struct(new(InjectorHTTP), "*"))

type InjectorHTTP struct {
	Engine  *gin.Engine
	Handler *api.Handler
	DB      *gorm.DB
	Redis   *redis.Client
}

var InjectorRPCSet = wire.NewSet(wire.Struct(new(InjectorRPC), "*"))

type InjectorRPC struct {
	HandlerRPC *api.HandlerRPC
	DB         *gorm.DB
	Redis      *redis.Client
}
