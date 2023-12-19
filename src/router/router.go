package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms-template-me/internal/middleware"
	"github.com/liuzhaomax/go-maxms-template-me/src/dataAPI/handler"
)

func Register(app *gin.Engine, handler *handler.HandlerData, mw *middleware.Middleware) {
	//app.Use(mw.Auth.VerifyToken())
	routerData := app.Group("")
	{
		routerData.GET("/:id", handler.GetDataById)
	}
}
