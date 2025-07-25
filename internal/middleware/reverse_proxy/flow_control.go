package reverse_proxy

import (
	"net/http"
	"net/http/httputil"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/gin-gonic/gin"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

func (rp *ReverseProxy) Throttle(target string, c *gin.Context, proxy *httputil.ReverseProxy) {
	configuration := config.NewDefaultConfig()
	configuration.Sentinel.Log.Logger = logging.NewConsoleLogger()

	err := sentinel.InitWithConfig(configuration)
	if err != nil {
		panic(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               target,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 匀速限流
			Threshold:              100,
			StatIntervalInMs:       1000, // 1000ms允许100个，QPS=100
			MaxQueueingTimeMs:      500,  // 500ms最大队列时长
			WarmUpPeriodSec:        30,   // 30s预热
		},
	})
	if err != nil {
		panic(err)
	}
	// 埋点
	entry, blockError := sentinel.Entry(target, sentinel.WithTrafficType(base.Inbound))
	if entry != nil {
		defer entry.Exit()
	}

	if blockError != nil {
		rp.AbortWithError(c, http.StatusForbidden, ext.Forbidden, "请求被限流", blockError)

		return
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
