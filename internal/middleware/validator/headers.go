package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var ValidatorSet = wire.NewSet(wire.Struct(new(Validator), "*"))

type Validator struct {
	Logger core.ILogger
	Redis  *redis.Client
}

func (v *Validator) ValidateHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := core.ValidateHeaders(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, v.GenErrMsg(c, "请求头错误", err))
		}
		c.Next()
	}
}

func (v *Validator) GenOkMsg(c *gin.Context, desc string) string {
	v.Logger.SucceedWithField(c, desc)
	return core.FormatInfo(desc)
}

func (v *Validator) GenErrMsg(c *gin.Context, desc string, err error) error {
	v.Logger.FailWithField(c, core.MissingParameters, desc, err)
	return core.FormatError(core.MissingParameters, desc, err)
}
