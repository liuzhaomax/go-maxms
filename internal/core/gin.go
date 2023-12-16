package core

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms-template-me/internal/api"
	"github.com/mattn/go-colorable"
)

func init() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.ForceConsoleColor()
}

func InitGinEngine(api api.API) *gin.Engine {
	gin.SetMode(GetConfig().Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	app.Use(LoggerToFile())
	api.Register(app)
	return app
}
