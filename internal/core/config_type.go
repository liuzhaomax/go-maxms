package core

import (
	"crypto/rsa"
	"github.com/sirupsen/logrus"
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
	Downstreams []Downstream
}

type App struct {
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
	Enabled       Enabled     `mapstructure:"enabled"`
	WhiteList     []WhiteList `mapstructure:"white_list"`
}

type Enabled struct {
	Vault            bool `mapstructure:"vault"`
	RSA              bool `mapstructure:"rsa"`
	Signature        bool `mapstructure:"signature"`
	RandomPort       bool `mapstructure:"random_port"`
	ServiceDiscovery bool `mapstructure:"service_discovery"`
}

type WhiteList struct {
	Name   string `mapstructure:"name"`
	Domain string `mapstructure:"domain"`
}

type Lib struct {
	Log
	Vault
	Gin
	DB
	Redis
	ETCD
	Consul
	Jaeger
	Rocketmq
}

type ETCD struct {
	DialTimeout          int `mapstructure:"dial_timeout"`
	DialKeepAliveTime    int `mapstructure:"dial_keep_alive_time"`
	DialKeepAliveTimeout int `mapstructure:"dial_keep_alive_timeout"`
	Endpoint
}

type Endpoint struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Server struct {
	Protocol        string `mapstructure:"protocol"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	BaseUrl         string `mapstructure:"base_url"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	IdleTimeout     int    `mapstructure:"idle_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

type Downstream struct {
	Id     string
	Secret string
	Name   string `mapstructure:"name"`
	Endpoint
}
