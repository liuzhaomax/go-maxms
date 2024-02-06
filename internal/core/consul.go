package core

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
)

type Consul struct {
	Timeout         string `mapstructure:"timeout"`
	Interval        string `mapstructure:"interval"`
	DeregisterAfter string `mapstructure:"deregister_after"`
	Endpoint
}

// ServiceRegister 服务注册
func (c *Consul) ServiceRegister() error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%s", cfg.Lib.Consul.Endpoint.Host, cfg.Lib.Consul.Endpoint.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	agentServiceRegistration := new(api.AgentServiceRegistration)
	agentServiceRegistration.Address = defaultConfig.Address
	agentServiceRegistration.Name = cfg.App.Name
	agentServiceRegistration.ID = ShortUUID()
	intPort, _ := strconv.Atoi(cfg.Server.Port)
	agentServiceRegistration.Port = intPort
	agentServiceRegistration.Tags = []string{cfg.App.Name, cfg.Server.Protocol}
	serverAddr := fmt.Sprintf("http://%s:%s/health", cfg.Server.Host, cfg.Server.Port)
	check := api.AgentServiceCheck{
		// GRPC: serverAddr,
		HTTP:                           serverAddr,
		Timeout:                        cfg.Lib.Consul.Timeout,
		Interval:                       cfg.Lib.Consul.Interval,
		DeregisterCriticalServiceAfter: cfg.Lib.Consul.DeregisterAfter,
	}
	agentServiceRegistration.Check = &check
	return client.Agent().ServiceRegister(agentServiceRegistration)
}
