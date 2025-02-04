package collector

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type MeterConfig struct {
	metricName string
	meter      metric.Meter
	observable metric.Observable
	callback   func(context.Context, metric.Observer) error
}
