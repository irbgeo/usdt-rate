package metrics

import "github.com/prometheus/client_golang/prometheus"

type service struct {
	controller *prometheus.HistogramVec
	storage    *prometheus.HistogramVec
	rate       *prometheus.HistogramVec
}

func NewService() *service {
	return &service{
		controller: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "controller_duration_seconds",
			Help: "Duration of controller requests",
		}, []string{"token_a", "token_b", "error"}),
		storage: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "storage_duration_seconds",
			Help: "Duration of storage requests",
		}, []string{"token_a", "token_b", "error"}),
		rate: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "rate_duration_seconds",
			Help: "Duration of rate requests",
		}, []string{"token_a", "token_b", "error"}),
	}
}

func (s *service) ObserveControllerRate(tokenA, tokenB string, duration float64, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	s.controller.WithLabelValues(tokenA, tokenB, e).Observe(duration)
}

func (s *service) ObserveStorageCreate(tokenA, tokenB string, duration float64, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	s.storage.WithLabelValues(tokenA, tokenB, e).Observe(duration)
}

func (s *service) ObserveRateRate(tokenA, tokenB string, duration float64, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	s.rate.WithLabelValues(tokenA, tokenB, e).Observe(duration)
}
