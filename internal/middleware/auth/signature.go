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
		userId, _ := c.Cookie(core.UserID)         // 允许为空，不需处理err
		nonce := c.Request.Header.Get(core.SpanId) // 在core.SetHeadersForDownstream之前，所以是spanId
		// 根据headers里给定的信息，生成签名并比对
		signatureGen := core.GenAppSignature(cfg.App.Id, cfg.App.Secret, userId, nonce)
		result, err := auth.Redis.SAdd(context.Background(), core.Signature, signatureGen).Result()
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
		err = auth.Redis.Expire(context.Background(), core.Signature, time.Second*5).Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, auth.GenErrMsg(c, "签名过期时间设置失败", err))
			return
		}
		auth.GenOkMsg(c, "签名已写入缓存")
		c.Next()
	}
}
