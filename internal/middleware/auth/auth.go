package auth

import (
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
	Logger      *logrus.Logger
	RedisClient *redis.Client
}

func (auth *Auth) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := core.NewJWT()
		// token in req header
		headerToken := c.Request.Header.Get(core.Authorization)
		if headerToken == core.EmptyString || len(headerToken) == 0 {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(nil))
			return
		}
		headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		if err != nil {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
			return
		}
		headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]
		userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
		if err != nil {
			if err.Error() != core.TokenExpired {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
				return
			}
			refreshedToken, err := j.RefreshToken(headerDecryptedTokenRemoveBearer)
			if err != nil {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
				return
			}
			userID, clientIP, err = j.ParseToken(refreshedToken)
			if err != nil {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
				return
			}
			// 验证refreshedToken
			result := auth.CompareCombination(c, userID, clientIP)
			if result == false {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
				return
			}
			c.Next()
			return
		}
		// 验证headerParsedToken
		result := auth.CompareCombination(c, userID, clientIP)
		if result == false {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genTokenErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genTokenErrMsg(err))
			return
		}
		c.Next()
	}
}

// 验证规则：
// 1. 当前请求IP或是header中的clientIP，与JWT中当初token签发IP相同
// 2. cookie中的userID与JWT中userID相同
func (auth *Auth) CompareCombination(c *gin.Context, userID string, clientIP string) bool {
	userIDInCookie, _ := c.Cookie(core.UserID)
	currentIP := core.GetClientIP(c)
	if currentIP == clientIP && userIDInCookie == userID {
		return true
	}
	return false
}

func genTokenErrMsg(err error) error {
	return core.FormatError(core.Unauthorized, "权限验证失败", err)
}
