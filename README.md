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

mysql：
```shell
mysql -u root -p
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
github.com/lithammer/shortuuid
github.com/redis/go-redis/v9
go.etcd.io/etcd/client/v3
github.com/hashicorp/consul/api
google.golang.org/grpc/metadata
github.com/uber/jaeger-client-go
github.com/prometheus/client_golang/prometheus
github.com/alibaba/sentinel-golang/api
github.com/apache/rocketmq-client-go/v2

## TODO
1. ~~防抖与节流(完成-redis签名方式)，sentinel实现~~
2. ~~etcd注册中心，服务注册与服务发现，心跳检查，consul实现~~
3. ~~etcd配置中心，不实现~~
4. ~~Prometheus，grafana，jaeger，ELK，OCP4，熔断限流降级~~
5. ~~采用第一种鉴权方式，做SWG，修改auth中间件~~
6. ~~consul 服务发现随机端口~~
7. ~~SGW反向代理~~
8. SGW限流熔断降级 - 日志
9. SGW防爬虫
10. dtm
11. 消息队列 - 日志
12. SSO

## TODO 以后
1. vault与k8s集成，在登录后读取jwtsecret
2. vault与k8s集成，动态数据库账号密码  https://www.youtube.com/watch?v=otNkDHFNWt0
3. vault使用production模式启动，链接换为https
4. 根据contract生成代码，包括type，不含集成其他工具的代码，（必填的非指针，可选的是指针）
5. redis账号密码登录
6. redis主从哨兵
7. redis TTL，持久化，布隆过滤器，雪崩击穿穿透
8. redis日志收集
9. RPC的中间件，包括token和签名
10. sentinel golang 不支持dashboard，需要二次开发：需建立与sentinel-dashboard通信的客户端（在config.Sentinel中定义地址），监听dashboard配置变化
11. prom监控 ELKf MQ jaeger consul vault
12. sgw负载均衡
