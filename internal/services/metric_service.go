package services

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "code", "method"},
	)
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests made.",
		}, []string{"handler", "code", "method"})

	serverCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_count",
			Help: "Count of game server",
		}, []string{})
)

type MetricService interface {
	Register()
	Handler() http.Handler
	IncServerCount()
	DecServerCount()
	ObserveRequestDuration(handler, code, method string, duration float64)
	IncRequestTotal(handler, code, method string)
}

type metricService struct {
	requestDuration *prometheus.HistogramVec
	requestTotal    *prometheus.CounterVec
	serverCount     *prometheus.GaugeVec
}

func NewMetricService() MetricService {
	return &metricService{
		requestDuration: requestDuration,
		requestTotal:    requestTotal,
		serverCount:     serverCount,
	}
}

func (s *metricService) Register() {
	prometheus.MustRegister(s.requestDuration)
	prometheus.MustRegister(s.requestTotal)
	prometheus.MustRegister(s.serverCount)
}

func (s *metricService) Handler() http.Handler {
	return promhttp.Handler()
}

func (s *metricService) IncServerCount() {
	s.serverCount.WithLabelValues().Inc()
}

func (s *metricService) DecServerCount() {
	s.serverCount.WithLabelValues().Dec()
}

func (s *metricService) ObserveRequestDuration(handler, code, method string, duration float64) {
	s.requestDuration.WithLabelValues(handler, code, method).Observe(duration)
}

func (s *metricService) IncRequestTotal(handler, code, method string) {
	s.requestTotal.WithLabelValues(handler, code, method).Inc()
}
