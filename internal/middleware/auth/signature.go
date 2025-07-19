package auth

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
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
			auth.AbortWithError(c, fmt.Errorf("签名不匹配"))
			return
		}
		result, err := auth.Redis.SAdd(context.Background(), core.Nonce, nonce).Result()
		// 如果直接使用返回值，(*result).Val()，1是set里原来没有，加入成功，0是set里原来有，加入失败
		if err != nil {
			auth.AbortWithError(c, err)
			return
		}
		if result == 0 {
			auth.AbortWithError(c, fmt.Errorf("set已有该值"))
			return
		}
		// 设置过期时间
		err = auth.Redis.Expire(context.Background(), core.Nonce, time.Second*5).Err()
		if err != nil {
			auth.AbortWithError(c, err)
			return
		}
		c.Next()
	}
}
