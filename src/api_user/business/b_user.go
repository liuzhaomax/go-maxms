package business

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"github.com/liuzhaomax/go-maxms/src/api_user/schema"
	"github.com/redis/go-redis/v9"
	"time"
)

var BusinessUserSet = wire.NewSet(wire.Struct(new(BusinessUser), "*"))

type BusinessUser struct {
	Model *model.ModelUser
	Tx    *core.Trans
	Redis *redis.Client
}

func (b *BusinessUser) PostLogin(c *gin.Context) (*schema.TokenRes, error) {
	res := &schema.TokenRes{}
	loginReq := &schema.LoginReq{}
	err := c.ShouldBind(loginReq)
	if err != nil {
		return res, core.FormatError(core.ParseIssue, "请求体无效", err)
	}
	decryptedUsername, err := core.RSADecrypt(core.GetPrivateKey(), loginReq.Username)
	loginReq.Username = decryptedUsername
	if err != nil {
		return res, core.FormatError(core.PermissionDenied, "请求体解码异常", err)
	}
	decryptedPassword, err := core.RSADecrypt(core.GetPrivateKey(), loginReq.Password)
	loginReq.Password = decryptedPassword
	if err != nil {
		return res, core.FormatError(core.PermissionDenied, "请求体解码异常", err)
	}
	user := &model.User{}
	err = b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err = b.Model.QueryUserByUsername(c, loginReq.Username, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return res, core.FormatError(core.PermissionDenied, "登录失败", err)
	}
	// loginReq.Password是从SGW经过RSA解码后得到密码
	result := core.VerifyEncodedPwd(loginReq.Password, core.GetConfig().App.Salt, user.Password)
	if !result {
		return res, core.FormatError(core.PermissionDenied, "登录验证失败", err)
	}
	// 定义过期时长
	maxAge := 60 * 60 * 24 * 7 // 一周
	duration := time.Second * time.Duration(maxAge)
	// 生成Bearer jwt，使用userID与ip签发
	j := core.NewJWT()
	token, err := j.GenerateToken(user.UserID, core.GetClientIP(c), duration)
	if err != nil {
		return res, core.FormatError(core.PermissionDenied, "Token生成失败", err)
	}
	bearerToken := core.Bearer + token
	// 对Bearer jwt 进行RSA加密
	encryptedBearerToken, err := core.RSAEncrypt(core.GetPublicKey(), bearerToken)
	if err != nil {
		return res, core.FormatError(core.PermissionDenied, "Token加密失败", err)
	}
	// 拼接响应
	res.Token = encryptedBearerToken
	res.UserID = user.UserID
	return res, nil
}

func (b *BusinessUser) DeleteLogin(c *gin.Context) error {
	// maxAge := int(time.Millisecond)
	// domain := core.GetConfig().App.Domain
	// c.SetSameSite(http.SameSiteNoneMode)
	// c.SetCookie(
	//     core.UserID,
	//     core.EmptyString,
	//     maxAge,
	//     "/",
	//     domain,
	//     true,
	//     true)
	return nil
}

func (b *BusinessUser) GetUserByUserID(c *gin.Context) (*schema.UserRes, error) {
	userID := c.Param(core.UserID)
	user := &model.User{}
	err := b.Tx.ExecTrans(c, func(ctx context.Context) error {
		err := b.Model.QueryUserByUserID(c, userID, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	userRes := schema.MapUser2UserRes(user)
	return userRes, nil
}
