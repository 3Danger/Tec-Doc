package internalserver

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Collector   *prometheus.CounterVec
	LeadTime    *prometheus.HistogramVec
	Rating      *prometheus.GaugeVec
	LeadTimeQua *prometheus.SummaryVec
}

const (
	ExternalServerComponent = "external server"
	InternalServerComponent = "internal server"

	Component = "component"
	Struct    = "struct"
	Method    = "method"
	Status    = "status"
)

func NewMetric(namespace, subsystem string) *Metrics {
	var (
		collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "collector",
				Help:      "The total number of err/success/events/http",
			},
			[]string{Component, Struct, Method, Status},
		)

		leadTime = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "lead_time",
				Help:      "Execution time of something",
			},
			[]string{Component, Struct, Method, Status},
		)

		rating = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "rating",
				Help:      "Current indicator of something",
			},
			[]string{Component, Struct, Method, Status},
		)

		leadTimeQua = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "lead_time_qua",
				Help:      "Execution time of something and quantiles",
			},
			[]string{Component, Struct, Method, Status},
		)
	)
	prometheus.MustRegister(collector)
	prometheus.MustRegister(leadTime)
	prometheus.MustRegister(rating)
	prometheus.MustRegister(leadTimeQua)

	return &Metrics{
		Collector:   collector,
		LeadTime:    leadTime,
		Rating:      rating,
		LeadTimeQua: leadTimeQua,
	}
}

//func (m *Metric) Metrics(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		wi := &http.ResponseWriter()
//	})
//}
