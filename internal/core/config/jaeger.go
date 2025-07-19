package config

import (
	"fmt"

	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
)

type jaegerConfig struct {
	Endpoint endpoint
}

func InitTracer() *jConfig.Configuration {
	return &jConfig.Configuration{
		Sampler: &jConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jConfig.ReporterConfig{
			LogSpans: false,
			LocalAgentHostPort: fmt.Sprintf(
				"%s:%s",
				cfg.Lib.Jaeger.Endpoint.Host,
				cfg.Lib.Jaeger.Endpoint.Port,
			),
		},
		ServiceName: cfg.App.Name,
	}
}
