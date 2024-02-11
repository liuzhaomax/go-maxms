package core

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
)

type Consul struct {
	Timeout         int `mapstructure:"timeout"`
	Interval        int `mapstructure:"interval"`
	DeregisterAfter int `mapstructure:"deregister_after"`
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
	agentServiceRegistration.Address = cfg.Server.Host
	agentServiceRegistration.Name = cfg.App.Name
	agentServiceRegistration.ID = ShortUUID()
	intPort, _ := strconv.Atoi(cfg.Server.Port)
	agentServiceRegistration.Port = intPort
	agentServiceRegistration.Tags = []string{cfg.App.Name, cfg.Server.Protocol}
	serverAddr := fmt.Sprintf("http://%s:%s/health", cfg.Server.Host, cfg.Server.Port)
	check := api.AgentServiceCheck{
		// GRPC: serverAddr,
		HTTP:                           serverAddr,
		Timeout:                        fmt.Sprintf("%ds", cfg.Lib.Consul.Timeout),
		Interval:                       fmt.Sprintf("%ds", cfg.Lib.Consul.Interval),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", cfg.Lib.Consul.DeregisterAfter),
	}
	agentServiceRegistration.Check = &check
	return client.Agent().ServiceRegister(agentServiceRegistration)
}

// ServiceDiscover 服务发现
func (c *Consul) ServiceDiscover() error {
	if cfg.Downstreams == nil || len(cfg.Downstreams) == 0 {
		return nil
	}
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%s", cfg.Lib.Consul.Endpoint.Host, cfg.Lib.Consul.Endpoint.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}
	for i, downstream := range cfg.Downstreams {
		services, _, err := client.Catalog().Service(downstream.Name, cfg.Server.Protocol, nil)
		if err != nil {
			return err
		}
		if len(services) == 0 {
			return errors.New("未发现可用服务: " + downstream.Name)
		}
		for _, service := range services {
			if downstream.Name == service.ServiceName {
				cfg.Downstreams[i].Endpoint.Host = service.ServiceAddress
				cfg.Downstreams[i].Endpoint.Port = strconv.Itoa(service.ServicePort)
			}
		}
	}
	return nil
}
