package common

import (
	"encoding/json"
	"github.com/liuzhaomax/go-maxms/internal/core"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
)

const (
	Get            = "GET"
	GetPukEndpoint = "/login"
)

var DefaultHeaders = map[string]string{
	"App_id":     "prismer@ai",
	"Request_id": "123",
	"Span_id":    "123",
	"Trace_id":   "123",
	"Parent_id":  "123",
	"Client_ip":  "123.123.123.123",
	"User-Agent": "postman",
	"Signature":  core.GenAppSignature("prismer@ai", "prismer@ai", "", "123"),
}

func RunTest(
	appURL, endpoint, requestMethod, requestBody string,
	headers http.Header,
	expectedStatusCode int,
	expectedResponse map[string]any,
) {
	request, err := BuildHttpRequest(requestMethod, appURL, endpoint, headers, requestBody)
	if err != nil {
		log.Fatal("API测试建立请求失败: ", err)
	}
	httpResponse, err := MakeHttpRequest(request)
	if err != nil {
		log.Fatal("API测试响应错误: ", err)
	}
	defer httpResponse.Body.Close()
	responseBody := MustReadBody(httpResponse)
	Expect(httpResponse.StatusCode).Should(Equal(expectedStatusCode))
	var actualResponse map[string]any
	err = json.Unmarshal(responseBody, &actualResponse)
	if err != nil {
		log.Println("API测试实际响应体解析错误: ", err)
		return
	}
	Expect(actualResponse).To(HaveLen(len(expectedResponse)))
	Expect(actualResponse).Should(Equal(expectedResponse)) // 底层是reflect.DeepEqual
}
