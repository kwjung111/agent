package collector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	name = "cpu"
)

type CpuCollector struct {
	Usages []float64
}

func init() {
	Register(&CpuCollector{})
}

func (c *CpuCollector) GetName() string {
	const metricName = "cpu"
	return metricName
}

func (c *CpuCollector) Update(ctx context.Context) error {
	cpuUsages, err := cpu.Percent(time.Second, true)
	if err != nil {
		fmt.Printf("err while get cpu : %v", err)
		return err
	}
	c.Usages = cpuUsages

	return nil
}

func (c *CpuCollector) InitMeter() error {
	meter := otel.Meter("os")

	c.Update(context.Background())

	var gauges []metric.Float64ObservableGauge

	for num, _ := range c.Usages {
		gauge, err := meter.Float64ObservableGauge(
			fmt.Sprintf("cpu%d_usage", num),
			metric.WithDescription("at est"),
		)
		if err != nil {
			log.Fatalf("failed to create meter")
		}
		gauges = append(gauges, gauge)
	}

	callback := func(ctx context.Context, observer metric.Observer) error {
		if err := c.Update(ctx); err != nil {
			return err
		}

		for num, gauge := range gauges {
			observer.ObserveFloat64(gauge, c.Usages[num], metric.WithAttributes(AttrUnitPercent()))
		}
		return nil
	}

	var instruments []metric.Observable
	for _, gauge := range gauges {
		instruments = append(instruments, gauge)
	}

	_, err := meter.RegisterCallback(callback, instruments...)
	if err != nil {
		log.Fatalf("failed to register Callback : %v", err)
	}

	return nil
}
