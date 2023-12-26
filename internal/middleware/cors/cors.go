package cors

import (
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
)

func Cors() gin.HandlerFunc {
	whiteList := core.GetConfig().WhiteList // TODO 需要测试
	return func(ctx *gin.Context) {
		if core.In(whiteList, ctx.Request.Header.Get("Origin")) {
			ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("Origin"))
		}
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, Set-Cookie, X-Requested-With, Access-Control-Allow-Origin, Content-Security-Policy")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		// let go all options request
		method := ctx.Request.Method
		if method == "OPTIONS" {
			ctx.Header("Access-Control-Max-Age", "86400") // one day
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}
