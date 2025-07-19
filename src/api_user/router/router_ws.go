package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
)

func RegisterWs(root *gin.RouterGroup, handler *handler.HandlerUser, mw *middleware.Middleware) {
	root.Use(mw.Auth.ValidateToken())
	root.Use(mw.WsUpgrader.Upgrade())
	root.GET("/login", ext.WrapperRes(func(c *gin.Context) (any, error) {
		return handler.GetPuk(c)
	}))
}
