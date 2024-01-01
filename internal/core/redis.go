package core

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	redisPassword = "123456"
)

func InitRedis() (*redis.Client, func(), error) {
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo("Redis连接启动"))
	client, clean, err := cfg.LoadRedis()
	if err != nil {
		cfg.App.Logger.WithField(FAILURE, GetFuncName()).Fatal(FormatError(Unknown, "Redis连接失败", err))
		return nil, clean, err
	}
	cfg.App.Logger.WithField(SUCCESS, GetFuncName()).Info(FormatInfo("Redis连接成功"))
	return client, clean, nil
}

func (cfg *Config) LoadRedis() (*redis.Client, func(), error) {
	addr := fmt.Sprintf("%s:%s", cfg.Lib.Redis.Host, cfg.Lib.Redis.Port)
	opts := &redis.Options{
		Network:               "",
		Addr:                  addr,
		ClientName:            "",
		Dialer:                nil,
		OnConnect:             nil,
		Protocol:              0,
		Username:              "",
		Password:              redisPassword,
		CredentialsProvider:   nil,
		DB:                    0,
		MaxRetries:            0,
		MinRetryBackoff:       0,
		MaxRetryBackoff:       0,
		DialTimeout:           0,
		ReadTimeout:           0,
		WriteTimeout:          0,
		ContextTimeoutEnabled: false,
		PoolFIFO:              false,
		PoolSize:              0,
		PoolTimeout:           0,
		MinIdleConns:          0,
		MaxIdleConns:          0,
		MaxActiveConns:        0,
		ConnMaxIdleTime:       0,
		ConnMaxLifetime:       0,
		TLSConfig:             nil,
		Limiter:               nil,
		DisableIndentity:      false,
	}
	client := redis.NewClient(opts)
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, nil, err
	}
	clean := func() {
		err := client.Close()
		if err != nil {
			cfg.App.Logger.WithField(FAILURE, GetFuncName()).Fatal(FormatError(Unknown, "Redis断开连接失败", err))
		}
	}
	return client, clean, nil
}

// client.Set(...)
// client.Get(...)
// client.SAdd(...)
