package core

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
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
	App
	Lib
	Server
	Downstream []Downstream
}

type App struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	PublicKeyStr string
	WhiteList    []WhiteList
}

type WhiteList struct {
	Name   string `mapstructure:"name"`
	Domain string `mapstructure:"domain"`
}

type Lib struct {
	Log
	Gin
}

type Log struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Color    bool   `mapstructure:"color"`
	Payload  bool   `mapstructure:"payload"`
	FileName string `mapstructure:"file_name"`
}

type Gin struct {
	RunMode string `mapstructure:"run_mode"`
}

type Server struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	BaseUrl         string `mapstructure:"base_url"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	IdleTimeout     int    `mapstructure:"idle_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

type Downstream struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

const configDir = "environment/config"

// 加载环境变量ENV，设置配置文件路径
func LoadEnv() string {
	// 读取环境变量 Mac和linux可以使用 export ENV=dev 直接设置环境变量，Windows要配环境变量并重启IDEA
	v := viper.New()
	v.AutomaticEnv()
	env := v.GetString("ENV")
	// 也可以通过添加flag “c”，执行命令行，来手动修改运行环境
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, env), "配置文件")
	flag.Parse()
	logrus.WithField("path", *configFile).Info("配置文件已识别")
	return *configFile
}

// 加载配置
func (cfg *Config) LoadConfig(configFile string) {
	v := viper.New()
	// 读取Config
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		logrus.WithField("path", configFile).WithField("失败方法", GetFuncName()).Panic("配置文件读取失败")
		panic(err)
	}
	err = v.Unmarshal(cfg)
	if err != nil {
		logrus.WithField("path", configFile).WithField("失败方法", GetFuncName()).Panic("配置文件反序列化失败")
		panic(err)
	}
	// 配置RSA密钥对
	cfg.SetRSAKeys()
}

// TODO 将公钥str存入ctx，将私钥存入vault
func (cfg *Config) SetRSAKeys() {
	//prk, puk, err := GenRSAKeyPair(2048)
	//if err != nil {
	//    logrus.WithField("失败方法", GetFuncName()).Panic("生成RSA密钥对失败")
	//    panic(err)
	//}
	//// TODO 存入vault
	//ctx.PublicKey = puk
	//ctx.PrivateKey = prk
	//publicKeyStr, err := PublicKeyToString()
	//if err != nil {
	//    logrus.WithField("失败方法", GetFuncName()).Panic("公钥转字符串失败")
	//    panic(err)
	//}
	//cfg.App.PublicKeyStr = publicKeyStr
}
