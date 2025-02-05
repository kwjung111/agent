package collector

import (
	"go.opentelemetry.io/otel/attribute"
)

func AttrUnitByte() attribute.KeyValue {
	return attribute.String("unit", "bb")
}

func AttrUnitPercent() attribute.KeyValue {
	return attribute.String("unit", "percent")
}
