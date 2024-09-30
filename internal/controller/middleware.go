package controller

import (
	"context"
	"time"

	"github.com/irbgeo/usdt-rate/internal/utils/logging"
)

type middlewareService struct {
	svc *controller

	metrics metrics
}

//go:generate mockery --name metrics --structname Metrics
type metrics interface {
	ObserveControllerRate(in Pair, duration time.Duration, err error)
}

func NewService(
	r rateProvider,
	s storage,
	m metrics,
) *middlewareService {
	return &middlewareService{
		svc:     newService(r, s),
		metrics: m,
	}
}

func (m *middlewareService) Rate(ctx context.Context, in Pair) (*Rate, error) {
	var (
		startTime = time.Now()

		out *Rate
		err error
	)
	defer func() {
		m.metrics.ObserveControllerRate(in, time.Since(startTime), err)
	}()

	out, err = m.svc.Rate(ctx, in)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	return out, nil
}
