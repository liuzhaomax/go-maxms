package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (auth *Auth) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := core.NewJWT()
		// token in req header
		headerToken := c.Request.Header.Get(core.Authorization)
		if headerToken == core.EmptyString || len(headerToken) == 0 {
			auth.AbortWithError(c, fmt.Errorf("权限验证失败"))
			return
		}
		headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		if err != nil {
			auth.AbortWithError(c, err)
			return
		}
		headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]
		userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
		if err != nil {
			if err.Error() != core.TokenExpired {
				auth.AbortWithError(c, err)
				return
			}
			refreshedToken, err := j.RefreshToken(headerDecryptedTokenRemoveBearer)
			if err != nil {
				auth.AbortWithError(c, err)
				return
			}
			userID, clientIP, err = j.ParseToken(refreshedToken)
			if err != nil {
				auth.AbortWithError(c, err)
				return
			}
			// 验证refreshedToken
			if !auth.CompareCombination(c, userID, clientIP) {
				auth.AbortWithError(c, fmt.Errorf("权限验证失败"))
				return
			}
			c.Next()
			return
		}
		// 验证headerParsedToken
		if !auth.CompareCombination(c, userID, clientIP) {
			auth.AbortWithError(c, fmt.Errorf("权限验证失败"))
			return
		}
		c.Next()
	}
}

func (auth *Auth) AbortWithError(c *gin.Context, args ...any) {
	msg := &core.MiddlewareMessage{
		StatusCode: 500,
		Code:       core.InternalServerError,
		Desc:       core.EmptyString,
		Err:        nil,
	}
	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusUnauthorized
		msg.Code = core.Unauthorized
		msg.Desc = "Not authenticated"
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(core.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		auth.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := core.FormatError(msg.Code, msg.Desc, msg.Err)
	auth.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, core.GenErrMsg(formattedErr))
}

// 验证规则：
// 1. 当前请求IP或是header中的clientIP，与JWT中当初token签发IP相同
// 2. header中的userID与JWT中userID相同
func (auth *Auth) CompareCombination(c *gin.Context, userID string, clientIP string) bool {
	userIdInHeaders := c.Request.Header.Get(core.UserId)
	currentIP := core.GetClientIP(c)
	if currentIP == clientIP && userIdInHeaders == userID {
		return true
	}
	return false
}
