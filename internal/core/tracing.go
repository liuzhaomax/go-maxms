package core

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
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
		c.Request.Header.Get(SpanId) == EmptyString ||
		c.Request.Header.Get(RequestId) == EmptyString {
		return errors.New("缺失链路信息")
	}
	if c.Request.Header.Get(AppId) == EmptyString {
		return errors.New("缺失接口签名信息")
	}
	return nil
}

func ValidateMetadata(md metadata.MD) error {
	if SelectFromMetadata(md, TraceId) == EmptyString ||
		SelectFromMetadata(md, SpanId) == EmptyString ||
		SelectFromMetadata(md, RequestId) == EmptyString {
		return errors.New("缺失链路信息")
	}
	if SelectFromMetadata(md, AppId) == EmptyString ||
		SelectFromMetadata(md, RequestURI) == EmptyString {
		return errors.New("缺失接口签名信息")
	}
	return nil
}

func SelectFromMetadata(md metadata.MD, key string) string {
	for k, v := range md {
		if strings.EqualFold(k, key) {
			return v[0]
		}
	}
	return EmptyString
}

func SetMetadataForDownstreamFromHttpHeaders(ctx context.Context, c *gin.Context, downstreamName string) (context.Context, error) {
	var mdMap = map[string]string{}
	mdMap[ClientIp] = c.Request.Header.Get(ClientIp)
	mdMap[UserAgent] = c.Request.Header.Get(UserAgent)
	mdMap[RequestId] = c.Request.Header.Get(RequestId)
	mdMap[TraceId] = c.Request.Header.Get(TraceId)
	mdMap[ParentId] = c.Request.Header.Get(SpanId)
	mdMap[SpanId] = c.Request.Header.Get(SpanId)
	mdMap[AppId] = cfg.App.Id
	mdMap[UserId] = c.Request.Header.Get(UserId)
	mdMap[Authorization] = c.Request.Header.Get(Authorization)
	mdMap[RequestURI] = c.Request.RequestURI
	mdMap[UberTraceId] = c.Request.Header.Get(UberTraceId)
	// 接口签名用
	userId := c.Request.Header.Get(UserId)
	nonce := c.Request.Header.Get(SpanId)
	downstreamAppId := EmptyString
	downstreamAppSecret := EmptyString
	for _, downstream := range cfg.Downstreams {
		if downstream.Name == downstreamName {
			downstreamAppId = downstream.Id
			downstreamAppSecret = downstream.Secret
			break
		}
	}
	// 生成签名并写入redis
	signature := GenAppSignature(downstreamAppId, downstreamAppSecret, userId, nonce)
	mdMap[Signature] = signature
	md := metadata.New(mdMap)
	newCtx := metadata.NewOutgoingContext(ctx, md)
	return newCtx, nil
}

func SetHeadersForDownstream(c *gin.Context, downstreamName string) error {
	c.Request.Header.Set(ClientIp, c.Request.Header.Get(ClientIp))
	c.Request.Header.Set(UserAgent, c.Request.Header.Get(UserAgent))
	c.Request.Header.Set(RequestId, c.Request.Header.Get(RequestId))
	c.Request.Header.Set(TraceId, c.Request.Header.Get(TraceId))
	c.Request.Header.Set(ParentId, c.Request.Header.Get(SpanId))
	c.Request.Header.Set(SpanId, c.Request.Header.Get(SpanId))
	c.Request.Header.Set(AppId, cfg.App.Id)
	c.Request.Header.Set(Authorization, c.Request.Header.Get(Authorization))
	c.Request.Header.Set(UserId, c.Request.Header.Get(UserId))
	// 接口签名
	userId := c.Request.Header.Get(UserId)
	nonce := c.Request.Header.Get(SpanId) // 当前ms的spanId
	downstreamAppId := EmptyString
	downstreamAppSecret := EmptyString
	for _, downstream := range cfg.Downstreams {
		if downstream.Name == downstreamName {
			downstreamAppId = downstream.Id
			downstreamAppSecret = downstream.Secret
			break
		}
	}
	signature := GenAppSignature(downstreamAppId, downstreamAppSecret, userId, nonce)
	c.Request.Header.Set(Signature, signature)
	return nil
}
