package reverse_proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
)

var ReverseProxySet = wire.NewSet(wire.Struct(new(ReverseProxy), "*"))

type ReverseProxy struct {
	Logger core.ILogger
}

func ReverseProxyRedirect(target *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = path.Join(target.Path, "/api", c.Param("action"))
		}
		proxy := &httputil.ReverseProxy{
			Director: director,
		}
		// TODO 加入authorisation SGW
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// 使用
//proxyUrl, _ := url.Parse("http://127.0.0.1:8080")
//r.GET("/api/:action", ReverseProxyRedirect(proxyUrl))
