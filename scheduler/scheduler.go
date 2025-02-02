package scheduler

import (
	"agent/collector"
	"time"
)

const (
	interval = 1 * time.Second
)

func Run() {
	// TODO need graceful shutdown
	// TODO configurable interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		collector.UpdateAll()
	}
}
