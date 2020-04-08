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

// Exporter implements prometheus.Collector
type Exporter struct {
	data  map[string]BatchCollector
	mutex sync.RWMutex
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.RLock()
	for _, mv := range e.data {
		mv.Collect(ch)
	}
	e.mutex.RUnlock()
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, mv := range e.data {
		mv.Describe(ch)
	}
}

type promMetric struct {
	value  float64
	labels []string
}

type BatchGaugeVec struct {
	desc    *prometheus.Desc
	metrics []*promMetric
	mutex   sync.RWMutex
}

func (bgv *BatchGaugeVec) BatchUpdate(data []*promMetric) {
	bgv.mutex.Lock()
	bgv.metrics = data
	bgv.mutex.Unlock()
}

func (bgv *BatchGaugeVec) Collect(ch chan<- prometheus.Metric) {
	bgv.mutex.RLock()
	for _, m := range bgv.metrics {
		ch <- prometheus.MustNewConstMetric(
			bgv.desc, prometheus.GaugeValue, m.value, m.labels...,
		)
	}
	bgv.mutex.RUnlock()
}

func (bgv *BatchGaugeVec) Describe(ch chan<- *prometheus.Desc) {
	ch <- bgv.desc
}

func NewBatchGaugeVec(opts prometheus.Opts, labels []string) *BatchGaugeVec {
	name := prometheus.BuildFQName(opts.Namespace, opts.Subsystem, opts.Name)
	return &BatchGaugeVec{
		desc:    prometheus.NewDesc(name, opts.Help, labels, prometheus.Labels{}),
		metrics: []*promMetric{},
	}
}

type BatchCounterVec struct {
	c *prometheus.CounterVec
}

func (bcv *BatchCounterVec) BatchUpdate(data []*promMetric) {
	for _, nm := range data {
		bcv.c.WithLabelValues(nm.labels...).Add(nm.value)
	}
}

func (bcv *BatchCounterVec) Collect(ch chan<- prometheus.Metric) {
	bcv.c.Collect(ch)
}

func (bcv *BatchCounterVec) Describe(ch chan<- *prometheus.Desc) {
	bcv.c.Describe(ch)
}

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
