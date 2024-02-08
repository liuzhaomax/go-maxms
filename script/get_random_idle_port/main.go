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
	UpdateYamlConfig()
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

func UpdateYamlConfig() {
	v := viper.New()
	cfg := core.GetConfig()
	v.AutomaticEnv()
	env := flag.String("e", "dev", "环境")
	flag.Parse()
	fmt.Printf("%s/%s.yaml\n", configDir, *env)
	// 也可以通过添加flag “c”，执行命令行，来手动修改运行环境
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, *env), "配置文件")
	flag.Parse()
	fmt.Printf("%s/%s.yaml\n", configDir, *env)
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
	cfg.Server.Port = GetRandomIdlePort()
	fmt.Printf("现在的port是：%s\n", cfg.Server.Port)
	// 修改yaml文件
	if err = v.WriteConfig(); err != nil {
		log.Fatalf("写入配置文件时出错: %v", err)
	}
}
