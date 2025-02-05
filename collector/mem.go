package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/mem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type MemoryCollector struct{}

type MemoryInfo struct {
	Total     uint64
	Available uint64
	Used      uint64
}

func init() {
	Register(&MemoryCollector{})
}

func getMem() (MemoryInfo, error) {
	virtMem, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	fmt.Println("")
	ret, err := parse(virtMem)
	if err != nil {
		return MemoryInfo{}, err
	}
	return ret, nil
}

func parse(virtMem *mem.VirtualMemoryStat) (MemoryInfo, error) {
	total := virtMem.Total
	available := virtMem.Available
	used := virtMem.Used

	memoryInfo := MemoryInfo{
		Total:     total,
		Available: available,
		Used:      used,
	}

	return memoryInfo, nil
}

func (m *MemoryCollector) InitMeter() error {

	meter := otel.Meter("os")

	gauge1, err := meter.Int64ObservableGauge(
		"gauge1",
		metric.WithDescription("this is a test Counter"),
	)
	if err != nil {
		log.Fatalf("failed to create meter")
	}

	gauge2, err := meter.Float64ObservableGauge(
		"gauge2",
		metric.WithDescription("this is guage2"),
	)
	if err != nil {
		log.Fatalf("failed to create meter")
	}

	callback := func(ctx context.Context, observer metric.Observer) error {
		observer.ObserveInt64(gauge1, 11, metric.WithAttributes(AttrUnitByte()))
		observer.ObserveFloat64(gauge2, 123.45, metric.WithAttributes(AttrUnitPercent()))
		return nil
	}

	_, err = meter.RegisterCallback(callback, gauge1, gauge2)
	if err != nil {
		log.Fatalf("error!")
	}

	return nil
}

func (m *MemoryCollector) GetName() string {
	const metricName = "memory"
	return metricName
}

func (m *MemoryCollector) Update(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		mem, err := getMem()
		if err != nil {
			return MemoryInfo{}, nil
		}
		return mem, nil
	}
}
