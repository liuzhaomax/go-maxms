package config

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

type consulConfig struct {
	Timeout         int `mapstructure:"timeout"`
	Interval        int `mapstructure:"interval"`
	DeregisterAfter int `mapstructure:"deregister_after"`
	Endpoint        endpoint
}

// ServiceRegister 服务注册
func (c *consulConfig) ServiceRegister() error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf(
		"%s:%s",
		cfg.Lib.Consul.Endpoint.Host,
		cfg.Lib.Consul.Endpoint.Port,
	)

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}

	agentServiceRegistration := new(api.AgentServiceRegistration)
	agentServiceRegistration.Address = cfg.Server.Http.Host
	agentServiceRegistration.Name = cfg.App.Name
	agentServiceRegistration.ID = ext.ShortUUID()
	intPort, _ := strconv.Atoi(cfg.Server.Http.Port)
	agentServiceRegistration.Port = intPort
	agentServiceRegistration.Tags = []string{cfg.App.Name, cfg.Server.Http.Protocol}
	check := api.AgentServiceCheck{}
	check.Timeout = fmt.Sprintf("%ds", cfg.Lib.Consul.Timeout)
	check.Interval = fmt.Sprintf("%ds", cfg.Lib.Consul.Interval)
	check.DeregisterCriticalServiceAfter = fmt.Sprintf("%ds", cfg.Lib.Consul.DeregisterAfter)
	serverAddrHTTP := fmt.Sprintf("http://%s:%s/health", cfg.Server.Http.Host, cfg.Server.Http.Port)

	serverAddrRPC := fmt.Sprintf("%s:%s", cfg.Server.Http.Host, cfg.Server.Http.Port)

	switch cfg.Server.Http.Protocol {
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
func (c *consulConfig) ServiceDiscover() error {
	if len(cfg.Downstreams) == 0 || cfg.Downstreams == nil {
		return nil
	}

	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf(
		"%s:%s",
		cfg.Lib.Consul.Endpoint.Host,
		cfg.Lib.Consul.Endpoint.Port,
	)

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return err
	}

	for i, downstream := range cfg.Downstreams {
		services, _, err := client.Catalog().Service(downstream.Name, "", nil)
		if err != nil {
			return err
		}

		if len(services) == 0 {
			return fmt.Errorf(
				"未发现可用服务: %s: %s:%s",
				downstream.Name,
				downstream.Endpoint.Host,
				downstream.Endpoint.Port,
			)
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
