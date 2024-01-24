package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/api"
	businessRpc "github.com/liuzhaomax/go-maxms/src/api_user_rpc/business"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	RPCService *businessRpc.BusinessUser
	Engine     *gin.Engine
	Handler    *api.Handler
	DB         *gorm.DB
	Redis      *redis.Client
}
