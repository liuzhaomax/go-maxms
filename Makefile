SHELL = /bin/bash

CONTRACT_URL=https://raw.githubusercontent.com/liuzhaomax/maxblog-devops/main/contracts/user.yaml
CONTRACT_PATH=../spec/user.yaml
BUILT_FILE=bin/main
TEST_INCLUSION=$(shell go list ./... | grep -Ewv 'main|test|internal|src/router|src/dataAPI/handler')
API_ENV=local
SCENARIO=all


# 读取contract
spec:
	go test -v ./script -run TestGetContract -url=$(CONTRACT_URL) -path=$(CONTRACT_PATH)

# 打包
build:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILT_FILE) main/main.go

# 本地vendor
ven:
	go mod tidy
	go mod vendor

# 依赖注入
wire:
	wire ./internal/app

# 运行
run:
	go run main/main.go -c environment/config/dev.yaml

# 语法检查
# vendor确保lint不会启用下载，不然在ci过程中会timeout
# golangci-lint run -v -c ./.golangci.yml ./...
lint:
	gofmt -w -s .
	golangci-lint run -c ./.golangci.yml ./...

# 单元测试
unit:
	go test -v -timeout 1000s -covermode=atomic -coverpkg=./... -coverprofile=unit_test.out $(TEST_INCLUSION)
	go tool cover -html=unit_test.out

# 打测试桩
stub:
	go test ./test/imposter/imposter_test.go -run TestUpdateImposter

# 接口测试
api:
	go test -v -tags $(SCENARIO) ./test -args -env=$(ENV)

.PHONY: spec
.PHONY: stub
