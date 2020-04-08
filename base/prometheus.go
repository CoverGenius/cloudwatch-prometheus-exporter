package base

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	exporter = Exporter{data: make(map[string]BatchCollector)}
)

func init() {
	prometheus.MustRegister(&exporter)
}

// BatchCollector is a prometheus.Collector which allows a metric to be
// atomically updated with multiple label combinations at once
type BatchCollector interface {
	prometheus.Collector
	BatchUpdate([]*promMetric)
}

// Exporter collects Cloudwatch metrics and exports them using the prometheus.Collector interface
type Exporter struct {
	data  map[string]BatchCollector
	mutex sync.RWMutex
}

// Collect fetches all the cached metrics stored by the CloudWatch exporter.
// Implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	for _, mv := range e.data {
		mv.Collect(ch)
	}
}

// Describe describes all the metrics exported by the CloudWatch exporter.
// Implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, mv := range e.data {
		mv.Describe(ch)
	}
}

type promMetric struct {
	value  float64
	labels []string
}

// BatchGaugeVec is a prometheus.GaugeVec which implements BatchCollector
type BatchGaugeVec struct {
	desc    *prometheus.Desc
	metrics []*promMetric
	mutex   sync.RWMutex
}

// BatchUpdate replaces the metric data for the BatchGaugeVec with the input data.
//
// Replacing with a slice of values means that any series previously populated
// which are no longer relevant will not be re-exported.
// TODO convert this to BatchSet/BatchAdd
func (bgv *BatchGaugeVec) BatchUpdate(data []*promMetric) {
	bgv.mutex.Lock()
	defer bgv.mutex.Unlock()
	bgv.metrics = data
}

// Collect implements prometheus.Collect for BatchGaugeVec
func (bgv *BatchGaugeVec) Collect(ch chan<- prometheus.Metric) {
	bgv.mutex.RLock()
	defer bgv.mutex.RUnlock()
	for _, m := range bgv.metrics {
		ch <- prometheus.MustNewConstMetric(
			bgv.desc, prometheus.GaugeValue, m.value, m.labels...,
		)
	}
}

// Describe implements prometheus.Describe for BatchGaugeVec
func (bgv *BatchGaugeVec) Describe(ch chan<- *prometheus.Desc) {
	ch <- bgv.desc
}

// NewBatchGaugeVec creates a new BatchGaugeVec based on the provided Opts and partitioned by the given label names.
func NewBatchGaugeVec(opts prometheus.Opts, labels []string) *BatchGaugeVec {
	name := prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name)
	return &BatchGaugeVec{
		desc:    prometheus.NewDesc(name, opts.Help, labels, prometheus.Labels{}),
		metrics: []*promMetric{},
	}
}

// BatchCounterVec is a prometheus.CounterVec which implements BatchCollector
type BatchCounterVec struct {
	c     *prometheus.CounterVec
	mutex sync.RWMutex
}

// BatchUpdate replaces the metric data for the BatchCounterVec with the input data.
// TODO convert this to BatchSet/BatchAdd
func (bcv *BatchCounterVec) BatchUpdate(data []*promMetric) {
	bcv.mutex.Lock()
	defer bcv.mutex.Unlock()
	for _, nm := range data {
		bcv.c.WithLabelValues(nm.labels...).Add(nm.value)
	}
}

// Collect implements prometheus.Collect for BatchCounterVec
func (bcv *BatchCounterVec) Collect(ch chan<- prometheus.Metric) {
	bcv.c.Collect(ch)
}

// Describe implements prometheus.Describe for BatchCounterVec
func (bcv *BatchCounterVec) Describe(ch chan<- *prometheus.Desc) {
	bcv.c.Describe(ch)
}

// NewBatchCounterVec creates a new BatchCounterVec based on the provided Opts and partitioned by the given label names.
func NewBatchCounterVec(opts prometheus.Opts, labels []string) *BatchCounterVec {
	return &BatchCounterVec{
		c: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      opts.Name,
				Subsystem: opts.Subsystem,
				Namespace: opts.Namespace,
				Help:      opts.Help,
			},
			labels,
		),
	}

}
