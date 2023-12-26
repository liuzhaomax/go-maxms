package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms-template/internal/core"
	"github.com/sirupsen/logrus"
	"net/http"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger *logrus.Logger
}

func (auth *Auth) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := core.NewJWT()
		// token in req header
		headerB64Token := c.Request.Header.Get("Authorisation")
		if headerB64Token == "" || len(headerB64Token) == 0 {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(nil))
			return
		}
		headerToken, err := core.BASE64DecodeStr(headerB64Token)
		if err != nil {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		if err != nil {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		headerParsedToken, err := j.ParseToken(headerDecryptedToken)
		if err != nil {
			if err.Error() == core.TokenExpired {
				auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		// token in req cookie
		cookieB64Token, err := c.Cookie("TOKEN")
		if cookieB64Token == "" || len(cookieB64Token) == 0 {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		cookieToken, err := core.BASE64DecodeStr(cookieB64Token)
		if err != nil {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		cookieDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), cookieToken)
		if err != nil {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		cookieParsedToken, err := j.ParseToken(cookieDecryptedToken)
		if err != nil {
			if err.Error() == core.TokenExpired {
				auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
				return
			}
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(err))
			return
		}
		// checking tokens info
		if headerParsedToken != cookieParsedToken {
			auth.Logger.WithField("失败方法", core.GetFuncName()).Debug(genErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genErrMsg(nil))
			return
		}
		c.Next()
	}
}

func genErrMsg(err error) string {
	return core.FormatError(core.PermissionDenied, "权限验证失败", err)
}
