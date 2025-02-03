package main

import (
	"agent/exporter"
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	initMeter()
	//scheduler.Run()
}

func initMeter() {
	ctx := context.Background()

	//TODO structuring / decoupling

	mp, err := exporter.InitMeterProvider(ctx)
	if err != nil {
		log.Fatalf("failed to initialized meter propvider : %v", err)
	}
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown meter provider : %v", err)
		}
	}()

	meter := otel.Meter("test-meter")

	testCounter, err := meter.Int64ObservableCounter(
		"test_counter",
		metric.WithDescription("this is a test Counter"),
	)
	if err != nil {
		log.Fatalf("failed to create meter")
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			inc := int64(1)
			observer.ObserveInt64(testCounter, inc, metric.WithAttributes(attribute.String("endpoint", "/example")))
			return nil
		},
		testCounter,
	)
	if err != nil {
		log.Fatalf("error!")
	}

	select {}
}
