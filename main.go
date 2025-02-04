package main

import (
	"agent/collector"
	"agent/exporter"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	//init
	mp, err := exporter.InitMeterProvider(ctx)
	if err != nil {
		log.Fatalf("failed to initialize meterProvider : %v", err)
	}
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown meter provider : %v", err)
		}
	}()

	err = collector.InitCollectors(ctx)
	if err != nil {
		log.Fatalf("failed to initialize collectors : %v", err)
	}

	select {}
}
