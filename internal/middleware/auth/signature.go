package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

func (auth *Auth) ValidateSignature() gin.HandlerFunc {
	cfg := core.GetConfig()

	return func(c *gin.Context) {
		userId := c.Request.Header.Get(config.UserId)
		nonceForValidation := c.Request.Header.Get(config.ParentId)

		nonce := c.Request.Header.Get(config.SpanId)
		if nonce == "" {
			nonce = ext.ShortUUID()
		}
		// 根据headers里给定的信息，生成签名并比对
		signatureGen := ext.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonceForValidation)

		signatureHeader := c.Request.Header.Get(config.Signature)

		if signatureGen != signatureHeader {
			auth.AbortWithError(c, errors.New("签名不匹配"))

			return
		}

		result, err := auth.Redis.SAdd(context.Background(), config.Nonce, nonce).Result()
		// 如果直接使用返回值，(*result).Val()，1是set里原来没有，加入成功，0是set里原来有，加入失败
		if err != nil {
			auth.AbortWithError(c, err)

			return
		}

		if result == 0 {
			auth.AbortWithError(c, errors.New("set已有该值"))

			return
		}
		// 设置过期时间
		err = auth.Redis.Expire(context.Background(), config.Nonce, time.Second*5).Err()
		if err != nil {
			auth.AbortWithError(c, err)

			return
		}

		c.Next()
	}
}
