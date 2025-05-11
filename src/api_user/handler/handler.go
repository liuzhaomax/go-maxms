package handler

import (
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var HandlerUserSet = wire.NewSet(wire.Struct(new(HandlerUser), "*"))

type HandlerUser struct {
	Model    *model.ModelUser
	Logger   *logrus.Logger
	RocketMQ core.IRocketMQ
	Tx       *core.Trans
	Redis    *redis.Client
}
