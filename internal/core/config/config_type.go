package config

import (
	"crypto/rsa"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	cfg  *Config
	once sync.Once
)

func init() {
	once.Do(func() {
		cfg = &Config{}
	})
}

func GetConfig() *Config {
	return cfg
}

type Config struct {
	App         app
	Lib         lib
	Server      server
	Secret      secret
	Downstreams []downstream
}

type app struct {
	Id            string
	Secret        string
	Name          string `mapstructure:"name"`
	Version       string `mapstructure:"version"`
	Domain        string `mapstructure:"domain"`
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyStr  string
	PrivateKeyStr string
	Salt          string
	JWTSecret     string
	Logger        *logrus.Logger
	Enabled       enabled     `mapstructure:"enabled"`
	WhiteList     []whiteList `mapstructure:"white_list"`
}

type enabled struct {
	Vault            bool `mapstructure:"vault"`
	RSA              bool `mapstructure:"rsa"`
	Signature        bool `mapstructure:"signature"`
	HeaderParams     bool `mapstructure:"header_params"`
	RandomPort       bool `mapstructure:"random_port"`
	ServiceDiscovery bool `mapstructure:"service_discovery"`
	HealthCheck      bool `mapstructure:"health_check"`
	Prometheus       bool `mapstructure:"prometheus"`
	Jaeger           bool `mapstructure:"jaeger"`
}

type whiteList struct {
	Name   string `mapstructure:"name"`
	Domain string `mapstructure:"domain"`
}

type lib struct {
	Log        logConfig
	Vault      vaultConfig
	Gin        ginConfig
	DB         dbConfig
	Redis      redisConfig
	WebSocket  webSocketConfig
	ETCD       etcdConfig
	Consul     consulConfig
	Jaeger     jaegerConfig
	Rocketmq   rocketmqConfig
	Mountebank mountebankConfig
}

type server struct {
	Protocol        string `mapstructure:"protocol"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	BaseUrl         string `mapstructure:"base_url"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	IdleTimeout     int    `mapstructure:"idle_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

type secret struct {
	Mysql  mysqlSecret
	Redis  redisSecret
	Wechat wechatSecret
}

type wechatSecret struct {
	AppId     string
	AppSecret string
}

type redisSecret struct {
	Password string
}

type mysqlSecret struct {
	Name     string
	UserName string
	PassWord string
}

type downstream struct {
	Id       string
	Secret   string
	Name     string `mapstructure:"name"`
	Endpoint endpoint
}

type endpoint struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
