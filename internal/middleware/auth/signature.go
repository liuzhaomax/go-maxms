package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
)

func (auth *Auth) ValidateSignature() gin.HandlerFunc {
	cfg := core.GetConfig()
	return func(c *gin.Context) {
		userId, _ := c.Cookie(core.UserID)         // 允许为空，不需处理err
		nonce := c.Request.Header.Get(core.SpanId) // 在core.SetHeadersForDownstream之前，所以是spanId
		// 根据headers里给定的信息，生成签名并比对
		signatureGen := core.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonce)
		result := auth.RedisClient.SAdd(context.Background(), core.Signature, signatureGen)
		// 1是set里原来没有，加入成功，0是set里原来有，加入失败
		if (*result).Val() == 0 {
			auth.Logger.WithField(core.FAILURE, core.GetFuncName()).Info(genSignatureErrMsg(nil))
			c.AbortWithStatusJSON(http.StatusUnauthorized, genSignatureErrMsg(nil))
			return
		}
		c.Next()
	}
}

func genSignatureErrMsg(err error) error {
	return core.FormatError(core.Unauthorized, "签名验证失败", err)
}
