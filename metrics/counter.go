package metrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/necroin/golibs/concurrent"
)

type CounterOpts struct {
	Name string
	Help string
}

type Counter struct {
	description *Description
	value       *concurrent.AtomicNumber[float64]
}

func NewCounter(opts CounterOpts) *Counter {
	return &Counter{
		description: &Description{
			Name: opts.Name,
			Type: "counter",
			Help: opts.Help,
		},
		value: concurrent.NewAtomicNumber[float64](),
	}
}

func (counter *Counter) Get() float64 {
	return counter.value.Get()
}

func (counter *Counter) Add(value float64) {
	counter.value.Add(value)
}

func (counter *Counter) Inc() {
	counter.value.Add(1)
}

func (counter *Counter) Description() *Description {
	return counter.description
}

func (counter *Counter) Write(writer io.Writer) {
	writer.Write([]byte(fmt.Sprintf("%s %v\n", counter.description.Name, counter.value.Get())))
}

type CounterVector struct {
	*MetricVector[*Counter]
	description *Description
}

func NewCounterVector(opts CounterOpts, labels ...string) *CounterVector {
	return &CounterVector{
		NewMetricVector[*Counter](func() *Counter { return NewCounter(CounterOpts{}) }, labels...),
		&Description{
			Name: opts.Name,
			Type: "counter_vector",
			Help: opts.Help,
		},
	}
}

func (counterVector *CounterVector) Description() *Description {
	return counterVector.description
}

func (counterVector *CounterVector) Write(writer io.Writer) {
	counterVector.data.Iterate(func(key string, counter *Counter) {
		labels := []string{}
		keyLabels := strings.Split(key, ",")
		for labelIndex, labelValue := range keyLabels {
			labelName := counterVector.labels[labelIndex]
			label := fmt.Sprintf("%s=\"%v\"", labelName, labelValue)
			labels = append(labels, label)
		}
		writer.Write([]byte(fmt.Sprintf("%s{%s} %v\n", counterVector.description.Name, strings.Join(labels, ","), counter.Get())))
	})
}
