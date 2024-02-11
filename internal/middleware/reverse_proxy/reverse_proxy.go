package reverse_proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var ReverseProxySet = wire.NewSet(wire.Struct(new(ReverseProxy), "*"))

type ReverseProxy struct {
	Logger      core.ILogger
	RedisClient *redis.Client
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
		err = core.SetHeadersForDownstream(c, cfg.Downstreams[0].Name, rp.RedisClient)
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

func (rp *ReverseProxy) GenOkMsg(c *gin.Context, desc string) string {
	rp.Logger.SucceedWithField(c, desc)
	return core.FormatInfo(desc)
}

func (rp *ReverseProxy) GenErrMsg(c *gin.Context, desc string, err error) error {
	rp.Logger.FailWithField(c, core.Forbidden, desc, err)
	return core.FormatError(core.Forbidden, desc, err)
}
