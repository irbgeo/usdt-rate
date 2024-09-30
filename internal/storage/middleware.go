package storage

import (
	"context"
	"time"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

type middlewareService struct {
	svc *storage

	metrics metrics
}

type metrics interface {
	ObserveStorageCreate(in controller.Rate, duration time.Duration, err error)
}

func New(
	d driver,
	m metrics,
) *middlewareService {
	return &middlewareService{
		svc:     newService(d),
		metrics: m,
	}
}

func (m *middlewareService) Create(ctx context.Context, in controller.Rate) error {
	var (
		startTime = time.Now()

		err error
	)
	defer func() {
		m.metrics.ObserveStorageCreate(in, time.Since(startTime), err)
	}()

	err = m.svc.Create(ctx, in)
	if err != nil {
		return err
	}

	return nil
}
