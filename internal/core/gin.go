package core

import (
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
)

func init() {
	gin.DefaultWriter = colorable.NewColorableStdout()
	gin.ForceConsoleColor()
}

type Gin struct {
	RunMode string `mapstructure:"run_mode"`
}

// InitGinEngine Gin引擎的provider
func InitGinEngine() *gin.Engine {
	gin.SetMode(GetConfig().Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	return app
}
