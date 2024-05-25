package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type Mountebank struct {
	Protocol string
	Mb
	Imposter
}

type Mb struct {
	Endpoint
}

type Imposter struct {
	Endpoint
}

func (m *Mountebank) CreateImposter(stubDir string) {
	// 读取imposter请求体
	v := viper.New()
	v.SetConfigFile(stubDir)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("imposter json 读取错误: %v", err)
	}
	// 将读取的配置转换为 map[string]interface{}
	var imposterMap map[string]interface{}
	if err := v.Unmarshal(&imposterMap); err != nil {
		log.Fatalf("解析imposter为map失败: %v", err)
	}
	// 序列化 map 为 JSON
	imposterJSON, err := json.Marshal(imposterMap)
	if err != nil {
		log.Fatalf("序列化imposter map为json失败: %v", err)
	}
	// 请求的其他参数
	mbURL := fmt.Sprintf("%s://%s:%s/imposters", cfg.Lib.Mountebank.Protocol, cfg.Lib.Mountebank.Mb.Endpoint.Host, cfg.Lib.Mountebank.Mb.Endpoint.Port)
	contentType := "application/json"
	// 发送请求
	resp, err := http.Post(mbURL, contentType, bytes.NewBuffer(imposterJSON))
	if err != nil {
		log.Fatalf("发送imposter失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("创建imposter失败, status code: %v", resp.StatusCode)
	}
	log.Println("imposter创建成功")
}
