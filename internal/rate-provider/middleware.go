package rateprovider

import (
	"context"
	"time"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

type middlewareService struct {
	svc *rateProvider

	metrics metrics
}

type metrics interface {
	ObserveRateProviderGet(in controller.Pair, duration time.Duration, err error)
}

func New(e exchange, m metrics) *middlewareService {
	return &middlewareService{
		svc:     newRateProvider(e),
		metrics: m,
	}
}

func (m *middlewareService) Get(ctx context.Context, in controller.Pair) (*controller.Rate, error) {
	var (
		startTime = time.Now()

		out *controller.Rate
		err error
	)
	defer func() {
		m.metrics.ObserveRateProviderGet(in, time.Since(startTime), err)
	}()

	out, err = m.svc.Get(ctx, in)
	if err != nil {
		return nil, err
	}

	return out, nil
}
