package collector

import (
	"context"
)

type Collector interface {
	InitMeter() error
	GetName() string
}

var (
	collectors = make(map[string]Collector)
)

func InitCollectors(ctx context.Context) error {
	for _, collector := range collectors {
		err := collector.InitMeter()
		if err != nil {
			return err
		}
	}
	return nil
}

func Register(c Collector) {
	if c.GetName() == "" {
		panic("Collector name cannot be empty")
	}

	if collectors[c.GetName()] != nil {
		panic("Collector already exists with name: " + c.GetName())
	}

	collectors[c.GetName()] = c
}

func GetCollectors() map[string]Collector {
	return collectors
}
