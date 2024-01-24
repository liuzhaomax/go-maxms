package core

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
)

// ServiceRegister 服务注册
func ServiceRegister() error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%s", cfg.Lib.Consul.Endpoint.Host, cfg.Lib.Consul.Endpoint.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	agentServiceRegistration := new(api.AgentServiceRegistration)
	agentServiceRegistration.Address = defaultConfig.Address
	agentServiceRegistration.Name = cfg.App.Name
	agentServiceRegistration.ID = cfg.App.Name
	intPort, _ := strconv.Atoi(cfg.Server.Port)
	agentServiceRegistration.Port = intPort
	agentServiceRegistration.Tags = []string{cfg.App.Name, cfg.Server.Protocol}
	serverAddr := fmt.Sprintf("http://%s:%s/health", cfg.Server.Host, cfg.Server.Port)
	check := api.AgentServiceCheck{
		// GRPC: serverAddr,
		HTTP:                           serverAddr,
		Timeout:                        "3s",
		Interval:                       "30s",
		DeregisterCriticalServiceAfter: "5s",
	}
	agentServiceRegistration.Check = &check
	return client.Agent().ServiceRegister(agentServiceRegistration)
}
