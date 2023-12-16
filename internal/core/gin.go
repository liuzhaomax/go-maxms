package core

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func init() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.ForceConsoleColor()
}

func InitGinEngine() *gin.Engine {
	gin.SetMode(GetConfig().Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	app.Use(LoggerToFile())
	return app
}
