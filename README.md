# go-maxms微服务脚手架

## 主要功能

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

## TODO 以后
1. vault与k8s集成，在登录后读取jwtsecret
2. vault与k8s集成，动态数据库账号密码  https://www.youtube.com/watch?v=otNkDHFNWt0
3. ~~vault使用production模式启动，链接换为https~~
4. 根据contract生成代码，包括type，不含集成其他工具的代码，（必填的非指针，可选的是指针），AST
5. redis账号密码登录
6. redis主从哨兵
7. redis TTL，持久化，布隆过滤器，雪崩击穿穿透
8. redis日志收集
9. RPC的中间件，包括token和签名
10. sentinel golang 不支持dashboard，需要二次开发：需建立与sentinel-dashboard通信的客户端（在config.Sentinel中定义地址），监听dashboard配置变化
11. prom监控 ELKf MQ jaeger consul vault
12. sgw负载均衡
