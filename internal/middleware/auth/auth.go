package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger *logrus.Logger
}

func (auth *Auth) CheckTokens() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	j := core.NewJWT()
		//	// token in req header
		//	headerB64Token := c.Request.Header.Get("Authorization")
		//	if headerB64Token == "" || len(headerB64Token) == 0 {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, nil))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, nil))
		//		return
		//	}
		//	headerToken, err := core.BASE64DecodeStr(headerB64Token)
		//	if err != nil {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
		//	if err != nil {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	headerParsedToken, err := j.ParseToken(headerDecryptedToken)
		//	if err != nil {
		//		if err.Error() == core.TokenExpired {
		//			auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//			c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//			return
		//		}
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	// token in req cookie
		//	cookieB64Token, err := c.Cookie("TOKEN")
		//	if cookieB64Token == "" || len(cookieB64Token) == 0 {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, nil))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	cookieToken, err := core.BASE64DecodeStr(cookieB64Token)
		//	if err != nil {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	cookieDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), cookieToken)
		//	if err != nil {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	cookieParsedToken, err := j.ParseToken(cookieDecryptedToken)
		//	if err != nil {
		//		if err.Error() == core.TokenExpired {
		//			auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//			c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//			return
		//		}
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, err))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, err))
		//		return
		//	}
		//	// checking tokens info
		//	if headerParsedToken != cookieParsedToken {
		//		auth.ILogger.LogFailure(core.GetFuncName(), core.FormatError(206, nil))
		//		c.AbortWithStatusJSON(http.StatusUnauthorized, core.FormatError(206, nil))
		//		return
		//	}
		//	c.Next()
	}
}
