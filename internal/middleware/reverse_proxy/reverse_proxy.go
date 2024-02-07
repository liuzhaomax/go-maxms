package reverse_proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var ReverseProxySet = wire.NewSet(wire.Struct(new(ReverseProxy), "*"))

type ReverseProxy struct {
	Logger core.ILogger
}

// Redirect URL使用通配符，例如/api/*
func (rp *ReverseProxy) Redirect(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxyUrl, err := url.Parse(target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, rp.GenErrMsg(c, "反向代理URL解析错误", err))
		}
		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
		rp.GenOkMsg(c, fmt.Sprintf("反向代理到地址: %s", target))
		proxy.ServeHTTP(c.Writer, c.Request)
		// rp.Throttle(target, c, proxy)  # 限流
		// rp.Break(target, c, proxy)     # 熔断
	}
}

// 使用
// root.GET("/login", mw.ReverseProxy.Redirect("http://172.30.64.1:9999"))

func (rp *ReverseProxy) GenOkMsg(c *gin.Context, desc string) string {
	rp.Logger.SucceedWithField(c, desc)
	return core.FormatInfo(desc)
}

func (rp *ReverseProxy) GenErrMsg(c *gin.Context, desc string, err error) error {
	rp.Logger.FailWithField(c, core.Forbidden, desc, err)
	return core.FormatError(core.Forbidden, desc, err)
}
