package handler

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user_rpc/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var HandlerUserSet = wire.NewSet(wire.Struct(new(HandlerUser), "*"))

type HandlerUser struct {
	Model    *model.ModelUser
	Tx       *core.Trans
	Redis    *redis.Client
	RocketMQ core.IRocketMQ
	Logger   *logrus.Logger
}
