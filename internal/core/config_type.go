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
	Downstream []Downstream
}

type App struct {
	Name          string `mapstructure:"name"`
	Version       string `mapstructure:"version"`
	PublicKey     *rsa.PublicKey
	PrivateKey    *rsa.PrivateKey
	PublicKeyStr  string
	PrivateKeyStr string
	Salt          string
	JWTSecret     string
	Logger        *logrus.Logger
	Domain        string `mapstructure:"domain"`
	Enabled
	WhiteList []WhiteList
}

type Enabled struct {
	Vault bool `mapstructure:"vault"`
	RSA   bool `mapstructure:"rsa"`
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
}

type Log struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Color    bool   `mapstructure:"color"`
	FileName string `mapstructure:"file_name"`
}

type Vault struct {
	Address string `mapstructure:"address"`
}

type Gin struct {
	RunMode string `mapstructure:"run_mode"`
}

type DB struct {
	Type         string `mapstructure:"type"`
	Debug        bool   `mapstructure:"debug"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	Params       string `mapstructure:"params"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
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
