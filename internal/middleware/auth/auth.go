package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger core.ILogger
	Redis  *redis.Client
}

func (auth *Auth) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := core.NewJWT()
		// token in req header
		headerToken := c.Request.Header.Get(core.Authorization)
		if headerToken == core.EmptyString || len(headerToken) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", errors.New("没找到token")))
			return
		}
		headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
			return
		}
		headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]
		userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
		if err != nil {
			if err.Error() != core.TokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
				return
			}
			refreshedToken, err := j.RefreshToken(headerDecryptedTokenRemoveBearer)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
				return
			}
			userID, clientIP, err = j.ParseToken(refreshedToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
				return
			}
			// 验证refreshedToken
			result := auth.CompareCombination(c, userID, clientIP)
			if !result {
				c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
				return
			}
			c.Next()
			return
		}
		// 验证headerParsedToken
		result := auth.CompareCombination(c, userID, clientIP)
		if !result {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "权限验证失败", err))
			return
		}
		c.Next()
	}
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

func (auth *Auth) GenOkMsg(c *gin.Context, desc string) string {
	auth.Logger.SucceedWithField(c, desc)
	return core.FormatInfo(desc)
}

func (auth *Auth) GenErrMsg(c *gin.Context, desc string, err error) error {
	auth.Logger.FailWithField(c, core.Unauthorized, desc, err)
	return core.FormatError(core.Unauthorized, desc, err)
}
