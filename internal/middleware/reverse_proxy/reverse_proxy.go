package reverse_proxy

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var ReverseProxySet = wire.NewSet(wire.Struct(new(ReverseProxy), "*"))

type ReverseProxy struct {
	Logger      core.ILogger
	RedisClient *redis.Client
}

// 防抖持续时间
const debounceDuration = 50 * time.Millisecond

// DebounceMiddleware 接口防抖
func (rp *ReverseProxy) DebounceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("上次请求时间-%s", clientIP)
		// 检查上次请求时间
		lastRequest, err := rp.RedisClient.Get(context.Background(), key).Int64()
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, rp.GenErrMsg(c, "接口防抖查询redis错误", err))
			return
		}
		// 对比上次请求时间到现在的时间，小于防抖时间，则报错429
		if err == nil && time.Since(time.Unix(lastRequest, 0)) < debounceDuration {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, rp.GenErrMsg(c, "请求过于频繁接口防抖生效", err))
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
				addr = fmt.Sprintf("http://%s:%s", downstream.Endpoint.Host, downstream.Endpoint.Port)
				break
			}
		}
		proxyUrl, err := url.Parse(addr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, rp.GenErrMsg(c, "反向代理URL解析错误", err))
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
		rp.GenOkMsg(c, fmt.Sprintf("反向代理到: %s, 地址: %s", serviceName, addr))
		err = core.SetHeadersForDownstream(c, cfg.Downstreams[0].Name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, rp.GenErrMsg(c, "反向代理请求头设置失败", err))
			return
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		// rp.Throttle(target, c, proxy)  # 限流
		// rp.Break(target, c, proxy)     # 熔断
	}
}

// 使用
// root.GET("/login", mw.ReverseProxy.Redirect("maxblog-user"))

func (rp *ReverseProxy) GenOkMsg(c *gin.Context, desc string) any {
	rp.Logger.SucceedWithField(c, desc)
	return gin.H{
		"status": gin.H{
			"code": core.OK,
			"desc": core.FormatInfo(desc),
		},
	}
}

func (rp *ReverseProxy) GenErrMsg(c *gin.Context, desc string, err error) any {
	rp.Logger.FailWithField(c, core.Forbidden, desc, err)
	return gin.H{
		"status": gin.H{
			"code": core.OK,
			"desc": core.FormatError(core.Forbidden, desc, err).Error(),
		},
	}
}
