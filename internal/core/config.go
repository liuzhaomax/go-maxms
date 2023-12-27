package core

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 加载环境变量ENV，设置配置文件路径
func LoadEnv() string {
	// 读取环境变量 Mac和linux可以使用 export ENV=dev 直接设置环境变量，Windows要配环境变量并重启IDEA
	v := viper.New()
	v.AutomaticEnv()
	env := v.GetString("ENV")
	// 也可以通过添加flag “c”，执行命令行，来手动修改运行环境
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, env), "配置文件")
	flag.Parse()
	logrus.WithField("path", *configFile).Info(FormatInfo("配置文件已识别"))
	return *configFile
}

// 加载配置
func (cfg *Config) LoadConfig(configFile string) func() {
	v := viper.New()
	// 读取Config
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		logrus.WithField("path", configFile).WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "配置文件读取失败", err))
		panic(err)
	}
	err = v.Unmarshal(cfg)
	if err != nil {
		logrus.WithField("path", configFile).WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "配置文件反序列化失败", err))
		panic(err)
	}
	// 配置日志
	cleanLogger := InitLogger()
	// 配置RSA密钥对
	cfg.SetRSAKeys()
	// 获取secret，包括JWT
	cfg.GetSecret()
	return cleanLogger
}

func (cfg *Config) SetRSAKeys() {
	prk, puk, err := GenRSAKeyPair(2048)
	if err != nil {
		logrus.WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "生成RSA密钥对失败", err))
		panic(err)
	}
	cfg.App.PublicKey = puk
	cfg.App.PrivateKey = prk
	publicKeyStr, err := PublicKeyToString(puk)
	if err != nil {
		logrus.WithField(FAILURE, GetFuncName()).Panic(FormatError(Unknown, "公钥转字符串失败", err))
		panic(err)
	}
	cfg.App.PublicKeyStr = publicKeyStr
}
