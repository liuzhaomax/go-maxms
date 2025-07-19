package config

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 加载环境变量ENV，设置配置文件路径
func LoadEnv() string {
	const configDir = "environment/config"
	// 读取环境变量 Mac和linux可以使用 export ENV=dev 直接设置环境变量，Windows要配环境变量并重启IDEA
	v := viper.New()
	v.AutomaticEnv()
	env := v.GetString("ENV")
	// 也可以通过添加flag “c”，执行命令行，来手动修改运行环境
	configFile := flag.String("c", fmt.Sprintf("%s/%s.yaml", configDir, env), "配置文件")
	flag.Parse()
	logrus.WithField("path", *configFile).Info(ext.FormatInfo("配置文件已识别"))

	return *configFile
}

// 加载配置
func (cfg *Config) LoadConfig(configFile string) {
	v := viper.New()
	// 读取Config
	v.SetConfigFile(configFile)

	err := v.ReadInConfig()
	if err != nil {
		logrus.WithField("path", configFile).
			WithField(FAILURE, ext.GetFuncName()).
			Panic(ext.FormatError(ext.ConfigError, "配置文件读取失败", err))
		panic(err)
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		logrus.WithField("path", configFile).
			WithField(FAILURE, ext.GetFuncName()).
			Panic(ext.FormatError(ext.ParseIssue, "配置文件反序列化失败", err))
		panic(err)
	}

	// 配置热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logrus.WithField("path", e.Name).Info("配置文件发生变化，重新加载配置")

		err = v.ReadInConfig()
		if err != nil {
			logrus.WithField("path", e.Name).
				WithField(FAILURE, ext.GetFuncName()).
				Error(ext.FormatError(ext.ConfigError, "重新加载配置文件失败", err))

			return
		}

		err = v.Unmarshal(cfg)
		if err != nil {
			logrus.WithField("path", e.Name).
				WithField(FAILURE, ext.GetFuncName()).
				Error(ext.FormatError(ext.ParseIssue, "配置文件反序列化失败", err))

			return
		}

		logrus.WithField("path", e.Name).Info("重新加载配置，成功")
	})

	// 处理加载的配置
	cfg.HandleLoadedConfig()
}

func (cfg *Config) HandleLoadedConfig() {
	// 配置日志
	InitLogger()
	// enabled几种情况（默认是第二种）
	// 1. 不使用RSA和vault：jwt_secret要自行设置（如下），salt会自动更新
	// 2. 不使用RSA，使用vault：jwt_secret，salt，puk，prk都要在vault预先设置好
	// 3. 使用RSA，不使用vault：jwt_secret要自行设置（如下），salt会自动更新，RSA自动存入内存
	// 4. 使用RSA，使用vault：RSA自动存入内存，并更新vault，salt读取于vault，jwt_secret要提前在vault设置好
	// AppId and AppSecret
	cfg.App.Id = ext.ShortUUID()

	cfg.App.Secret = ext.MD5Str(ext.ShortUUID())

	if cfg.App.Enabled.Vault {
		InitVault() // 配置Vault
		cfg.PutAppSecret()
	}

	if cfg.App.Enabled.RSA {
		// 生成密钥对，并将RSA结构体转为字符串，结构体与字符串都保存
		// 注意：这个RSA密钥对，是base64序列化后的pem block格式
		cfg.SetRSAKeys()
		// 写入secret
		if cfg.App.Enabled.Vault {
			cfg.PutRSA()
		}
	}

	if cfg.App.Enabled.Vault {
		go func() {
			for {
				// 包含RSA, JWT secret, Salt, DownstreamID和Secret
				cfg.GetSecret()
				// 将已保存的RSA字符串转为结构体，并保存
				cfg.ConvertRSAKeys()
				time.Sleep(time.Second * time.Duration(cfg.Lib.Vault.Interval))
			}
		}()
	} else {
		// 不适用vault需要自行设置jwt secret
		cfg.App.JWTSecret = "liuzhaomax@163.com"
	}
}

func (cfg *Config) SetRSAKeys() {
	prk, puk, err := ext.GenRSAKeyPair(2048)
	if err != nil {
		LogFailure(ext.Unknown, "生成RSA密钥对失败", err)
		panic(err)
	}

	cfg.App.PublicKey = puk
	cfg.App.PrivateKey = prk

	publicKeyStr, err := ext.PublicKeyToString(puk)
	if err != nil {
		LogFailure(ext.ParseIssue, "公钥转字符串失败", err)
		panic(err)
	}

	cfg.App.PublicKeyStr = publicKeyStr

	privateKeyStr, err := ext.PrivateKeyToString(prk)
	if err != nil {
		LogFailure(ext.ParseIssue, "私钥转字符串失败", err)
		panic(err)
	}

	cfg.App.PrivateKeyStr = privateKeyStr
}

func (cfg *Config) ConvertRSAKeys() {
	publicKey, err := ext.PublicKeyB64StrToStruct(cfg.App.PublicKeyStr)
	if err != nil {
		LogFailure(ext.ParseIssue, "公钥字符串转结构体失败", err)
		panic(err)
	}

	cfg.App.PublicKey = publicKey

	privateKey, err := ext.PrivateKeyB64StrToStruct(cfg.App.PrivateKeyStr)
	if err != nil {
		LogFailure(ext.ParseIssue, "私钥字符串转结构体失败", err)
		panic(err)
	}

	cfg.App.PrivateKey = privateKey
}

func GetPrivateKey() *rsa.PrivateKey {
	return cfg.App.PrivateKey
}

func GetPublicKey() *rsa.PublicKey {
	return cfg.App.PublicKey
}

func GetPublicKeyStr() string {
	return cfg.App.PublicKeyStr
}
