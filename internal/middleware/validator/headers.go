package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
)

var ValidatorSet = wire.NewSet(wire.Struct(new(Validator), "*"))

type Validator struct {
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (v *Validator) ValidateHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := core.ValidateHeaders(c)
		if err != nil {
			v.AbortWithError(c, err)
			return
		}
		c.Next()
	}
}

func (v *Validator) AbortWithError(c *gin.Context, args ...any) {
	msg := &core.MiddlewareMessage{
		StatusCode: 500,
		Code:       core.InternalServerError,
		Desc:       core.EmptyString,
		Err:        nil,
	}
	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusBadRequest
		msg.Code = core.MissingParameters
		msg.Desc = "请求头错误"
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(core.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		v.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := core.FormatError(msg.Code, msg.Desc, msg.Err)
	v.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, core.GenErrMsg(formattedErr))
}
