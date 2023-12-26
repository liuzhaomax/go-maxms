# Go-MaxMs 模板

pb code gen
> Download the protoc bin file.
> https://github.com/protocolbuffers/protobuf/releases
```shell
# protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go
# grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

```shell
protoc -I . --go_out=plugins=grpc:. *.proto
```

wire
```shell
# 安装
go install github.com/google/wire/main/wire@latest
go get -u github.com/google/wire/main/wire@v0.5.0
# 生成
cd internal/app
go run github.com/google/wire/main/wire
# 生成
wire
```

代码覆盖率检查：
```shell
# 查看pwd下所有go文件代码覆盖率
go test -cover
# 查看pwd下所有go文件代码覆盖率，并输出覆盖率报告文件unit_test.out
go test -cover -coverprofile=unit_test.out
# 用html方式读取报告文件unit_test.out，可查看具体是哪段代码没有覆盖
go tool cover -html=unit_test.out
```

## 使用的需要安装的包
github.com/anaskhan96/go-password-encoder
github.com/google/wire/main/wire
github.com/gin-gonic/gin
github.com/sirupsen/logrus
github.com/snowzach/rotatefilehook
github.com/spf13/viper
github.com/mattn/go-colorable
github.com/golang/protobuf/protoc-gen-go
google.golang.org/grpc
google.golang.org/protobuf
github.com/hashicorp/vault/api
gorm.io/gorm

## TODO
1. 根据contract生成代码，包括type，不含集成其他工具的代码，必填的非指针，可选的是指针
2. 采用第一种鉴权方式，先做SWG -> main -> user，修改auth中间件
3. 动态数据库账号密码

## TODO 以后
1. vault与k8s集成，在登录后读取jwtsecret
2. vault与k8s集成，动态数据库账号密码  https://www.youtube.com/watch?v=otNkDHFNWt0
3. vault使用production模式启动，链接换为https
