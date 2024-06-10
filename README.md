# go-maxms微服务脚手架

## 主要功能
+ 上手简单，中间件等工具可选
+ 支持多环境配置，配置热重载
+ 支持http/https/rpc协议，使用gin和grpc
+ 支持服务注册与发现，支持随机空闲端口部署
+ 支持vault密钥管理
+ 支持RSA加密解密，JWT生成，接口签名
+ 支持代理转发，熔断限流
+ 支持消息队列
+ 支持链路追踪
+ 具备监控指标接口
+ 具备日志系统
+ 具备自定义错误模式
+ 支持依赖注入
+ 具备统一的DB事务处理方式
+ 支持跨域访问
+ 支持分层架构
+ 支持CI/CD，jenkins流水线与docker自动化部署
+ 支持github工作流

## 使用指南
查看[使用指南](./init_common.md)了解更多详情。

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
13. ~~RPC服务集成prometheus~~
14. mysql集群
15. redis集群
16. 单元测试
17. ginkgo test
18. OCP4
19. grpcmock

## TODO 以后
1. vault与k8s集成，在登录后读取jwtsecret
2. vault与k8s集成，动态数据库账号密码  https://www.youtube.com/watch?v=otNkDHFNWt0
3. ~~vault使用production模式启动，链接换为https~~
4. 根据contract生成代码，包括type，不含集成其他工具的代码，（必填的非指针，可选的是指针），AST，需要注意是否使用指针类型
5. redis账号密码登录
6. redis主从哨兵
7. redis TTL，持久化，布隆过滤器，雪崩击穿穿透
8. ~~redis日志收集~~
9. ~~RPC的中间件，包括token和签名~~
10. sentinel golang 不支持dashboard，需要二次开发：需建立与sentinel-dashboard通信的客户端（在config.Sentinel中定义地址），监听dashboard配置变化
11. ~~prom监控 ELKf MQ jaeger consul vault~~
12. sgw负载均衡
13. mountebank 需要开一个服务，用读取list的方式，来显示当前运行的stub，提供一个接口（所要查询的端口是否被stub占用）
