package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger *logrus.Logger
}

func (auth *Auth) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := core.NewJWT()
		// token in req header
		headerToken := c.Request.Header.Get(core.Authorization)
		if headerToken == "" || len(headerToken) == 0 {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(nil))
			return
		}
		headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		if err != nil {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]
		userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
		if err != nil {
			if err.Error() != core.TokenExpired {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			refreshedToken, err := j.RefreshToken(headerDecryptedTokenRemoveBearer)
			if err != nil {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			userID, clientIP, err = j.ParseToken(refreshedToken)
			if err != nil {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			// 验证refreshedToken
			result := auth.CompareCombination(c, userID, clientIP)
			if result == false {
				auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			c.Next()
			return
		}
		// 验证headerParsedToken
		result := auth.CompareCombination(c, userID, clientIP)
		if result == false {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		c.Next()
	}
}

// 验证规则：
// 1. 当前请求IP与JWT中当初token签发IP相同
// 2. cookie中的userID与JWT中userID相同
func (auth *Auth) CompareCombination(c *gin.Context, userID string, clientIP string) bool {
	userIDInCookie, _ := c.Cookie(utils.UserID)
	if c.ClientIP() == clientIP && userIDInCookie == userID {
		return true
	}
	return false
}

func genErrMsg(err error) error {
	return core.FormatError(core.Unauthorized, "权限验证失败", err)
}
