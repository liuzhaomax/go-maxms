package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"strconv"
)

func Str2Uint32(str string) (uint32, error) {
	if str == "" {
		str = "0"
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

func SetHeaders(c *gin.Context) error {
	if c.Request.Header.Get(core.TraceId) == core.EmptyString || c.Request.Header.Get(core.SpanId) == core.EmptyString {
		return errors.New("缺失链路信息")
	}
	c.Request.Header.Set(core.TraceId, c.Request.Header.Get(core.TraceId))
	c.Request.Header.Set(core.ParentId, c.Request.Header.Get(core.SpanId))
	c.Request.Header.Set(core.SpanId, core.SpanID())
	c.Request.Header.Set(core.ClientIp, c.ClientIP())
	return nil
}
