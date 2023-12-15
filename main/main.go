package main

import (
	"context"
	"github.com/liuzhaomax/go-maxms-template-me/internal/app"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
)

func main() {
	app.Launch(
		context.Background(),
		app.SetConfigFile(core.LoadEnv()),
		app.SetWWWDir("www"),
	)

	core.InitLogger()

}
