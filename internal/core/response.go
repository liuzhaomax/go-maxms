package core

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WrapperHandler = func(c *gin.Context) (any, error)

func WrapperRes(handle WrapperHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handle(c)
		if err != nil {
			var apiError *ApiError
			errors.As(err, &apiError)
			statusCode := selectStatusCode(apiError.Code)
			c.JSON(statusCode, gin.H{
				"status": gin.H{
					"code": apiError.Code,
					"desc": apiError.Message,
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

func selectStatusCode(customizedCode int) int {
	if customizedCode >= 1000 && customizedCode < 2000 {
		return http.StatusInternalServerError
	}
	if customizedCode >= 2000 && customizedCode < 3000 {
		return http.StatusBadRequest
	}
	if customizedCode >= 10000 {
		return http.StatusFailedDependency
	}
	return http.StatusInternalServerError
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
