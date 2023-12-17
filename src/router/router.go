package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms-template-me/src/dataAPI/handler"
)

func Register(handler *handler.HData, app *gin.Engine) {
	//itcpt := &interceptor.Interceptor{}
	routerData := app.Group("")
	{
		routerData.GET("/:id", handler.GetDataById)
	}
}
