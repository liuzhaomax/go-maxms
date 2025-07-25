package reverse_proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var ReverseProxySet = wire.NewSet(wire.Struct(new(ReverseProxy), "*"))

type ReverseProxy struct {
	Logger      *logrus.Entry
	RedisClient *redis.Client
}

// 防抖持续时间
const debounceDuration = 50 * time.Millisecond

// DebounceMiddleware 接口防抖
func (rp *ReverseProxy) DebounceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := "上次请求时间-" + clientIP
		// 检查上次请求时间
		lastRequest, err := rp.RedisClient.Get(context.Background(), key).Int64()
		if err != nil && !errors.Is(err, redis.Nil) {
			rp.AbortWithError(
				c,
				http.StatusInternalServerError,
				ext.Forbidden,
				"接口防抖查询redis错误",
				err,
			)

			return
		}
		// 对比上次请求时间到现在的时间，小于防抖时间，则报错429
		if err == nil && time.Since(time.Unix(lastRequest, 0)) < debounceDuration {
			rp.AbortWithError(c, http.StatusTooManyRequests, ext.Forbidden, "请求过于频繁接口防抖生效", err)

			return
		}
		// 记录当前请求时间
		rp.RedisClient.Set(context.Background(), key, time.Now().Unix(), debounceDuration)
		c.Next()
	}
}

// Redirect URL使用通配符，例如/api/*
func (rp *ReverseProxy) Redirect(serviceName string) gin.HandlerFunc {
	cfg := core.GetConfig()

	return func(c *gin.Context) {
		var addr string

		for _, downstream := range cfg.Downstreams {
			if downstream.Name == serviceName {
				addr = fmt.Sprintf(
					"http://%s:%s",
					downstream.Endpoint.Host,
					downstream.Endpoint.Port,
				)

				break
			}
		}

		proxyUrl, err := url.Parse(addr)
		if err != nil {
			rp.AbortWithError(c, http.StatusForbidden, ext.Forbidden, "反向代理URL解析错误", err)

			return
		}

		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

		rp.Logger.Info(ext.FormatInfo(fmt.Sprintf("反向代理到: %s, 地址: %s", serviceName, addr)))

		err = config.SetHeadersForDownstream(c, cfg.Downstreams[0].Name)
		if err != nil {
			rp.AbortWithError(c, http.StatusForbidden, ext.Forbidden, "反向代理请求头设置失败", err)

			return
		}

		proxy.ServeHTTP(c.Writer, c.Request)
		// rp.Throttle(target, c, proxy)  # 限流
		// rp.Break(target, c, proxy)     # 熔断
	}
}

// 使用
// root.GET("/login", mw.ReverseProxy.Redirect("maxblog-user"))

func (rp *ReverseProxy) AbortWithError(c *gin.Context, args ...any) {
	loggerFormat := config.GenGinLoggerFields(c)
	rp.Logger = rp.Logger.WithFields(loggerFormat)

	msg := &ext.MiddlewareMessage{
		StatusCode: 500,
		Code:       ext.InternalServerError,
		Desc:       "",
		Err:        nil,
	}

	switch len(args) {
	case 1: // 简化调用：AbortWithError(c, err)
		msg.Err = args[0].(error)
	case 4: // 复杂调用：AbortWithError(c, statusCode, code, desc, err)
		msg.StatusCode = args[0].(int)
		msg.Code = args[1].(ext.Code)
		msg.Desc = args[2].(string)
		msg.Err = args[3].(error)
	default:
		rp.Logger.Error("invalid arguments for AbortWithError")
	}

	formattedErr := ext.FormatError(msg.Code, msg.Desc, msg.Err)
	rp.Logger.Error(formattedErr)
	c.AbortWithStatusJSON(msg.StatusCode, ext.GenErrMsg(formattedErr))
}
