# Go-MaxMs 中端模板

protobuf
```shell
# protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go
# grpc
go get -u google.golang.org/grpc
# protobuf
go get -u google.golang.org/protobuf
```

```shell
protoc -I . --go_out=plugins=grpc:. *.proto
```

wire
```shell
# 安装
go install github.com/google/wire/main/wire@latest
go get github.com/google/wire/main/wire@v0.5.0
# 生成
cd internal/app
go run github.com/google/wire/main/wire
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