package sched

import (
	"fmt"
	"sync"
	"time"

	"agent/mem"
	"agent/process"
)

type Collector interface {
	Update() (interface{}, error)
}

var collectors []Collector = make([]Collector, 0)

func init() {
	collectors = append(collectors, &process.ProcessCollector{})
	collectors = append(collectors, &mem.MemoryCollector{})
}

func update() {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)

	for _, collector := range collectors {
		go func(c Collector) {
			defer wg.Done()
			res, err := c.Update()
			if err != nil {
				errChan <- err
				return
			}
			fmt.Println(res)
		}(collector)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Printf("Error occurred: %v\n", err)
	}

}

func Run() {
	// TODO need graceful shutdown
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			update()
		}
	}
}
