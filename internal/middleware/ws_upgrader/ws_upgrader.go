package ws_upgrader

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/sirupsen/logrus"
	"net/http"
)

var WsUpgraderSet = wire.NewSet(wire.Struct(new(WsUpgrader), "*"))

type WsUpgrader struct {
	Logger   *logrus.Logger
	Upgrader *websocket.Upgrader
}

func (wsUpgrader *WsUpgrader) Upgrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			wsUpgrader.AbortWithError(c, http.StatusInternalServerError, core.ProtocolUpgradeFailed, "http升级ws未能生成连接", err)
			return
		}
		c.Set(core.WsConn, conn)
		c.Next()
	}
}

func (wsUpgrader *WsUpgrader) AbortWithError(c *gin.Context, args ...any) {
	msg := &core.MiddlewareMessage{
		StatusCode: 500,
		Code:       core.InternalServerError,
		Desc:       core.EmptyString,
		Err:        nil,
	}
	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(core.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		wsUpgrader.Logger.Error("invalid arguments for AbortWithError")
	}
	formattedErr := core.FormatError(msg.Code, msg.Desc, msg.Err)
	wsUpgrader.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, core.GenErrMsg(formattedErr))
}
