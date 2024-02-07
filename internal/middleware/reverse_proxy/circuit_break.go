package reverse_proxy

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	// fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	// fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	// fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (rp *ReverseProxy) Break(target string, c *gin.Context, proxy *httputil.ReverseProxy) {
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})
	_, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:                     target,
			Strategy:                     circuitbreaker.SlowRequestRatio, // 慢调用熔断策略
			RetryTimeoutMs:               3000,                            // 熔断触发后持续3秒
			MinRequestAmount:             10,                              // 周期内，触发熔断的最小请求数目，请求数小于此数值，即使到达熔断条件，也不会触发熔断
			StatIntervalMs:               1000,                            // 窗口长度
			StatSlidingWindowBucketCount: 10,                              // 随着桶数的增加，统计数据会更加精确，但内存成本也会增加，“StatIntervalMs % StatSlidingWindowBucketCount == 0”，否则StatSlidingWindowBucketCount将被1取代
			Threshold:                    0.5,                             // 触发慢调用熔断比例
			MaxAllowedRtMs:               100,                             // 慢调用判断条件 > 100ms的RT
		},
		// 测试数据
		// {
		// 	Resource:                     target,
		// 	Strategy:                     circuitbreaker.SlowRequestRatio, // 慢调用熔断策略
		// 	RetryTimeoutMs:               3000,                            // 熔断触发后持续3秒
		// 	MinRequestAmount:             0,                               // 周期内，触发熔断的最小请求数目，请求数小于此数值，即使到达熔断条件，也不会触发熔断
		// 	StatIntervalMs:               2,                               // 窗口长度
		// 	StatSlidingWindowBucketCount: 10,                              // 随着桶数的增加，统计数据会更加精确，但内存成本也会增加，“StatIntervalMs % StatSlidingWindowBucketCount == 0”，否则StatSlidingWindowBucketCount将被1取代
		// 	Threshold:                    0.01,                            // 触发慢调用熔断比例
		// 	MaxAllowedRtMs:               1,                               // 慢调用判断条件 > 100ms的RT
		// },
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
		c.AbortWithStatusJSON(http.StatusForbidden, rp.GenErrMsg(c, "服务被熔断", err))
		return
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
