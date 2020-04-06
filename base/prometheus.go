package base

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	results  = make(map[string]prometheus.Collector)
	registry = Exporter{make(map[string]*promMetric)}
)

type Exporter struct {
	registry map[string]*promMetricVec
}

type promMetricVec struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	metrics   []*promMetric
}

type promMetric struct {
	value  float64
	labels []string
}

func (e *Exporter) NewMetricVec(name string, help string, labels []string) {
	if _, ok := e.registry[name]; ok == true {
		// TODO what to do if already registerd
		return
	}
	desc := prometheus.NewDesc(name, help, labels, prometheus.Labels{})
	e.registry[name] = &promMetricVec{
		desc,
		prometheus.GaugeValue,
		[]*promMetric{},
	}

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.collectGauges(ch)
}

func (e *Exporter) collectGauges(ch chan<- prometheus.Metric) {
	for _, mv := range e.registry {
		for _, m := range mv.metrics {
			ch <- prometheus.MustNewConstMetric(
				mv.desc, mv.valueType, m.value, m.labels...,
			)
		}
	}
}
