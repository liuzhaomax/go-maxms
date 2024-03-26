package core

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/google/wire"
	"time"
)

type Rocketmq struct {
	Timeout int `mapstructure:"timeout"`
	Retry   int `mapstructure:"retry"`
	Endpoint
}

var RocketMQSet = wire.NewSet(wire.Struct(new(RocketMQ), "*"), wire.Bind(new(IRocketMQ), new(*RocketMQ)))

type IRocketMQ interface {
	GenProducer(string) (*rocketmq.Producer, error)
	GenPushConsumer(string) (*rocketmq.PushConsumer, error)
}

type RocketMQ struct {
}

func (r *RocketMQ) GenProducer(groupName string) (*rocketmq.Producer, error) {
	newProducer, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{fmt.Sprintf("%s:%s", cfg.Lib.Rocketmq.Endpoint.Host, cfg.Lib.Rocketmq.Endpoint.Port)}),
		producer.WithRetry(cfg.Lib.Rocketmq.Retry), // 尝试发送数据的次数
		producer.WithSendMsgTimeout(time.Second*time.Duration(cfg.Lib.Rocketmq.Timeout)),
		producer.WithGroupName(groupName),
	)
	if err != nil {
		return nil, err
	}
	return &newProducer, nil
}

func (r *RocketMQ) GenPushConsumer(groupName string) (*rocketmq.PushConsumer, error) {
	newConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%s", cfg.Lib.Rocketmq.Endpoint.Host, cfg.Lib.Rocketmq.Endpoint.Port)}),
		consumer.WithRetry(cfg.Lib.Rocketmq.Retry), // 尝试发送数据的次数
		consumer.WithConsumeTimeout(time.Second*time.Duration(cfg.Lib.Rocketmq.Timeout)),
		consumer.WithGroupName(groupName),
	)
	if err != nil {
		return nil, err
	}
	return &newConsumer, nil
}
