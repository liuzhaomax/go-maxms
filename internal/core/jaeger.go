package core

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	jConfig "github.com/uber/jaeger-client-go/config"
)

type Jaeger struct {
	Endpoint
}

func InitTracer() *jConfig.Configuration {
	return &jConfig.Configuration{
		Sampler: &jConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%s:%s", cfg.Jaeger.Endpoint.Host, cfg.Jaeger.Endpoint.Port),
		},
		ServiceName: cfg.App.Name,
	}
}
