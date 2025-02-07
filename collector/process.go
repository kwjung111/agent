package collector

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ProcessCollector struct {
}

type ProcessInfo struct {
	PID        int32
	Name       string
	CpuPercent float64
}

func init() {
	Register(&ProcessCollector{})
}

func getProcessInfo() ([]ProcessInfo, error) {

	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error fetching processes: %w", err)
	}

	var processInfos []ProcessInfo

	for _, p := range processes {
		cpuPercent, err := p.CPUPercent()
		if err != nil {
			continue
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CpuPercent: cpuPercent,
		})
	}
	return processInfos, nil
}

func (p *ProcessCollector) InitMeter() error {
	meter := otel.Meter("os")

	observable, err := meter.Float64ObservableGauge(
		"process",
		metric.WithDescription("this is a test Counter"),
	)
	if err != nil {
		log.Fatalf("failed to create meter")
	}

	callback := func(ctx context.Context, observer metric.Observer) error {
		data, err := p.Update(ctx)
		if err != nil {
			return err
		}
		processInfos, ok := data.([]ProcessInfo)
		if !ok {
			return fmt.Errorf("invalid process data")
		}

		for _, proc := range processInfos {
			observer.ObserveFloat64(
				observable,
				proc.CpuPercent,
				metric.WithAttributes(
					attribute.String("pid", fmt.Sprintf("%d", proc.PID)),
					attribute.String("name", proc.Name),
					AttrUnitPercent(),
				),
			)
		}
		return nil
	}

	_, err = meter.RegisterCallback(callback, observable)
	if err != nil {
		log.Fatalf("error!")
	}

	return nil
}

func (p *ProcessCollector) GetName() string {
	const metricName = "process"
	return metricName
}

func (p *ProcessCollector) Update(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		pInfo, err := getProcessInfo()
		if err != nil {
			return nil, err
		}
		top := getTopProcesses(pInfo)

		return top, nil
	}
}

// getTopProcesses 함수: CPU 사용률이 높은 상위 5개 프로세스를 정렬 및 출력
func getTopProcesses(pInfo []ProcessInfo) []ProcessInfo {
	sort.Slice(pInfo, func(i, j int) bool {
		return pInfo[i].CpuPercent > pInfo[j].CpuPercent
	})

	topN := min(len(pInfo), 5)

	return pInfo[:topN]
}

// min 함수: 두 숫자 중 작은 값을 반환
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
