package cors

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
)

func Cors() gin.HandlerFunc {
	whiteList := core.GetConfig().App.WhiteList
	return func(c *gin.Context) {
		if core.In(whiteList, c.Request.Header.Get("Origin")) {
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, Set-Cookie, X-Requested-With, Access-Control-Allow-Origin, Content-Security-Policy, Request_id, App_id, User_id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// let go all options request
		method := c.Request.Method
		if method == "OPTIONS" {
			c.Header("Access-Control-Max-Age", "86400") // one day
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
