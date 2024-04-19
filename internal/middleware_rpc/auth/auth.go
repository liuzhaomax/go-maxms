package auth

import (
	"context"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

var AuthRPCSet = wire.NewSet(wire.Struct(new(AuthRPC), "*"))

type AuthRPC struct {
	Logger core.ILogger
	Redis  *redis.Client
}

func (auth *AuthRPC) ValidateToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = auth.GenErrMsg(ctx, "元数据解析错误", err)
		return
	}
	j := core.NewJWT()
	// token in md
	if len(md[core.Authorization]) == 0 {
		err = auth.GenErrMsg(ctx, "权限验证失败", errors.New("没找到token"))
		return
	}
	headerToken := md[core.Authorization][0]
	if headerToken == core.EmptyString {
		err = auth.GenErrMsg(ctx, "权限验证失败", errors.New("没找到token"))
		return
	}
	headerDecryptedToken, err := core.RSADecrypt(core.GetPrivateKey(), headerToken)
	if err != nil {
		err = auth.GenErrMsg(ctx, "权限验证失败", err)
		return
	}
	headerDecryptedTokenRemoveBearer := (strings.Split(headerDecryptedToken, " "))[1]
	userID, clientIP, err := j.ParseToken(headerDecryptedTokenRemoveBearer)
	if err != nil {
		if err.Error() != core.TokenExpired {
			err = auth.GenErrMsg(ctx, "权限验证失败", err)
			return
		}
		refreshedToken, errNew := j.RefreshToken(headerDecryptedTokenRemoveBearer)
		if errNew != nil {
			err = auth.GenErrMsg(ctx, "权限验证失败", errNew)
			return
		}
		userID, clientIP, err = j.ParseToken(refreshedToken)
		if err != nil {
			err = auth.GenErrMsg(ctx, "权限验证失败", err)
			return
		}
		// 验证refreshedToken
		result := auth.CompareCombination(md, userID, clientIP)
		if !result {
			err = auth.GenErrMsg(ctx, "权限验证失败", err)
			return
		}
		resp, err = handler(ctx, req)
		return
	}
	// 验证headerParsedToken
	result := auth.CompareCombination(md, userID, clientIP)
	if !result {
		err = auth.GenErrMsg(ctx, "权限验证失败", err)
		return
	}
	resp, err = handler(ctx, req)
	return
}

// 验证规则：
// 1. 当前请求IP或是header中的clientIP，与JWT中当初token签发IP相同
// 2. header中的userID与JWT中userID相同
func (auth *AuthRPC) CompareCombination(md metadata.MD, userID string, clientIP string) bool {
	var userIdInMD string
	if len(md[core.UserId]) != 0 {
		userIdInMD = md[core.UserId][0]
	}
	var currentIP string
	if len(md[core.ClientIp]) != 0 {
		currentIP = md[core.ClientIp][0]
	}
	if currentIP == clientIP && userIdInMD == userID {
		return true
	}
	return false
}

func (auth *AuthRPC) GenOkMsg(ctx context.Context, desc string) string {
	auth.Logger.SucceedWithFieldForRPC(ctx, desc)
	return core.FormatInfo(desc)
}

func (auth *AuthRPC) GenErrMsg(ctx context.Context, desc string, err error) error {
	auth.Logger.FailWithFieldForRPC(ctx, core.Unauthorized, desc, err)
	return core.FormatError(core.Unauthorized, desc, err)
}
