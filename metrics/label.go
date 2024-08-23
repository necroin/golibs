package metrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/necroin/golibs/concurrent"
)

type LabelOpts struct {
	Name string
	Help string
}

type Label struct {
	description *Description
	value       *concurrent.AtomicValue[string]
}

func NewLabel(opts LabelOpts) *Label {
	return &Label{
		description: &Description{
			Name: opts.Name,
			Type: "label",
			Help: opts.Help,
		},
		value: concurrent.NewAtomicValue[string](),
	}
}

func (label *Label) Set(value string) {
	label.value.Set(value)
}

func (label *Label) Get() string {
	return label.value.Get()
}

func (label *Label) Description() *Description {
	return label.description
}

func (label *Label) Write(writer io.Writer) {
	writer.Write([]byte(fmt.Sprintf("%s %s\n", label.description.Name, label.value.Get())))
}

func (label *Label) JsonData() any {
	return MetricJsonData{
		Description: *label.description,
		Data:        label.value.Get(),
	}
}

func (label *Label) Reset() {
	label.Set("")
}

type LabelVector struct {
	*MetricVector[*Label]
	description *Description
}

func NewLabelVector(opts LabelOpts, labels ...string) *LabelVector {
	return &LabelVector{
		NewMetricVector[*Label](func() *Label { return NewLabel(LabelOpts{}) }, labels...),
		&Description{
			Name: opts.Name,
			Type: "label",
			Help: opts.Help,
		},
	}
}

func (labelVector *LabelVector) Description() *Description {
	return labelVector.description
}

func (labelVector *LabelVector) Write(writer io.Writer) {
	labelVector.data.Iterate(func(key string, label *Label) {
		labels := []string{}
		keyLabels := strings.Split(key, ",")
		for labelIndex, labelValue := range keyLabels {
			labelName := labelVector.labels[labelIndex]
			label := fmt.Sprintf("%s=\"%v\"", labelName, labelValue)
			labels = append(labels, label)
		}
		writer.Write([]byte(fmt.Sprintf("%s{%s} %v\n", labelVector.description.Name, strings.Join(labels, ","), label.value.Get())))
	})
}

func (labelVector *LabelVector) JsonData() any {
	data := map[string]string{}

	labelVector.data.Iterate(func(key string, label *Label) {
		data[key] = label.Get()
	})

	return MetricVectorJsonData{
		Description: Description{
			Name: labelVector.description.Name,
			Type: "label_vector",
			Help: labelVector.description.Help,
		},
		Labels: labelVector.labels,
		Data:   data,
	}
}

func (labelVector *LabelVector) Reset() {
	labelVector.data.Iterate(func(key string, label *Label) {
		label.Reset()
	})
}
