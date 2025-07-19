package core

import (
	"github.com/liuzhaomax/go-maxms/internal/core/config"
)

func GetConfig() *config.Config {
	return config.GetConfig()
}
