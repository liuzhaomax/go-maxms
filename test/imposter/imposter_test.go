package imposter

import (
	"flag"
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/spf13/viper"
	"log"
	"testing"
)

const configDir = "../../environment/config"
const stubDir = "../stub/imposter.json"
const ENV = "dev"

func TestUpdateImposter(t *testing.T) {
	cfg := getConfig()
	// 删除imposter
	cfg.DeleteImposter(stubDir)
	// 创建imposter
	cfg.CreateImposter(stubDir)
}

func TestCreateImposter(t *testing.T) {
	cfg := getConfig()
	// 创建imposter
	cfg.CreateImposter(stubDir)
}

func TestDeleteImposter(t *testing.T) {
	cfg := getConfig()
	// 删除imposter
	cfg.DeleteImposter(stubDir)
}

func getConfig() *core.Config {
	v := viper.New()
	cfg := core.GetConfig()
	v.AutomaticEnv()
	env := flag.String("e", ENV, "环境")
	flag.Parse() // 后面有*env，必须先解析
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, *env), "配置文件")
	flag.Parse()
	// 读取Config
	v.SetConfigFile(*configFile)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("读取配置文件时出错: %v", err)
	}
	err = v.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("解析配置文件时出错: %v", err)
	}
	return cfg
}
