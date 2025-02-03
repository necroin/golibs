package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	registry *Registry
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, onServeHandler := range handler.registry.onCollectandlers {
		onServeHandler()
	}

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

type JsonHandler struct {
	registry *Registry
}

func (handler JsonHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, onServeHandler := range handler.registry.onCollectandlers {
		onServeHandler()
	}

	datas := []any{}
	for _, metric := range handler.registry.metrics {
		datas = append(datas, metric.JsonData())
	}
	json.NewEncoder(writer).Encode(datas)
}

type Registry struct {
	metrics          []Metric
	onCollectandlers []func()
}

func NewRegistry() *Registry {
	return &Registry{
		metrics:          []Metric{},
		onCollectandlers: []func(){},
	}
}

func (registry *Registry) Register(metric Metric) {
	registry.metrics = append(registry.metrics, metric)
}

func (registry *Registry) Handler() Handler {
	return Handler{registry: registry}
}

func (registry *Registry) JsonHandler() JsonHandler {
	return JsonHandler{registry: registry}
}

func (registry *Registry) SetOnCollect(handlers ...func()) {
	registry.onCollectandlers = handlers
}
