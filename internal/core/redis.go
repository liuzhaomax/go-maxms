package core

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	redisPassword = "123456"
)

type Redis struct {
	Endpoint
}

func InitRedis() (*redis.Client, func(), error) {
	LogSuccess("Redis连接启动")
	client, clean, err := cfg.Lib.Redis.LoadRedis()
	if err != nil {
		LogFailure(ConnectionFailed, "Redis连接失败", err)
		return nil, clean, err
	}
	LogSuccess("Redis连接成功")
	return client, clean, nil
}

func (r *Redis) LoadRedis() (*redis.Client, func(), error) {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%s", r.Endpoint.Host, r.Endpoint.Port)
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
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, nil, err
	}
	// 设置必备键的过期时间
	err = client.Expire(ctx, Signature, 24*time.Hour).Err()
	if err != nil {
		LogFailure(CacheDenied, "过期时间设置失败", err)
		return nil, nil, err
	}
	clean := func() {
		err := client.Close()
		if err != nil {
			LogFailure(Unknown, "Redis断开连接失败", err)
		}
	}
	return client, clean, nil
}

// client.Set(...)
// client.Get(...)
// client.SAdd(...)

// func logAndExec(cmd func(context.Context, *redis.Client, ...interface{}) *redis.Cmd,
//    ctx context.Context, client *redis.Client, args []interface{}) (string, error) {
//
//    // 在执行前记录日志
//    log.Printf("Executing command: %v", args)
//
//    // 执行 Redis 命令
//    cmdResult := cmd(ctx, client, args...)
//
//    // 在执行后记录日志
//    log.Printf("Command result: %v", cmdResult)
//
//    // 返回结果或错误
//    return cmdResult.String(), cmdResult.Err()
// }
