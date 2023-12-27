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

func userpassLogin() (string, error) {
	// create a vault client
	client, err := vault.NewClient(&vault.Config{Address: VaultAddr, HttpClient: httpClient})
	if err != nil {
		return "", err
	}
	// to pass the Password
	options := map[string]interface{}{
		"password": Password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", Username)
	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}
	token := secret.Auth.ClientToken
	return token, nil
}

func (cfg *Config) GetSecret() {
	// Authenticate and obtain a token
	token, err := userpassLogin()
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(PermissionDenied, "vault用户登录token获取失败", err))
		panic(err)
	}
	// Create a new Vault client with the obtained token
	client, err := vault.NewClient(&vault.Config{Address: VaultAddr, HttpClient: httpClient})
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "创建vault连接失败", err))
		panic(err)
	}
	// Set the token for subsequent requests
	client.SetToken(token)
	client.SetNamespace("dev")
	kvPath := "kv"
	ctx := context.Background()

	// Write a secret
	//secretData := map[string]interface{}{
	//	"jwt_secret": "987654",
	//}
	//_, err = client.KVv2(kvPath).Put(ctx, client.Namespace(), secretData)
	//if err != nil {
	//	log.Fatalf("unable to write secret: %v", err)
	//}
	//
	//fmt.Println("Secret written successfully.")

	// Make a read request to Vault
	secret, err := client.KVv2(kvPath).Get(ctx, client.Namespace())
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Panic(FormatError(PermissionDenied, "获取secret失败", err))
		panic(err)
	}
	cfg.App.JWTSecret = secret.Data["jwt_secret"].(string)
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo("secret获取成功"))
}
