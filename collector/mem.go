package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/mem"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

func (m *MemoryCollector) GetMeterConfig() MeterConfig {

	metricName := "dd22"

	meter := otel.Meter("test-meter2")

	observable, err := meter.Int64ObservableCounter(
		metricName,
		metric.WithDescription("this is a test Counter"),
	)
	if err != nil {
		log.Fatalf("failed to create meter")
	}

	callback := func(ctx context.Context, observer metric.Observer) error {
		inc := int64(2)
		observer.ObserveInt64(observable, inc, metric.WithAttributes(attribute.String("endpoint", "/test2")))
		return nil
	}

	return MeterConfig{
		"memory", meter, observable, callback,
	}
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
