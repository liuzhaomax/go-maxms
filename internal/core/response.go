package core

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WrapperHandle = func(c *gin.Context) (any, error)

func WrapperRes(handle WrapperHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handle(c)
		if err != nil {
			var apiError *Error
			errors.As(err, &apiError)
			c.JSON(http.StatusOK, gin.H{
				"status": gin.H{
					"code": apiError.Code,
					"desc": apiError.Desc,
				},
				"data": data,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"code": 0,
				"desc": "success",
			},
			"data": data,
		})
	}
}

func GenErrMsg(err error) any {
	var apiError *Error
	errors.As(err, &apiError)
	return gin.H{
		"status": gin.H{
			"code": apiError.Code,
			"desc": apiError.Desc,
		},
		"data": nil,
	}
}

type MiddlewareMessage struct {
	StatusCode int
	Code       Code
	Desc       string
	Err        error
}
