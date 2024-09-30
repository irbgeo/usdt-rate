package controller

import (
	"context"

	rateerr "github.com/irbgeo/usdt-rate/internal/utils/rate-error"
)

var tokenA = "usdt"

type controller struct {
	rate    rateProvider
	storage storage
}

//go:generate mockery --name rateProvider --structname RateProvider
type rateProvider interface {
	Get(ctx context.Context, in Pair) (out *Rate, err error)
}

//go:generate mockery --name storage --structname Storage
type storage interface {
	Create(ctx context.Context, in Rate) error
}

func newService(r rateProvider, s storage) *controller {
	return &controller{
		rate:    r,
		storage: s,
	}
}

func (c *controller) Rate(ctx context.Context, in Pair) (*Rate, error) {
	in.TokenA = tokenA

	out, err := c.rate.Get(ctx, in)
	if err != nil {
		return nil, rateerr.New(err, "msg", "failed to get rate")
	}

	err = c.storage.Create(ctx, *out)
	if err != nil {
		return nil, rateerr.New(err, "msg", "failed to store rate")
	}

	return out, nil
}
