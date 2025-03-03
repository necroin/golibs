package metrics

import (
	"io"
	"strings"

	"github.com/necroin/golibs/libs/concurrent"
)

type Labels map[string]string

type MetricJsonData struct {
	Description Description `json:"description"`
	Data        any         `json:"data"`
}

type MetricVectorJsonData struct {
	Description Description `json:"description"`
	Labels      []string    `json:"labels"`
	Data        any         `json:"data"`
}

type Metric interface {
	Description() *Description
	Write(io.Writer)
	JsonData() any
	Reset()
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

func (metricVector *MetricVector[T]) With(labels Labels) T {
	labelValues := []string{}
	for _, labelName := range metricVector.labels {
		labelValues = append(labelValues, labels[labelName])
	}
	return metricVector.WithLabelValues(labelValues...)
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

func (metricVector *MetricVector[T]) IterateOverLabelValues(handler func(metric T, values ...string)) {
	metricVector.data.Iterate(func(key string, value T) {
		keyLabels := strings.Split(key, ",")
		handler(value, keyLabels...)
	})
}
