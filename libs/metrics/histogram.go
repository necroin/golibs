package metrics

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/necroin/golibs/libs/concurrent"
	"github.com/necroin/golibs/utils"
)

type HistogramJsonDataItem struct {
	Buckets  Buckets   `json:"buckets,omitempty"`
	MinusInf float64   `json:"minus_inf"`
	PlusInf  float64   `json:"plus_inf"`
	Values   []float64 `json:"values"`
}

type HistogramBucketView struct {
	Head  string
	Count string
	Value string
}

type HistogramSummary struct {
	Count   int64
	Sum     float64
	Min     float64
	Max     float64
	Average float64
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
	min         *concurrent.AtomicValue[float64]
	max         *concurrent.AtomicValue[float64]
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
		min:      concurrent.NewAtomicValue[float64](),
		max:      concurrent.NewAtomicValue[float64](),
	}

	for i := 0; i < int(histogram.buckets.Count); i++ {
		histogram.values.Append(NewCounter(CounterOpts{}))
	}

	return histogram
}

func (histogram *Histogram) Description() *Description {
	return histogram.description
}

func (histogram *Histogram) Buckets() Buckets {
	return histogram.buckets
}

func (histogram *Histogram) MinusInf() *Counter {
	return histogram.minusInf
}

func (histogram *Histogram) PlusInf() *Counter {
	return histogram.plusInf
}

func (histogram *Histogram) Values() []*Counter {
	return histogram.values.Data()
}

func (histogram *Histogram) Sum() *Counter {
	return histogram.sum
}

func (histogram *Histogram) Count() *Counter {
	return histogram.count
}

func (histogram *Histogram) divAllBuckets(value float64) {
	histogram.minusInf.set(histogram.minusInf.Get() / value)
	histogram.plusInf.set(histogram.plusInf.Get() / value)

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(bucketIterator)
		counter.set(counter.Get() / value)
	}
}

func (histogram *Histogram) Observe(value float64) {
	histogram.sum.Add(value)
	histogram.count.Inc()

	histogram.min.SetWithCondition(value, func(oldValue, newValue float64) bool { return newValue < oldValue })
	histogram.max.SetWithCondition(value, func(oldValue, newValue float64) bool { return newValue > oldValue })

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

	bucket, _ := histogram.values.At(int(bucketId))
	bucketValue := bucket.Get()
	if bucketValue+1 < 0 {
		histogram.divAllBuckets(divValue)
	}
	bucket.Inc()
}

func (histogram *Histogram) Write(writer io.Writer) {
	fmt.Fprintf(writer, "%s{le=\"-Inf\"} %v\n", histogram.description.Name, histogram.minusInf.Get())

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(bucketIterator)
		fmt.Fprintf(writer,
			"%s{ge=\"%v\",lt=\"%v\"} %v\n",
			histogram.description.Name,
			bucketIterator*int(histogram.buckets.Range),
			(bucketIterator+1)*int(histogram.buckets.Range),
			counter.Get(),
		)
	}

	fmt.Fprintf(writer, "%s{ge=\"+Inf\"} %v\n", histogram.description.Name, histogram.plusInf.Get())
	fmt.Fprintf(writer, "%s_sum %v\n", histogram.description.Name, histogram.sum.Get())
	fmt.Fprintf(writer, "%s_count %v\n", histogram.description.Name, histogram.count.Get())
}

func (histogram *Histogram) JsonData() any {
	values := []float64{}

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		counter, _ := histogram.values.At(bucketIterator)
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
		counter, _ := histogram.values.At(bucketIterator)
		counter.Reset()
	}
}

func (histogram *Histogram) String() string {
	result := ""

	bucketsViews := []*HistogramBucketView{}
	maxBucketViewHeadLen := 0
	maxBucketViewCountLen := 0

	for bucketIterator := 0; bucketIterator < int(histogram.buckets.Count); bucketIterator++ {
		bucketEnd := histogram.buckets.Range * uint(bucketIterator+1)
		counter, _ := histogram.values.At(bucketIterator)
		percent := int(utils.SafeDivide(counter.Get(), histogram.count.Get()) * 100)
		bucketView := &HistogramBucketView{
			Head:  fmt.Sprintf("%v", bucketEnd),
			Count: fmt.Sprintf("[%v]", counter.Get()),
			Value: strings.Repeat("▪", percent/2),
		}
		if len(bucketView.Head) > maxBucketViewHeadLen {
			maxBucketViewHeadLen = len(bucketView.Head)
		}
		if len(bucketView.Count) > maxBucketViewCountLen {
			maxBucketViewCountLen = len(bucketView.Count)
		}
		bucketsViews = append(bucketsViews, bucketView)
	}

	for _, bucketView := range bucketsViews {
		result += fmt.Sprintf(
			"%s%s %s%s | %s\n",
			bucketView.Head,
			strings.Repeat(" ", maxBucketViewHeadLen-len(bucketView.Head)),
			bucketView.Count,
			strings.Repeat(" ", maxBucketViewCountLen-len(bucketView.Count)),
			bucketView.Value,
		)
	}

	return result
}

func (histogram *Histogram) Summary() HistogramSummary {
	count := histogram.count.Get()
	sum := histogram.sum.Get()
	min := histogram.min.Get()
	max := histogram.max.Get()

	return HistogramSummary{
		Count:   int64(count),
		Sum:     sum,
		Min:     min,
		Max:     max,
		Average: utils.SafeDivide(sum, count),
	}
}

func (summary HistogramSummary) String() string {
	buffer := &bytes.Buffer{}
	summaryTemplate, err := template.New("HistogramSummary").Parse("Count: {{.Count}}\nSum: {{.Sum}}\nMin: {{.Min}}\nMax: {{.Max}}\nAverage: {{.Average}}\n")
	if err != nil {
		panic(err)
	}

	if err := summaryTemplate.Execute(buffer, summary); err != nil {
		panic(err)
	}
	return buffer.String()
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
			counter, _ := histogram.values.At(bucketIterator)
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
			counter, _ := histogram.values.At(bucketIterator)
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
