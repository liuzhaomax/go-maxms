package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/api"
	"github.com/liuzhaomax/go-maxms/internal/core/pool"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	InjectorHTTP
	InjectorWS
	InjectorRPC
}

var InjectorHTTPSet = wire.NewSet(wire.Struct(new(InjectorHTTP), "*"))

type InjectorHTTP struct {
	Engine  *gin.Engine
	Handler *api.Handler
	DB      *gorm.DB
	Redis   *redis.Client
}

var InjectorWSSet = wire.NewSet(wire.Struct(new(InjectorWS), "*"))

type InjectorWS struct {
	Engine    *gin.Engine
	HandlerWs *api.HandlerWs
	DB        *gorm.DB
	Redis     *redis.Client
	Pool      *pool.Pool
}

var InjectorRPCSet = wire.NewSet(wire.Struct(new(InjectorRPC), "*"))

type InjectorRPC struct {
	HandlerRPC *api.HandlerRPC
	DB         *gorm.DB
	Redis      *redis.Client
}
