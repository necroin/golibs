package metrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/necroin/golibs/concurrent"
)

type GaugeOpts struct {
	Name string
	Help string
}

type Gauge struct {
	description *Description
	value       *concurrent.AtomicNumber[float64]
}

func NewGauge(opts GaugeOpts) *Gauge {
	return &Gauge{
		description: &Description{
			Name: opts.Name,
			Type: "gauge",
			Help: opts.Help,
		},
		value: concurrent.NewAtomicNumber[float64](),
	}
}

func (gauge *Gauge) Set(value float64) {
	gauge.value.Set(value)
}

func (gauge *Gauge) Get() float64 {
	return gauge.value.Get()
}

func (gauge *Gauge) Add(value float64) {
	gauge.value.Add(value)
}

func (gauge *Gauge) Sub(value float64) {
	gauge.value.Sub(value)
}

func (gauge *Gauge) Inc() {
	gauge.value.Add(1)
}

func (gauge *Gauge) Dec() {
	gauge.value.Sub(1)
}

func (gauge *Gauge) Description() *Description {
	return gauge.description
}

func (gauge *Gauge) Write(writer io.Writer) {
	writer.Write([]byte(fmt.Sprintf("%s %v\n", gauge.description.Name, gauge.value.Get())))
}

type GaugeVector struct {
	*MetricVector[*Gauge]
	description *Description
}

func NewGaugeVector(opts GaugeOpts, labels ...string) *GaugeVector {
	return &GaugeVector{
		NewMetricVector[*Gauge](func() *Gauge { return NewGauge(GaugeOpts{}) }, labels...),
		&Description{
			Name: opts.Name,
			Type: "gauge",
			Help: opts.Help,
		},
	}
}

func (gaugeVector *GaugeVector) Description() *Description {
	return gaugeVector.description
}

func (gaugeVector *GaugeVector) Write(writer io.Writer) {
	gaugeVector.data.Iterate(func(key string, gauge *Gauge) {
		labels := []string{}
		keyLabels := strings.Split(key, ",")
		for labelIndex, labelValue := range keyLabels {
			labelName := gaugeVector.labels[labelIndex]
			label := fmt.Sprintf("%s=\"%v\"", labelName, labelValue)
			labels = append(labels, label)
		}
		writer.Write([]byte(fmt.Sprintf("%s{%s} %v\n", gaugeVector.description.Name, strings.Join(labels, ","), gauge.value.Get())))
	})
}
