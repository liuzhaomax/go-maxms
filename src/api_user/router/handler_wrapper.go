package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/liuzhaomax/go-maxms/src/api_user/code"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
	"net/http"
)

type wrapperHandler = func(c *gin.Context) (any, error)

func wrapHandler(handler *handler.HandlerUser, handle wrapperHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggerFormat := config.GenGinLoggerFields(c)
		handler.Logger = handler.Logger.WithFields(loggerFormat)
		data, err := handle(c)
		if err != nil {
			var apiError *ext.ApiError
			errors.As(err, &apiError)
			statusCode := code.SelectStatusCode(apiError.Code)
			c.JSON(statusCode, gin.H{
				"status": gin.H{
					"code": apiError.Code,
					"desc": apiError.Message,
				},
				"data": data,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"code": 0,
				"desc": "success",
			},
			"data": data,
		})
	}
}

func wrapHandlerWS(handler *handler.HandlerUser, handle wrapperHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggerFormat := config.GenGinLoggerFields(c)
		handler.Logger = handler.Logger.WithFields(loggerFormat)

		data, err := handle(c)

		conn := c.MustGet(config.MyWsConn).(*websocket.Conn)
		defer handler.Close(conn)

		if err != nil {
			var apiError *ext.ApiError
			errors.As(err, &apiError)
			err = conn.WriteJSON(gin.H{
				"status": gin.H{
					"code": apiError.Code,
					"desc": apiError.Message,
				},
				"data": data,
			})
			if err != nil {
				return
			}
			return
		}
	}
}
