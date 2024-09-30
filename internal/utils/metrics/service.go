package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

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
		}, []string{"token_b", "error"}),
		storage: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "storage_duration_seconds",
			Help: "Duration of storage requests",
		}, []string{"token_b", "error"}),
		rate: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "rate_duration_seconds",
			Help: "Duration of rate requests",
		}, []string{"token_b", "error"}),
	}
}

func (s *service) ObserveControllerRate(in controller.Pair, duration time.Duration, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}

	s.controller.WithLabelValues(in.TokenB, e).Observe(float64(duration))
}

func (s *service) ObserveStorageCreate(in controller.Rate, duration time.Duration, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	s.storage.WithLabelValues(in.TokenB, e).Observe(float64(duration))
}

func (s *service) ObserveRateProviderGet(in controller.Pair, duration time.Duration, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	s.rate.WithLabelValues(in.TokenB, e).Observe(float64(duration))
}
