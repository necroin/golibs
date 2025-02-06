package metrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/necroin/golibs/libs/concurrent"
)

type HistogramJsonDataItem struct {
	Buckets  Buckets   `json:"buckets,omitempty"`
	MinusInf float64   `json:"minus_inf"`
	PlusInf  float64   `json:"plus_inf"`
	Values   []float64 `json:"values"`
}

type Buckets struct {
	Start int
	Range uint
	Count uint
}

type HistogramOpts struct {
	Name    string
	Help    string
	Buckets Buckets
}

type Histogram struct {
	description *Description
	buckets     Buckets
	minusInf    *Counter
	plusInf     *Counter
	values      *concurrent.ConcurrentSlice[*Counter]
	sum         *Counter
	count       *Counter
}

func NewHistogram(opts HistogramOpts) *Histogram {
	histogram := &Histogram{
		description: &Description{
			Name: opts.Name,
			Help: opts.Help,
			Type: "histogram",
		},
		buckets:  opts.Buckets,
		minusInf: NewCounter(CounterOpts{}),
		plusInf:  NewCounter(CounterOpts{}),
		values:   concurrent.NewConcurrentSlice[*Counter](),
		sum:      NewCounter(CounterOpts{}),
		count:    NewCounter(CounterOpts{}),
	}

	for i := 0; i < int(histogram.buckets.Count); i++ {
		histogram.values.Append(NewCounter(CounterOpts{}))
	}

	return histogram
}

func (histogram *Histogram) Description() *Description {
	return histogram.description
}

func (histogram *Histogram) divAllBuckets(value float64) {
	histogram.minusInf.set(histogram.minusInf.Get() / value)
	histogram.plusInf.set(histogram.plusInf.Get() / value)

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(uint(bucketIterator))
		counter.set(counter.Get() / value)
	}
}

func (histogram *Histogram) Observe(value float64) {
	histogram.sum.Add(value)
	histogram.count.Inc()

	divValue := float64(2)
	offset := value - float64(histogram.buckets.Start)

	if offset < 0 {
		minusInfValue := histogram.minusInf.Get()
		if minusInfValue+1 < 0 {
			histogram.divAllBuckets(divValue)
		}
		histogram.minusInf.Inc()
		return
	}

	bucketId := offset / float64(histogram.buckets.Range)

	if bucketId >= float64(histogram.buckets.Count) {
		plusInfValue := histogram.plusInf.Get()
		if plusInfValue+1 < 0 {
			histogram.divAllBuckets(divValue)
		}
		histogram.plusInf.Inc()
		return
	}

	bucket, _ := histogram.values.At(uint(bucketId))
	bucketValue := bucket.Get()
	if bucketValue+1 < 0 {
		histogram.divAllBuckets(divValue)
	}
	bucket.Inc()
}

func (histogram *Histogram) Write(writer io.Writer) {
	writer.Write([]byte(fmt.Sprintf("%s{le=\"-Inf\"} %v\n", histogram.description.Name, histogram.minusInf.Get())))

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(uint(bucketIterator))
		writer.Write([]byte(fmt.Sprintf(
			"%s{ge=\"%v\",lt=\"%v\"} %v\n",
			histogram.description.Name,
			bucketIterator*int(histogram.buckets.Range),
			(bucketIterator+1)*int(histogram.buckets.Range),
			counter.Get(),
		)))
	}

	writer.Write([]byte(fmt.Sprintf("%s{ge=\"+Inf\"} %v\n", histogram.description.Name, histogram.plusInf.Get())))
	writer.Write([]byte(fmt.Sprintf("%s_sum %v\n", histogram.description.Name, histogram.sum.Get())))
	writer.Write([]byte(fmt.Sprintf("%s_count %v\n", histogram.description.Name, histogram.count.Get())))
}

func (histogram *Histogram) JsonData() any {
	values := []float64{}

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(uint(bucketIterator))
		values = append(values, counter.Get())
	}

	return MetricJsonData{
		Description: *histogram.description,
		Data: HistogramJsonDataItem{
			Buckets:  histogram.buckets,
			MinusInf: histogram.minusInf.Get(),
			PlusInf:  histogram.plusInf.Get(),
			Values:   values,
		},
	}
}

func (histogram *Histogram) Reset() {
	histogram.minusInf.Reset()
	histogram.plusInf.Reset()
	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(uint(bucketIterator))
		counter.Reset()
	}
}

type HistogramVector struct {
	*MetricVector[*Histogram]
	description *Description
	buckets     Buckets
}

func NewHistogramVector(opts HistogramOpts, labels ...string) *HistogramVector {
	return &HistogramVector{
		NewMetricVector[*Histogram](func() *Histogram { return NewHistogram(HistogramOpts{Buckets: opts.Buckets}) }, labels...),
		&Description{
			Name: opts.Name,
			Type: "histogram",
			Help: opts.Help,
		},
		opts.Buckets,
	}
}

func (histogramVector *HistogramVector) Description() *Description {
	return histogramVector.description
}

func (histogramVector *HistogramVector) Write(writer io.Writer) {
	histogramVector.data.Iterate(func(key string, histogram *Histogram) {
		labels := []string{}
		keyLabels := strings.Split(key, ",")
		for labelIndex, labelValue := range keyLabels {
			labelName := histogramVector.labels[labelIndex]
			label := fmt.Sprintf("%s=\"%v\"", labelName, labelValue)
			labels = append(labels, label)
		}
		labelsText := strings.Join(labels, ",")

		writer.Write([]byte(fmt.Sprintf("%s{%s,le=\"-Inf\"} %v\n", histogramVector.description.Name, labelsText, histogram.minusInf.Get())))

		for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
			counter, _ := histogram.values.At(uint(bucketIterator))
			writer.Write([]byte(fmt.Sprintf(
				"%s{%s,ge=\"%v\",lt=\"%v\"} %v\n",
				histogramVector.description.Name,
				labelsText,
				bucketIterator*int(histogram.buckets.Range),
				(bucketIterator+1)*int(histogram.buckets.Range),
				counter.value.Get(),
			)))
		}

		writer.Write([]byte(fmt.Sprintf("%s{%s,ge=\"+Inf\"} %v\n", histogramVector.description.Name, labelsText, histogram.plusInf.Get())))
		writer.Write([]byte(fmt.Sprintf("%s_sum{%s} %v\n", histogramVector.description.Name, labelsText, histogram.sum.Get())))
		writer.Write([]byte(fmt.Sprintf("%s_count{%s} %v\n", histogramVector.description.Name, labelsText, histogram.count.Get())))
	})
}

func (histogramVector *HistogramVector) JsonData() any {
	items := map[string]HistogramJsonDataItem{}

	histogramVector.data.Iterate(func(key string, histogram *Histogram) {
		values := []float64{}

		for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
			counter, _ := histogram.values.At(uint(bucketIterator))
			values = append(values, counter.Get())
		}

		items[key] = HistogramJsonDataItem{
			Buckets:  histogram.buckets,
			MinusInf: histogram.minusInf.Get(),
			PlusInf:  histogram.plusInf.Get(),
			Values:   values,
		}
	})

	return MetricVectorJsonData{
		Description: Description{
			Name: histogramVector.description.Name,
			Type: "histogram_vector",
			Help: histogramVector.description.Help,
		},
		Labels: histogramVector.labels,
		Data:   items,
	}
}

func (histogramVector *HistogramVector) Reset() {
	histogramVector.data.Iterate(func(key string, histogram *Histogram) {
		histogram.Reset()
	})
}
