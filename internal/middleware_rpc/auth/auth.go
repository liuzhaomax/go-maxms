package auth

import (
	"context"
	"strings"

	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var AuthRPCSet = wire.NewSet(wire.Struct(new(AuthRPC), "*"))

type AuthRPC struct {
	Logger *logrus.Logger
	Redis  *redis.Client
}

func (auth *AuthRPC) ValidateToken(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		auth.AbortWithError(ctx, errors.New("元数据解析错误"))

		return resp, err
	}

	j := config.NewJWT()
	// token in md
	if len(md[config.Authorization]) == 0 {
		auth.AbortWithError(ctx, errors.New("没找到token"))

		return resp, err
	}

	headerToken := md[config.Authorization][0]
	if headerToken == "" {
		auth.AbortWithError(ctx, errors.New("没找到token"))

		return resp, err
	}

	headerDecryptedToken, err := ext.RSADecrypt(config.GetPrivateKey(), headerToken)
	if err != nil {
		auth.AbortWithError(ctx, err)

		return resp, err
	}

	headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]

	userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
	if err != nil {
		if err.Error() != config.TokenExpired {
			auth.AbortWithError(ctx, err)

			return resp, err
		}

		refreshedToken, errNew := j.RefreshToken(headerDecryptedTokenRemoveBearer)
		if errNew != nil {
			auth.AbortWithError(ctx, errNew)

			return resp, err
		}

		userID, clientIP, err = j.ParseToken(refreshedToken)
		if err != nil {
			auth.AbortWithError(ctx, err)

			return resp, err
		}
		// 验证refreshedToken
		result := auth.CompareCombination(md, userID, clientIP)
		if !result {
			auth.AbortWithError(ctx, err)

			return resp, err
		}

		resp, err = handler(ctx, req)

		return resp, err
	}
	// 验证headerParsedToken
	result := auth.CompareCombination(md, userID, clientIP)
	if !result {
		auth.AbortWithError(ctx, err)

		return resp, err
	}

	resp, err = handler(ctx, req)

	return resp, err
}

// 验证规则：
// 1. 当前请求IP或是header中的clientIP，与JWT中当初token签发IP相同
// 2. header中的userID与JWT中userID相同
func (auth *AuthRPC) CompareCombination(md metadata.MD, userID string, clientIP string) bool {
	var userIdInMD string
	if len(md[config.UserId]) != 0 {
		userIdInMD = md[config.UserId][0]
	}

	var currentIP string
	if len(md[config.ClientIp]) != 0 {
		currentIP = md[config.ClientIp][0]
	}

	if currentIP == clientIP && userIdInMD == userID {
		return true
	}

	return false
}

func (auth *AuthRPC) AbortWithError(ctx context.Context, args ...any) {
	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Code = ext.Unauthorized
		msg.Desc = "Not authenticated"
		msg.Err = args[0].(error)
	case 3: // 复杂调用：AbortWithError(c, code, desc, err)
		msg.Code = args[0].(ext.Code)
		msg.Desc = args[1].(string)
		msg.Err = args[2].(error)
	default:
		auth.Logger.Error("invalid arguments for AbortWithError")
	}
	// 整理打印
	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	auth.Logger.Error(formattedErr)
}
