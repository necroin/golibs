package metrics

import (
	"io"
	"strings"

	"github.com/necroin/golibs/concurrent"
)

type Labels map[string]string

type Metric interface {
	Description() *Description
	Write(io.Writer)
}

type MetricVector[T Metric] struct {
	Metric
	data                *concurrent.ConcurrentMap[string, T]
	labels              []string
	defaultConsctructor func() T
}

func NewMetricVector[T Metric](defaultConsctructor func() T, labels ...string) *MetricVector[T] {
	return &MetricVector[T]{
		data:                concurrent.NewConcurrentMap[string, T](),
		labels:              labels,
		defaultConsctructor: defaultConsctructor,
	}
}

func (metricVector *MetricVector[T]) WithLabelValues(labels ...string) T {
	if len(labels) != len(metricVector.labels) {
		panic("[Metrics] [WithLabels] [Error] mismatch labels count")
	}

	key := strings.Join(labels, ",")
	result, ok := metricVector.data.Find(key)
	if !ok {
		newValue := metricVector.defaultConsctructor()
		metricVector.data.Insert(key, newValue)
		result = newValue
	}

	return result
}
