package main

import (
	"fmt"
	"github.com/liuzhaomax/go-maxms-template-me/internal/core"
)

func main() {
	core.GetConfig().LoadConfig()
	fmt.Println(core.GetConfig())

	//ctx := context.Background()

}
