package exporter

import (
	config "agent/common"
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	meterProvider *sdkmetric.MeterProvider
	once          sync.Once
	initErr       error
)

func InitMeterProvider(ctx context.Context) (*sdkmetric.MeterProvider, error) {

	config := config.GetConfig()
	interval := config.GetScrapeInterval()

	once.Do(func() {

		// Create a grpc Exporer
		exporter, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint("localhost:4317"),
			otlpmetricgrpc.WithInsecure())
		if err != nil {
			return
		}

		//TODO CONFIG
		res, err := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String("test"),
				attribute.String("environment", "production"),
			),
		)
		if err != nil {
			return
		}

		// Create a trace provider with the exporter and a resource
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(interval*time.Second))),
		)

		// Register the trace provider with the global trace provider
		otel.SetMeterProvider(meterProvider)
	})

	fmt.Println("meterProvider inited")
	return meterProvider, initErr
}

func GetMeterProvider() *sdkmetric.MeterProvider {
	return meterProvider
}
