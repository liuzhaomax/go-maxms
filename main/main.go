package main

import (
	"context"

	"github.com/liuzhaomax/go-maxms/internal/app"
	"github.com/liuzhaomax/go-maxms/internal/core/config"
)

func main() {
	app.Launch(
		context.Background(),
		app.SetConfigFile(config.LoadEnv()),
		app.SetWWWDir("www"),
	)
}
