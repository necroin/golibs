package metrics

import (
	"fmt"
	"net/http"
)

type Handler struct {
	registry *Registry
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, metric := range handler.registry.metrics {
		description := metric.Description()
		if description != nil {
			writer.Write([]byte(fmt.Sprintf("# TYPE %s %s\n", description.Name, description.Type)))
			if description.Help != "" {
				writer.Write([]byte(fmt.Sprintf("# HELP %s %s\n", description.Name, description.Help)))
			}
		}
		metric.Write(writer)
	}
}

type Registry struct {
	metrics []Metric
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: []Metric{},
	}
}

func (registry *Registry) Register(metric Metric) {
	registry.metrics = append(registry.metrics, metric)
}

func (registry *Registry) Handler() Handler {
	return Handler{registry: registry}
}
