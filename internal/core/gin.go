package core

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func init() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.ForceConsoleColor()
}

// InitGinEngine Gin引擎的provider
func InitGinEngine() *gin.Engine {
	gin.SetMode(GetConfig().Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	app.Use(LoggerToFile()) // Gin使用logrus中间件处理日志
	return app
}
