package metrics_tests

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/necroin/golibs/libs/metrics"
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
	label       = metrics.NewLabel(metrics.LabelOpts{Name: "test_label", Help: "Label help information"})
	labelVector = metrics.NewLabelVector(
		metrics.LabelOpts{Name: "test_label_vector", Help: "Label vector help information"},
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
	registry.Register(label)
	registry.Register(labelVector)
	registry.Register(histogram)
	registry.Register(histogramVector)

	http.Handle("/metrics", registry.Handler())
	http.Handle("/metrics/json", registry.JsonHandler())

	SimMetricsWork()

	go http.ListenAndServe("localhost:3301", nil)
	time.Sleep(5 * time.Second)
	fmt.Println(histogram.Summary().String())
	fmt.Println(histogram.String())
}

func SimMetricsWork() {
	go func() {
		for {
			counter.Inc()
			counterVector.WithLabelValues("test11", "test12").Inc()
			counterVector.WithLabelValues("test21", "test22").Inc()

			gauge.Set(rand.Float64())
			gaugeVector.WithLabelValues("test11", "test12").Set(rand.Float64())
			gaugeVector.WithLabelValues("test21", "test22").Set(rand.Float64())

			label.Set(RandomString(10))
			labelVector.WithLabelValues("test11", "test12").Set(RandomString(10))
			labelVector.WithLabelValues("test21", "test22").Set(RandomString(10))

			// histogram.Observe(math.Abs(rand.NormFloat64()) * 100)
			histogram.Observe(50)

			histogramVector.WithLabelValues("test11", "test12").Observe(math.Abs(rand.Float64()) * 100)
			histogramVector.WithLabelValues("test21", "test22").Observe(math.Abs(rand.Float64()) * 100)

			// time.Sleep(1 * time.Millisecond)
		}
	}()
}
