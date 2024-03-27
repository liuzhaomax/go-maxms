package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
)

func Register(root *gin.RouterGroup, handler *handler.HandlerUser, mw *middleware.Middleware) {
	root.GET("/login", handler.GetPuk)
	root.POST("/login", handler.PostLogin)

	root.Use(mw.Auth.ValidateToken())
	root.DELETE("/login", handler.DeleteLogin)
	routerUser := root.Group("/users")
	{
		routerUser.GET("/:userID", handler.GetUserByUserID)
	}
}
