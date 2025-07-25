package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
)

func Register(root *gin.RouterGroup, handler *handler.HandlerUser, mw *middleware.Middleware) {
	root.GET("/login", wrapHandler(handler, func(c *gin.Context) (any, error) {
		return handler.GetPuk(c)
	}))
	root.POST("/login", wrapHandler(handler, func(c *gin.Context) (any, error) {
		return handler.PostLogin(c)
	}))

	root.Use(mw.Auth.ValidateToken())
	root.DELETE("/login", wrapHandler(handler, func(c *gin.Context) (any, error) {
		return handler.DeleteLogin(c)
	}))

	routerUser := root.Group("/users")
	{
		routerUser.GET("/:userID", wrapHandler(handler, func(c *gin.Context) (any, error) {
			return handler.GetUserByUserID(c)
		}))
	}
}
