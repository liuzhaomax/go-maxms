package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ValidatorSet = wire.NewSet(wire.Struct(new(Validator), "*"))

type Validator struct {
	Logger *logrus.Entry
	Redis  *redis.Client
}

func (v *Validator) ValidateHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := config.ValidateHeaders(c)
		if err != nil {
			v.AbortWithError(c, err)

			return
		}

		c.Next()
	}
}

func (v *Validator) AbortWithError(c *gin.Context, args ...any) {
	loggerFormat := config.GenGinLoggerFields(c)
	v.Logger = v.Logger.WithFields(loggerFormat)

	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusBadRequest
		msg.Code = ext.MissingParameters
		msg.Desc = "请求头错误"
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(ext.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		v.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	v.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, ext.GenErrMsg(formattedErr))
}
