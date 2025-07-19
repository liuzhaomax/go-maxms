//go:build positive || all
// +build positive all

package test

import (
	"encoding/json"
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/test/common"
	. "github.com/onsi/ginkgo/v2"
	"net/http"
)

var _ = Describe("GET /login 获取公钥", func() {
	Context("200_success", func() {
		It("应该返回含有公钥字符串的json", func() {
			expectedRespBodyPath := "./base/get_login/response/200/200_success.json"
			expectedRespBody := common.ReadFile(expectedRespBodyPath)
			var expectedRespJson map[string]any
			_ = json.Unmarshal(expectedRespBody, &expectedRespJson)
			headers := common.BuildHttpHeaders(common.DefaultHeaders)
			expectedStatusCode := http.StatusOK
			appURL := fmt.Sprintf("%s://%s:%s", common.Cfg.Server.Protocol, common.Cfg.Server.Host, common.Cfg.Server.Port)
			common.RunTest(
				appURL,
				common.GetPukEndpoint,
				common.Get,
				core.EmptyString,
				headers,
				expectedStatusCode,
				expectedRespJson,
			)
		})
	})
})
