package storage

import (
	"context"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

type storage struct {
	driver driver
}

type driver interface {
	InsertRate(ctx context.Context, in controller.Rate) error
}

func newService(d driver) *storage {
	return &storage{
		driver: d,
	}
}

func (s *storage) Create(ctx context.Context, in controller.Rate) error {
	return s.driver.InsertRate(ctx, in)
}
