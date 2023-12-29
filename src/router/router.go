package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/middleware"
	"github.com/liuzhaomax/go-maxms/src/api_user/handler"
)

func Register(app *gin.Engine, handler *handler.HandlerUser, mw *middleware.Middleware) {
	//app.Use(mw.Auth.VerifyToken())
	routerUser := app.Group("/users")
	{
		routerUser.GET("/:userID", handler.GetUserByUserID)
	}
}
