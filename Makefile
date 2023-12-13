SHELL = /bin/bash

BuiltFile=bin/main
env=local
scenario=all

# 安装依赖
tidy:
	go mod tidy

# 打包
build:
	go build -o $(BuiltFile) -tags prod internal/cmd/main.go

# 运行
run:
	go run internal/cmd/main.go -c env/dev.yaml

# 语法检查
lint:
	golangci-lint run ./...

# 单元测试
unit:
	go test -v -race -timeout 1000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.out ./src/handler/... ./src/service/... ./src/utils/...

# API测试
api:
	go test -v -race -tags $(scenario) ./test/api -args -env=$(env)
