package common

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
)

const configDir = "../environment/config"

var Cfg *core.Config

func LoadConfig(env string) {
	v := viper.New()
	Cfg = core.GetConfig()
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, env), "配置文件")
	flag.Parse()
	// 读取Config
	v.SetConfigFile(*configFile)
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("读取配置文件时出错: %v", err)
		return
	}
	err = v.Unmarshal(Cfg)
	if err != nil {
		log.Printf("解析配置文件时出错: %v", err)
		return
	}
}

func BuildHttpHeaders(headerInfo map[string]string) http.Header {
	newHeader := http.Header{}
	for key, value := range headerInfo {
		newHeader.Add(key, value)
	}
	return newHeader
}

func BuildHttpRequest(requestMethod string, appURL string, endpoint string, headers http.Header, rawRequestBody string) (*http.Request, error) {
	requestURL := appURL + endpoint
	var requestBody *bytes.Buffer
	if isJson(rawRequestBody) {
		requestBody = bytes.NewBuffer([]byte(rawRequestBody))
	} else {
		requestBody = bytes.NewBuffer(ReadFile(rawRequestBody))
	}
	req, err := http.NewRequest(requestMethod, requestURL, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return req, nil
}

func isJson(s string) bool {
	var js map[string]any
	return json.Unmarshal([]byte(s), &js) == nil
}

func ReadFile(requestFile string) []byte {
	reqBodyByte, err := os.ReadFile(requestFile)
	if err != nil {
		return nil
	}
	return reqBodyByte
}

func MakeHttpRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func MustReadBody(response *http.Response) []byte {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}
