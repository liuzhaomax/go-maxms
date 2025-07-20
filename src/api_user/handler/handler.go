package handler

import (
	"context"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/liuzhaomax/go-maxms/internal/core/pool/ws"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

var HandlerUserSet = wire.NewSet(wire.Struct(new(HandlerUser), "*"))

type HandlerUser struct {
	Model    *model.ModelUser
	Logger   *logrus.Entry
	RocketMQ config.IRocketMQ
	Tx       *ext.Trans
	Redis    *redis.Client
	Pool     *ws.WsPool
}

// Close 错误处理统一在wrap中关闭连接
func (h *HandlerUser) Close(conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.Pool.CloseConn(ctx, conn)
	if err != nil {
		h.Logger.Error(ext.FormatError(ext.ConnectionFailed, "连接关闭失败", err))
		return
	}
	h.Logger.Info(ext.FormatInfo("连接关闭成功"))
}
