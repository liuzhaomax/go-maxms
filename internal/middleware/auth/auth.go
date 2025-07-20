package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var AuthSet = wire.NewSet(wire.Struct(new(Auth), "*"))

type Auth struct {
	Logger *logrus.Entry
	Redis  *redis.Client
}

func (auth *Auth) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		j := config.NewJWT()
		// token in req header
		headerToken := c.Request.Header.Get(config.Authorization)
		if headerToken == "" || len(headerToken) == 0 {
			auth.AbortWithError(c, errors.New("权限验证失败"))

			return
		}

		headerDecryptedToken, err := ext.RSADecrypt(config.GetPrivateKey(), headerToken)
		if err != nil {
			auth.AbortWithError(c, err)

			return
		}

		headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]

		userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
		if err != nil {
			if err.Error() != config.TokenExpired {
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
				auth.AbortWithError(c, errors.New("权限验证失败"))

				return
			}

			c.Next()

			return
		}
		// 验证headerParsedToken
		if !auth.CompareCombination(c, userID, clientIP) {
			auth.AbortWithError(c, errors.New("权限验证失败"))

			return
		}

		c.Next()
	}
}

func (auth *Auth) AbortWithError(c *gin.Context, args ...any) {
	auth.Logger = auth.Logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"client_ip":  config.GetClientIP(c),
		"user_agent": config.GetUserAgent(c),
		"token":      c.GetHeader(config.Authorization),
		"trace_id":   c.GetHeader(config.TraceId),
		"span_id":    c.GetHeader(config.SpanId),
		"parent_id":  c.GetHeader(config.ParentId),
		"app_id":     c.GetHeader(config.AppId),
		"request_id": c.GetHeader(config.RequestId),
		"user_id":    c.GetHeader(config.UserId),
	})

	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.StatusCode = http.StatusUnauthorized
		msg.Code = ext.Unauthorized
		msg.Desc = "Not authenticated"
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(ext.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		auth.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	auth.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, ext.GenErrMsg(formattedErr))
}

// 验证规则：
// 1. 当前请求IP或是header中的clientIP，与JWT中当初token签发IP相同
// 2. header中的userID与JWT中userID相同
func (auth *Auth) CompareCombination(c *gin.Context, userID string, clientIP string) bool {
	userIdInHeaders := c.GetHeader(config.UserId)

	currentIP := config.GetClientIP(c)

	if currentIP == clientIP && userIdInHeaders == userID {
		return true
	}

	return false
}
