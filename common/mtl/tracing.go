package mtl

import (
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	//"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func InitTracing(serviceName string) provider.OtelProvider {

	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		//	provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
		provider.WithEnableMetrics(false),
	)
	//defer p.Shutdown(context.Background())
	return p
}
