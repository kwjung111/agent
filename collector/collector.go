package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Collector interface {
	GetMeterConfig() MeterConfig
	GetName() string
	Update(ctx context.Context) (interface{}, error)
}

var (
	collectors = make(map[string]Collector)
)

func InitCollectors(ctx context.Context) error {
	for _, collector := range collectors {
		cfg := collector.GetMeterConfig()
		//init
		_, err := cfg.meter.RegisterCallback(cfg.callback, cfg.observable)
		if err != nil {
			log.Fatalf("error!")
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

func UpdateAll() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var clen = len(collectors)

	errChan := make(chan error, clen)
	wg.Add(clen)

	for _, collector := range collectors {
		go func(c Collector) {
			defer wg.Done()
			res, err := c.Update(ctx)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Println(res) // FOR DEBUGGING
		}(collector)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Printf("Error occurred: %v\n", err)
	}

}
