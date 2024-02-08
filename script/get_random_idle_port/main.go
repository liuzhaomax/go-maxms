package main

import (
	"flag"
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/spf13/viper"
	"log"
	"net"
	"strconv"
)

const configDir = "environment/config"

func main() {
	fmt.Print(UpdateYamlConfig())
}

func GetRandomIdlePort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port)
}

func UpdateYamlConfig() string {
	v := viper.New()
	cfg := core.GetConfig()
	v.AutomaticEnv()
	env := flag.String("e", "dev", "环境")
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
	// 修改port
	if cfg.App.Enabled.RandomPort {
		cfg.Server.Port = GetRandomIdlePort()
	}
	// 修改yaml文件
	v.Set("server.port", cfg.Server.Port)
	if err = v.WriteConfig(); err != nil {
		log.Fatalf("写入配置文件时出错: %v", err)
	}
	return cfg.Server.Port
}
