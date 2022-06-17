package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Collector   *prometheus.CounterVec
	LeadTime    *prometheus.HistogramVec
	Rating      *prometheus.GaugeVec
	LeadTimeQua *prometheus.SummaryVec
}

const (
	//MainComponent           = "main"
	//ServiceComponent        = "service"
	//StoreComponent          = "store"
	//NatsComponent           = "nats"
	//ProviderComponent       = "provider"
	//WatcherComponent        = "watcher"
	ExternalServerComponent = "external server"
	InternalServerComponent = "internal server"

	//StatusSuccess  = "success"
	//StatusError    = "error"
	//StatusLeadTime = "lead time"
	//StatusHttp     = "http request"

	Component = "component"
	Method    = "method"
	Status    = "status"
	Path      = "path"
)

func NewMetrics(namespace, subsystem string) *Metrics {
	var (
		collector = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "collector",
				Help:      "The total number of err/success/events/http",
			},
			[]string{Component, Method, Path, Status},
		)

		leadTime = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "lead_time",
				Help:      "Execution time of something",
			},
			[]string{Component, Method, Path, Status},
		)

		rating = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "rating",
				Help:      "Current indicator of something",
			},
			[]string{Component, Method, Path, Status},
		)

		leadTimeQua = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "lead_time_qua",
				Help:      "Execution time of something and quantiles",
			},
			[]string{Component, Method, Path, Status},
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
