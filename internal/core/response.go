package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ResponseSet = wire.NewSet(wire.Struct(new(Response), "*"), wire.Bind(new(IResponse), new(*Response)))

type IResponse interface {
	ResSuccess(*gin.Context, any)
	ResFailure(*gin.Context, int, Code, string, error)
	ResSuccessForRPC(context.Context)
	ResFailureForRPC(context.Context, Code, string, error)
}

type Response struct {
	Logger ILogger
}

func (res *Response) ResSuccess(c *gin.Context, sth any) {
	res.Logger.SucceedWithField(c, "响应正常")
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

func (res *Response) ResFailure(c *gin.Context, statusCode int, code Code, desc string, err error) {
	res.Logger.FailWithField(c, code, fmt.Sprintf("响应异常: %s", desc), err)
	res.ResJson(c, statusCode, gin.H{
		"status": gin.H{
			"code": code,
			"desc": desc,
			"err":  err.Error(),
		},
	})
}

func (res *Response) ResJson(c *gin.Context, statusCode int, sth any) {
	c.JSON(statusCode, sth)
}

func (res *Response) ResSuccessForRPC(ctx context.Context) {
	res.Logger.SucceedWithFieldForRPC(ctx, "响应正常")
}

func (res *Response) ResFailureForRPC(ctx context.Context, code Code, desc string, err error) {
	res.Logger.FailWithFieldForRPC(ctx, code, fmt.Sprintf("响应异常: %s", desc), err)
}
