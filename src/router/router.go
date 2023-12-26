package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms-template/internal/middleware"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/handler"
)

func Register(app *gin.Engine, handler *handler.HandlerData, mw *middleware.Middleware) {
	//app.Use(mw.Auth.VerifyToken())
	routerData := app.Group("")
	{
		routerData.GET("/:id", handler.GetDataById)
	}
}
