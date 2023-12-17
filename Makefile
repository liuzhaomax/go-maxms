SHELL = /bin/bash

BuiltFile=bin/main
TestInclusion=$(shell go list ./... | grep -Ewv 'main|test|internal|src/router|src/dataAPI/handler')
env=local
scenario=all

# 安装依赖
tidy:
	go mod tidy

# 打包
build:
	go build -o $(BuiltFile) main/main.go

# 依赖注入
wire:
	wire ./internal/app

# 运行
run:
	go run main/main.go -c environment/config/dev.yaml

# 语法检查
lint:
	golangci-lint run ./...

# 单元测试
unit:
	go test -v -timeout 1000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.out $(TestInclusion)
	go tool cover -html=unit_test.out

# 接口测试
api:
	go test -v -race -tags $(scenario) ./test/api -args -env=$(env)
