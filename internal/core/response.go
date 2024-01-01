package core

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var ResponseSet = wire.NewSet(wire.Struct(new(Response), "*"), wire.Bind(new(IResponse), new(*Response)))

type IResponse interface {
	ResSuccess(*gin.Context, string, interface{})
	ResFailure(*gin.Context, string, int, Code, string, error)
}

type Response struct {
	Logger *logrus.Logger
}

func (res *Response) ResSuccess(c *gin.Context, funcName string, sth interface{}) {
	res.Logger.WithField(SUCCESS, funcName).WithField("trace_id", c.Request.Header.Get(TraceId)).Debug(FormatInfo("响应成功"))
	if sth != nil {
		res.ResJson(c, 200, gin.H{
			"status": gin.H{
				"code": OK,
				"desc": "success",
			},
			"data": sth,
		})
		return
	}
	res.ResJson(c, 200, gin.H{
		"status": gin.H{
			"code": OK,
			"desc": "success",
		},
	})
}

func (res *Response) ResFailure(c *gin.Context, funcName string, statusCode int, code Code, desc string, err error) {
	res.Logger.WithField(FAILURE, funcName).WithField("trace_id", c.Request.Header.Get(TraceId)).Debug(FormatError(code, desc, err))
	res.ResJson(c, statusCode, gin.H{
		"status": gin.H{
			"code": code,
			"desc": desc,
			"err":  err.Error(),
		},
	})
}

func (res *Response) ResJson(c *gin.Context, statusCode int, sth interface{}) {
	c.JSON(statusCode, sth)
}
