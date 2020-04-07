package base

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	results  = make(map[string]prometheus.Collector)
	exporter = Exporter{data: make(map[string]*promMetricVec)}
)

func init() {
	prometheus.MustRegister(&exporter)
}

type Exporter struct {
	data  map[string]*promMetricVec
	mutex sync.RWMutex
}

type promMetricVec struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	metrics   []*promMetric
	counter   *prometheus.CounterVec
	mutex     sync.RWMutex
}

type promMetric struct {
	value  float64
	labels []string
}

func newPromMetric(name string, help string, vt prometheus.ValueType, labels ...string) *promMetricVec {
	var counter *prometheus.CounterVec
	if vt == prometheus.CounterValue {
		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: help,
			},
			labels,
		)
	}
	return &promMetricVec{
		desc:      prometheus.NewDesc(name, help, labels, prometheus.Labels{}),
		valueType: vt,
		metrics:   []*promMetric{},
		counter:   counter,
	}
}

func (m *promMetricVec) Update(new []*promMetric) {
	// TODO do multiple regions work???
	if m.valueType == prometheus.GaugeValue {
		m.mutex.Lock()
		m.metrics = new
		m.mutex.Unlock()
	} else {
		for _, nm := range new {
			m.counter.WithLabelValues(nm.labels...).Add(nm.value)
		}
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, mv := range e.data {
		if mv.valueType == prometheus.GaugeValue {
			mv.mutex.Lock()
			for _, m := range mv.metrics {
				ch <- prometheus.MustNewConstMetric(
					mv.desc, mv.valueType, m.value, m.labels...,
				)
			}
			mv.mutex.Unlock()
		} else {
			mv.counter.Collect(ch)
		}
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, mv := range e.data {
		ch <- mv.desc
	}
}
