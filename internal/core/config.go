package core

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var cfg *Config
var once sync.Once

func init() {
	once.Do(func() {
		cfg = &Config{}
	})
}

func GetConfig() *Config {
	return cfg
}

type Config struct {
	CommonConfig
	AppConfig
}

type CommonConfig struct {
	App
}

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type AppConfig struct {
	Server
}

type Server struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

const configDir = "environment/config"
const commonConfigFile = "environment/config/common.yaml"

// 加载环境变量ENV，设置配置文件路径
func (cfg *Config) LoadEnv(v *viper.Viper) string {
	// 读取环境变量 Mac和linux可以使用 export ENV=dev 直接设置环境变量，Windows要配环境变量并重启IDEA
	v.AutomaticEnv()
	env := v.GetString("ENV")
	// 也可以通过添加flag “c”，执行命令行，来手动修改运行环境
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, env), "配置文件")
	flag.Parse()
	// TODO 日志输出读取的文件
	return *configFile
}

// 加载配置
func (cfg *Config) LoadConfig() {
	v := viper.New()
	// 读取AppConfig
	appConfigFile := cfg.LoadEnv(v)
	v.SetConfigFile(appConfigFile)
	err := v.ReadInConfig()
	if err != nil {
		panic(err) // TODO 日志修改
	}
	err = v.Unmarshal(&cfg.AppConfig)
	if err != nil {
		panic(err) // TODO 日志修改
	}
	// 读取CommonConfig
	v.SetConfigFile(commonConfigFile)
	err = v.ReadInConfig()
	if err != nil {
		panic(err) // TODO 日志修改
	}
	err = v.Unmarshal(&cfg.CommonConfig)
	if err != nil {
		panic(err) // TODO 日志修改
	}
}
