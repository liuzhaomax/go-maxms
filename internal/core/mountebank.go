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
		log.Printf("imposter json 读取错误: %v", err)
		return
	}
	// 将读取的配置转换为 map[string]interface{}
	var imposterMap map[string]interface{}
	if err := v.Unmarshal(&imposterMap); err != nil {
		log.Printf("解析imposter为map失败: %v", err)
		return
	}
	stubPort := imposterMap["port"]
	// 序列化 map 为 JSON
	imposterJSON, err := json.Marshal(imposterMap)
	if err != nil {
		log.Printf("序列化imposter map为json失败: %v", err)
		return
	}
	// 请求的其他参数
	mbURL := fmt.Sprintf("%s://%s:%s/imposters", cfg.Lib.Mountebank.Protocol, cfg.Lib.Mountebank.Mb.Endpoint.Host, cfg.Lib.Mountebank.Mb.Endpoint.Port)
	contentType := "application/json"
	// 发送请求
	resp, err := http.Post(mbURL, contentType, bytes.NewBuffer(imposterJSON))
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Printf("创建imposter失败, status code: %v", resp.StatusCode)
		return
	}
	log.Printf("imposter创建成功，stub运行在端口：%v", stubPort)
}

func (m *Mountebank) DeleteImposter(stubDir string) {
	// 读取imposter请求体
	v := viper.New()
	v.SetConfigFile(stubDir)
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("imposter json 读取错误: %v", err)
		return
	}
	// 将读取的配置转换为 map[string]interface{}
	var imposterMap map[string]interface{}
	if err := v.Unmarshal(&imposterMap); err != nil {
		log.Printf("解析imposter为map失败: %v", err)
		return
	}
	stubPort := imposterMap["port"]
	// 请求的其他参数
	mbURL := fmt.Sprintf("%s://%s:%s/imposters", cfg.Lib.Mountebank.Protocol, cfg.Lib.Mountebank.Mb.Endpoint.Host, cfg.Lib.Mountebank.Mb.Endpoint.Port)
	// 发送请求
	req, err := http.NewRequest("DELETE", mbURL, nil)
	if err != nil {
		log.Printf("请求创建失败: %v", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("删除imposter失败, status code: %v", resp.StatusCode)
		return
	}
	log.Printf("位于%v端口的imposter删除成功\n", stubPort)
}
