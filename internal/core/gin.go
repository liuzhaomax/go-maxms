package core

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func init() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.ForceConsoleColor()
}

// InitGinEngine Gin引擎的provider
func InitGinEngine() *gin.Engine {
	gin.SetMode(GetConfig().Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	logHandlerFunc, _ := LoggerToFile() // Gin使用logrus日志处理中间件
	app.Use(logHandlerFunc)
	return app
}

// InitGinLogger 日志对象的provider，从而使用额外的自定义日志，且可以自定义输出不同级别的日志
func InitGinLogger() *logrus.Logger {
	_, logger := LoggerToFile()
	return logger
}
