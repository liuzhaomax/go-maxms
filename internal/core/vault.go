package core

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"net/http"
	"time"
)

// vault access
const (
	Username  = "liuzhao"
	Password  = "123456"
	VaultAddr = "http://127.0.0.1:8200"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// userpassLogin 使用用户名密码登录vault，获取token
func userpassLogin() (string, error) {
	// 创建vault连接客户端
	client, err := vault.NewClient(&vault.Config{Address: VaultAddr, HttpClient: httpClient})
	if err != nil {
		return "", err
	}
	// 配置密码
	options := map[string]interface{}{
		"password": Password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", Username)
	// PUT 登录vault，获取token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}
	token := secret.Auth.ClientToken
	return token, nil
}

// GetSecret 使用client token访问vault，获取secret，存入配置对象
func (cfg *Config) GetSecret() {
	// 获取登录vault的token
	token, err := userpassLogin()
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ConnectionFailed, "vault用户登录token获取失败", err))
		panic(err)
	}
	// 创建vault连接客户端
	client, err := vault.NewClient(&vault.Config{Address: VaultAddr, HttpClient: httpClient})
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "创建vault连接失败", err))
		panic(err)
	}
	// 配置请求
	client.SetToken(token)
	client.SetNamespace("dev")
	kvPath := "kv"
	ctx := context.Background()

	// 写入secret
	//secretData := map[string]interface{}{
	//	"jwt_secret": "987654",
	//}
	//_, err = client.KVv2(kvPath).Put(ctx, client.Namespace(), secretData)
	//if err != nil {
	//	cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ConnectionFailed, "写入secret失败", err))
	//}

	// 读取secret
	secret, err := client.KVv2(kvPath).Get(ctx, client.Namespace())
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ConnectionFailed, "获取secret失败", err))
		panic(err)
	}
	cfg.App.JWTSecret = secret.Data["jwt_secret"].(string)
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo("secret获取成功"))
}
