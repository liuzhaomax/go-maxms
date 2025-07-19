package imposter_test

import (
	"flag"
	"fmt"
	"log"
	"testing"

	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/spf13/viper"
)

const (
	configDir = "../../environment/config"
	stubDir   = "../stub/imposter.json"
	ENV       = "dev"
)

func TestUpdateImposter(t *testing.T) {
	cfg := getConfig()
	// 删除imposter
	cfg.Lib.Mountebank.DeleteImposter(stubDir)
	// 创建imposter
	cfg.Lib.Mountebank.CreateImposter(stubDir)
}

func TestCreateImposter(t *testing.T) {
	cfg := getConfig()
	// 创建imposter
	cfg.Lib.Mountebank.CreateImposter(stubDir)
}

func TestDeleteImposter(t *testing.T) {
	cfg := getConfig()
	// 删除imposter
	cfg.Lib.Mountebank.DeleteImposter(stubDir)
}

func getConfig() *config.Config {
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
		log.Printf("读取配置文件时出错: %v", err)

		return nil
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		log.Printf("解析配置文件时出错: %v", err)

		return nil
	}

	return cfg
}
