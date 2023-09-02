package tests

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/necroin/golibs/metrics"
)

var (
	counter       = metrics.NewCounter(metrics.CounterOpts{Name: "test_counter", Help: "Counter help information"})
	counterVector = metrics.NewCounterVector(
		metrics.CounterOpts{Name: "test_counter_vector", Help: "Counter vector help information"},
		"label1", "label2",
	)
	gauge       = metrics.NewGauge(metrics.GaugeOpts{Name: "test_gauge", Help: "Gauge help information"})
	gaugeVector = metrics.NewGaugeVector(
		metrics.GaugeOpts{Name: "test_gauge_vector", Help: "Gauge vector help information"},
		"label1", "label2",
	)
	histogram = metrics.NewHistogram(metrics.HistogramOpts{
		Name: "test_histogram", Help: "Histogram help information",
		Buckets: metrics.Buckets{Start: 0, Range: 10, Count: 10},
	})
	histogramVector = metrics.NewHistogramVector(
		metrics.HistogramOpts{
			Name: "test_histogram_vector", Help: "Histogram vector help information",
			Buckets: metrics.Buckets{Start: 0, Range: 10, Count: 10},
		},
		"label1", "label2",
	)
)

func TestMetrics(t *testing.T) {
	registry := metrics.NewRegistry()
	registry.Register(counter)
	registry.Register(counterVector)
	registry.Register(gauge)
	registry.Register(gaugeVector)
	registry.Register(histogram)
	registry.Register(histogramVector)

	http.Handle("/metrics", registry.Handler())

	SimMetricsWork()

	go http.ListenAndServe("localhost:3301", nil)
	time.Sleep(10 * time.Second)
}

func SimMetricsWork() {
	go func() {
		for {
			counter.Inc()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			counterVector.WithLabelValues("test11", "test12").Inc()
			counterVector.WithLabelValues("test21", "test22").Inc()
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			gauge.Set(rand.Float64())
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			gaugeVector.WithLabelValues("test11", "test12").Set(rand.Float64())
			gaugeVector.WithLabelValues("test21", "test22").Set(rand.Float64())
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			histogram.Observe(rand.Float64() * 100)
			histogram.Observe(1000)

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			histogramVector.WithLabelValues("test11", "test12").Observe(rand.Float64() * 100)
			histogramVector.WithLabelValues("test21", "test22").Observe(rand.Float64() * 100)
			time.Sleep(1 * time.Second)
		}
	}()
}
