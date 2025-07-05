package main

import (
	"context"
	"github.com/liuzhaomax/go-maxms/internal/app"
	"github.com/liuzhaomax/go-maxms/internal/core"
)

func main() {
	app.Launch(
		context.Background(),
		app.SetConfigFile(core.LoadEnv()),
		app.SetWWWDir("www"),
	)
}
