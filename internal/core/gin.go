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
	RunMode            string `mapstructure:"run_mode"`
	MaxMultipartMemory int64  `mapstructure:"max_multipart_memory"`
}

// InitGinEngine Gin引擎的provider
func InitGinEngine() *gin.Engine {
	gin.SetMode(cfg.Lib.Gin.RunMode) // debug, test, release
	app := gin.Default()
	app.MaxMultipartMemory = cfg.Lib.Gin.MaxMultipartMemory << 20
	return app
}
