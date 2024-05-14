package core

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"net/http"
	"time"
)

type Vault struct {
	Address  string `mapstructure:"address"`
	Interval int    `mapstructure:"interval"`
}

// vault access
const (
	Username = "liuzhao"
	Password = "123456"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var vaultClient *vault.Client

func InitVault() {
	// 获取登录vault的token
	token, err := userpassLogin()
	if err != nil {
		LogFailure(VaultDenied, "vault用户登录token获取失败", err)
		panic(err)
	}
	// 创建vault连接客户端
	client, err := vault.NewClient(&vault.Config{Address: cfg.Lib.Vault.Address, HttpClient: httpClient})
	if err != nil {
		LogFailure(Unknown, "创建vault连接失败", err)
		panic(err)
	}
	// 配置请求
	client.SetToken(token)
	vaultClient = client
}

// userpassLogin 使用用户名密码登录vault，获取token
func userpassLogin() (string, error) {
	// 创建vault连接客户端
	client, err := vault.NewClient(&vault.Config{Address: cfg.Lib.Vault.Address, HttpClient: httpClient})
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

func getSecret(ctx context.Context, secretEngineName string, namespace string, key string) (string, error) {
	vaultClient.SetNamespace(namespace)
	// 读取secret
	secret, err := vaultClient.KVv2(secretEngineName).Get(ctx, vaultClient.Namespace())
	if err != nil {
		return EmptyString, err
	}
	value, ok := secret.Data[key].(string)
	if !ok {
		return EmptyString, fmt.Errorf("获取的%s不是字符串", key)
	}
	return value, nil
}

func putSecret(ctx context.Context, secretEngineName string, namespace string, secretData map[string]interface{}) error {
	vaultClient.SetNamespace(namespace)
	_, err := vaultClient.KVv2(secretEngineName).Put(ctx, vaultClient.Namespace(), secretData)
	if err != nil {
		return err
	}
	return nil
}

// GetSecret 使用client token访问vault，获取secret，存入配置对象
func (cfg *Config) GetSecret() {
	ctx := context.Background()
	// 读取jwt_secret
	jwtSecret, err := getSecret(ctx, KV, JWT, SECRET)
	if err != nil {
		LogFailure(VaultDenied, "Vault: jwt_secret获取失败", err)
		panic(err)
	}
	cfg.App.JWTSecret = jwtSecret
	// 读取salt
	salt, err := getSecret(ctx, KV, PWD, SALT)
	if err != nil {
		LogFailure(VaultDenied, "Vault: salt获取失败", err)
		panic(err)
	}
	cfg.App.Salt = salt
	// 读取rsa string
	puk, err := getSecret(ctx, KV, RSA, PUK)
	if err != nil {
		LogFailure(VaultDenied, "Vault: puk获取失败", err)
		panic(err)
	}
	prk, err := getSecret(ctx, KV, RSA, PRK)
	if err != nil {
		LogFailure(VaultDenied, "Vault: prk获取失败", err)
		panic(err)
	}
	cfg.App.PublicKeyStr = puk
	cfg.App.PrivateKeyStr = prk
	// 读取downstream app id 和 secret
	for i, downstream := range cfg.Downstreams {
		cfg.Downstreams[i].Id, err = getSecret(ctx, KV, fmt.Sprintf("%s/%s", APP, downstream.Name), ID)
		if err != nil {
			LogFailure(VaultDenied, "Vault: downstream信息获取失败", err)
			panic(err)
		}
		cfg.Downstreams[i].Secret, err = getSecret(ctx, KV, fmt.Sprintf("%s/%s", APP, downstream.Name), SECRET)
		if err != nil {
			LogFailure(VaultDenied, "Vault: downstream信息获取失败", err)
			panic(err)
		}
	}
}

// PutRSA 新增和修改rsa
func (cfg *Config) PutRSA() {
	ctx := context.Background()
	// 写入rsa
	secretData := map[string]interface{}{
		PUK: cfg.App.PublicKeyStr,
		PRK: cfg.App.PrivateKeyStr,
	}
	err := putSecret(ctx, KV, RSA, secretData)
	if err != nil {
		LogFailure(VaultDenied, fmt.Sprintf("Vault: %s写入失败", RSA), err)
		panic(err)
	}
	// 打印日志
	LogSuccess(fmt.Sprintf("Vault: %s写入成功", RSA))
}

// PutSalt 新增和修改salt
func (cfg *Config) PutSalt() {
	ctx := context.Background()
	// 写入salt
	secretData := map[string]interface{}{
		SALT: cfg.App.Salt,
	}
	err := putSecret(ctx, KV, PWD, secretData)
	if err != nil {
		LogFailure(VaultDenied, fmt.Sprintf("Vault: %s写入失败", SALT), err)
		panic(err)
	}
	// 打印日志
	LogSuccess(fmt.Sprintf("Vault: %s写入成功", SALT))
}

// PutAppSecret 新增修改AppId和AppSecret
func (cfg *Config) PutAppSecret() {
	ctx := context.Background()
	// 写入服务信息
	secretData := map[string]interface{}{
		ID:     cfg.App.Id,
		SECRET: cfg.App.Secret,
	}
	namespace := fmt.Sprintf("%s/%s", APP, cfg.App.Name)
	err := putSecret(ctx, KV, namespace, secretData)
	if err != nil {
		LogFailure(VaultDenied, fmt.Sprintf("Vault: %s写入失败", namespace), err)
		panic(err)
	}
	// 打印日志
	LogSuccess(fmt.Sprintf("Vault: %s写入成功", namespace))
}
