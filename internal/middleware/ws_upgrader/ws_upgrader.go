package ws_upgrader

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/sirupsen/logrus"
)

var WsUpgraderSet = wire.NewSet(wire.Struct(new(WsUpgrader), "*"))

type WsUpgrader struct {
	Logger   *logrus.Entry
	Upgrader *websocket.Upgrader
}

func (wsUpgrader *WsUpgrader) Upgrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := wsUpgrader.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			wsUpgrader.AbortWithError(
				c,
				http.StatusInternalServerError,
				ext.ProtocolUpgradeFailed,
				"http升级ws未能生成连接",
				err,
			)
			return
		}

		// 设置心跳处理 传输层通信控制
		conn.SetPingHandler(func(appData string) error {
			return conn.WriteControl(
				websocket.PongMessage,
				[]byte(appData),
				time.Now().Add(time.Second),
			)
		})

		// 设置读取超时（心跳超时检测）
		err = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		if err != nil {
			wsUpgrader.AbortWithError(c, err)
			return
		}
		conn.SetPongHandler(func(string) error {
			err = conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
			if err != nil {
				wsUpgrader.AbortWithError(c, err)
				return err
			}
			return nil
		})

		c.Set(config.MyWsConn, conn)
		c.Next()
	}
}

func (wsUpgrader *WsUpgrader) AbortWithError(c *gin.Context, args ...any) {
	loggerFormat := config.GenGinLoggerFields(c)
	wsUpgrader.Logger = wsUpgrader.Logger.WithFields(loggerFormat)

	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(ext.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		wsUpgrader.Logger.Error("invalid arguments for AbortWithError")
	}

	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	wsUpgrader.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, ext.GenErrMsg(formattedErr))
}
