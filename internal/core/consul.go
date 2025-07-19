package core

import (
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
	check := api.AgentServiceCheck{}
	check.Timeout = fmt.Sprintf("%ds", cfg.Lib.Consul.Timeout)
	check.Interval = fmt.Sprintf("%ds", cfg.Lib.Consul.Interval)
	check.DeregisterCriticalServiceAfter = fmt.Sprintf("%ds", cfg.Lib.Consul.DeregisterAfter)
	serverAddrHTTP := fmt.Sprintf("http://%s:%s/health", cfg.Server.Host, cfg.Server.Port)
	serverAddrRPC := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	switch cfg.Server.Protocol {
	case "http":
		check.HTTP = serverAddrHTTP
	case "rpc":
		check.GRPC = serverAddrRPC
	default:
		check.HTTP = serverAddrHTTP
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
		services, _, err := client.Catalog().Service(downstream.Name, EmptyString, nil)
		if err != nil {
			return err
		}
		if len(services) == 0 {
			return fmt.Errorf("未发现可用服务: %s: %s:%s", downstream.Name, downstream.Host, downstream.Port)
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
