package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
	"time"
)

func (auth *Auth) ValidateSignature() gin.HandlerFunc {
	cfg := core.GetConfig()
	return func(c *gin.Context) {
		userId := c.Request.Header.Get(core.UserId)
		nonceForValidation := c.Request.Header.Get(core.ParentId)
		nonce := c.Request.Header.Get(core.SpanId)
		if nonce == core.EmptyString {
			nonce = core.ShortUUID()
		}
		// 根据headers里给定的信息，生成签名并比对
		signatureGen := core.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonceForValidation)
		signatureHeader := c.Request.Header.Get(core.Signature)
		if signatureGen != signatureHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "签名验证失败", errors.New("签名不匹配")))
			return
		}
		result, err := auth.Redis.SAdd(context.Background(), core.Nonce, nonce).Result()
		// 如果直接使用返回值，(*result).Val()，1是set里原来没有，加入成功，0是set里原来有，加入失败
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "签名验证失败", err))
			return
		}
		if result == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "签名验证失败", errors.New("set已有该值")))
			return
		}
		// 设置过期时间
		err = auth.Redis.Expire(context.Background(), core.Nonce, time.Second*5).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "签名过期时间设置失败", err))
			return
		}
		auth.GenOkMsg(c, "签名已写入缓存")
		c.Next()
	}
}
