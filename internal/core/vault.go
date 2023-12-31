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

var vaultClient *vault.Client

func init() {
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
	vaultClient = client
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

func getSecret(ctx context.Context, secretEngineName string, namespace string, key string) string {
	vaultClient.SetNamespace(namespace)
	// 读取secret
	secret, err := vaultClient.KVv2(secretEngineName).Get(ctx, vaultClient.Namespace())
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ConnectionFailed, fmt.Sprintf("获取%s失败", key), err))
		panic(err)
	}
	value, ok := secret.Data[key].(string)
	if !ok {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ParseIssue, fmt.Sprintf("获取的%s不是字符串", key), err))
		panic(err)
	}
	return value
}

func putSecret(ctx context.Context, secretEngineName string, namespace string, secretData map[string]interface{}) {
	vaultClient.SetNamespace(namespace)
	_, err := vaultClient.KVv2(secretEngineName).Put(ctx, vaultClient.Namespace(), secretData)
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(ConnectionFailed, fmt.Sprintf("写入%s失败", namespace), err))
		panic(err)
	}
}

// GetSecret 使用client token访问vault，获取secret，存入配置对象
func (cfg *Config) GetSecret() {
	ctx := context.Background()
	// 读取jwt_secret
	cfg.App.JWTSecret = getSecret(ctx, Kv, Jwt, Secret)
	// 读取salt
	cfg.App.Salt = getSecret(ctx, Kv, Pwd, Salt)
	// 读取rsa string
	cfg.App.PublicKeyStr = getSecret(ctx, Kv, Rsa, Puk)
	cfg.App.PrivateKeyStr = getSecret(ctx, Kv, Rsa, Prk)
	// 打印日志
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo("secret获取成功"))
}

// PutRSA 新增和修改rsa
func (cfg *Config) PutRSA() {
	ctx := context.Background()
	// 写入rsa
	secretData := map[string]interface{}{
		Puk: cfg.App.PublicKeyStr,
		Prk: cfg.App.PrivateKeyStr,
	}
	putSecret(ctx, Kv, Rsa, secretData)
	// 打印日志
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo(fmt.Sprintf("%s写入成功", Rsa)))
}

// PutSalt 新增和修改salt
func (cfg *Config) PutSalt() {
	ctx := context.Background()
	// 写入salt
	secretData := map[string]interface{}{
		Salt: cfg.App.Salt,
	}
	putSecret(ctx, Kv, Pwd, secretData)
	// 打印日志
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo(fmt.Sprintf("%s写入成功", Salt)))
}
