package reverse_proxy

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http/httputil"
	"time"
)

func Throttle(target string, c *gin.Context, proxy *httputil.ReverseProxy) {
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
	if blockError != nil {
		fmt.Println("流量过大，开启限流")
		time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
	} else {
		fmt.Println("正常通过")
		proxy.ServeHTTP(c.Writer, c.Request)
		time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
		entry.Exit()
	}
}
