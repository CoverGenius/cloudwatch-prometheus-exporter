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

type BatchCollector interface {
	prometheus.Collector
	BatchUpdate([]*promMetric)
}

type Exporter struct {
	data  map[string]BatchCollector
	mutex sync.RWMutex
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, mv := range e.data {
		mv.Collect(ch)
	}
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

type BatchCounterVec struct {
	c *prometheus.CounterVec
}

func (mv *BatchGaugeVec) BatchUpdate(data []*promMetric) {
	mv.mutex.Lock()
	mv.metrics = data
	mv.mutex.Unlock()
}

func (mv *BatchGaugeVec) Collect(ch chan<- prometheus.Metric) {
	mv.mutex.Lock()
	for _, m := range mv.metrics {
		ch <- prometheus.MustNewConstMetric(
			mv.desc, prometheus.GaugeValue, m.value, m.labels...,
		)
	}
	mv.mutex.Unlock()
}

func (mv *BatchGaugeVec) Describe(ch chan<- *prometheus.Desc) {
	ch <- mv.desc
}

func (c *BatchCounterVec) BatchUpdate(data []*promMetric) {
	for _, nm := range data {
		c.c.WithLabelValues(nm.labels...).Add(nm.value)
	}
}

func (mv *BatchCounterVec) Collect(ch chan<- prometheus.Metric) {
	mv.c.Collect(ch)
}

func (mv *BatchCounterVec) Describe(ch chan<- *prometheus.Desc) {
	mv.c.Describe(ch)
}

func NewBatchGaugeVec(name string, help string, labels ...string) *BatchGaugeVec {
	return &BatchGaugeVec{
		desc:    prometheus.NewDesc(name, help, labels, prometheus.Labels{}),
		metrics: []*promMetric{},
	}

}

func NewBatchCounterVec(name string, help string, labels ...string) *BatchCounterVec {
	return &BatchCounterVec{
		c: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: help,
			},
			labels,
		),
	}

}
