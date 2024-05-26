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
const DEV = "dev"

func TestCreateImposter(t *testing.T) {
	v := viper.New()
	cfg := core.GetConfig()
	v.AutomaticEnv()
	env := flag.String("e", DEV, "环境")
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

	cfg.CreateImposter(stubDir)
}