package core

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"github.com/redis/go-redis/v9"
	"github.com/satori/go.uuid"
	"strings"
)

func TraceID() string {
	return UUIDInUpper()
}

func SpanID() string {
	return UUIDInLower()
}

func UUIDInUpper() string {
	return strings.ToUpper(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}

func UUIDInLower() string {
	return strings.ToLower(strings.ReplaceAll(uuid.NewV1().String(), "-", ""))
}

func ShortUUID() string {
	return shortuuid.New()
}

func GetClientIP(c *gin.Context) string {
	clientIP := c.Request.Header.Get(ClientIp)
	if clientIP == EmptyString {
		clientIP = c.ClientIP()
		c.Request.Header.Set(ClientIp, clientIP)
	}
	return clientIP
}

func GetUserAgent(c *gin.Context) string {
	userAgent := c.Request.Header.Get(UserAgent)
	if userAgent == EmptyString {
		userAgent = c.Request.UserAgent()
		c.Request.Header.Set(UserAgent, userAgent)
	}
	return userAgent
}

func ValidateHeaders(c *gin.Context) error {
	if c.Request.Header.Get(TraceId) == EmptyString ||
		c.Request.Header.Get(SpanId) == EmptyString {
		return errors.New("缺失链路信息")
	}
	if c.Request.Header.Get(AppId) == EmptyString {
		return errors.New("缺失接口签名信息")
	}
	return nil
}

func SetHeadersForDownstream(c *gin.Context, downstreamName string, client *redis.Client) error {
	c.Request.Header.Set(ClientIp, c.Request.Header.Get(ClientIp))
	c.Request.Header.Set(UserAgent, c.Request.Header.Get(UserAgent))
	c.Request.Header.Set(TraceId, c.Request.Header.Get(TraceId))
	c.Request.Header.Set(ParentId, c.Request.Header.Get(SpanId))
	c.Request.Header.Set(SpanId, SpanID())
	c.Request.Header.Set(AppId, cfg.App.Id)
	userId, _ := c.Cookie(UserID)
	nonce := c.Request.Header.Get(ParentId) // ParentId已被赋值为req headers里的spanId
	downstreamAppId := EmptyString
	downstreamAppSecret := EmptyString
	for _, downstream := range cfg.Downstream {
		if downstream.Name == downstreamName {
			downstreamAppId = downstream.Id
			downstreamAppSecret = downstream.Secret
			break
		}
	}
	// 生成签名并写入redis
	signature := GenAppSignature(downstreamAppId, downstreamAppSecret, userId, nonce)
	result, err := client.SAdd(context.Background(), Signature, signature).Result()
	if err != nil {
		return FormatError(CacheDenied, "缓存写入失败", err)
	}
	if result == 0 {
		return FormatError(CacheDenied, "缓存写入失败", errors.New("set已有该值"))
	}
	return nil
}
