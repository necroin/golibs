package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/necroin/golibs/metrics"
)

var (
	counter       = metrics.NewCounter(metrics.CounterOpts{Name: "test_counter", Help: "Counter help information"})
	counterVector = metrics.NewCounterVector(metrics.CounterOpts{Name: "test_counter_vector", Help: "Counter vector help information"}, "label1", "label2")
)

func TestMetrics(t *testing.T) {
	registry := metrics.NewRegistry()
	registry.Register(counter)
	registry.Register(counterVector)

	http.Handle("/metrics", registry.Handler())

	go func() {
		for {
			counter.Inc()
		}
	}()

	go func() {
		for {
			counterVector.WithLabelValues("test11", "test12").Inc()
			counterVector.WithLabelValues("test21", "test22").Inc()
		}
	}()
	go http.ListenAndServe("localhost:3301", nil)
	time.Sleep(10 * time.Second)
}
