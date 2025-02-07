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
	Free      uint64
	Available uint64
	Used      uint64
}

func init() {
	Register(&MemoryCollector{})
}

func getMemInfo() (MemoryInfo, error) {
	virtMem, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	ret, err := parse(virtMem)
	if err != nil {
		return MemoryInfo{}, err
	}
	return ret, nil
}

func parse(virtMem *mem.VirtualMemoryStat) (MemoryInfo, error) {
	total := virtMem.Total
	free := virtMem.Free
	available := virtMem.Available
	used := virtMem.Used

	memoryInfo := MemoryInfo{
		Total:     total,
		Free:      free,
		Available: available,
		Used:      used,
	}

	return memoryInfo, nil
}

func (m *MemoryCollector) InitMeter() error {

	meter := otel.Meter("os")

	total, err := meter.Int64ObservableGauge(
		"total_memory",
		metric.WithDescription("total_memory"),
	)
	if err != nil {
		log.Fatalf("failed to create meter : total memory : %v", err)
	}

	free, err := meter.Int64ObservableGauge(
		"free_memory",
		metric.WithDescription("free_memory"),
	)
	if err != nil {
		log.Fatalf("failed to create meter : free memory : %v", err)
	}

	used, err := meter.Int64ObservableGauge(
		"used_memory",
		metric.WithDescription("used_memory"),
	)
	if err != nil {
		log.Fatalf("failed to create meter : used memory : %v", err)
	}

	available, err := meter.Int64ObservableGauge(
		"available_memory",
		metric.WithDescription("available_memory"),
	)
	if err != nil {
		log.Fatalf("failed to create meter : available memory : %v", err)
	}

	callback := func(ctx context.Context, observer metric.Observer) error {
		memory, err := m.Update(ctx)
		if err != nil {
			return fmt.Errorf("invalid memory data")
		}

		observer.ObserveInt64(total, int64(memory.Total), metric.WithAttributes(AttrUnitByte()))
		observer.ObserveInt64(free, int64(memory.Free), metric.WithAttributes(AttrUnitByte()))
		observer.ObserveInt64(used, int64(memory.Used), metric.WithAttributes(AttrUnitByte()))
		observer.ObserveInt64(available, int64(memory.Available), metric.WithAttributes(AttrUnitByte()))
		return nil
	}

	_, err = meter.RegisterCallback(callback, total, free, used, available)
	if err != nil {
		log.Fatalf("error!")
	}

	return nil
}

func (m *MemoryCollector) GetName() string {
	const metricName = "memory"
	return metricName
}

func (m *MemoryCollector) Update(ctx context.Context) (MemoryInfo, error) {
	select {
	case <-ctx.Done():
		return MemoryInfo{}, ctx.Err()
	default:
		mInfo, err := getMemInfo()
		if err != nil {
			return MemoryInfo{}, nil
		}
		return mInfo, nil
	}
}
