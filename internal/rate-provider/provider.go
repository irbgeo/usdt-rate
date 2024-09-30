package rateprovider

import (
	"context"
	"strings"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

type rateProvider struct {
	exchange exchange
}

type exchange interface {
	OrderBook(market string, depth int) (*OrderBook, error)
}

func newRateProvider(e exchange) *rateProvider {
	return &rateProvider{
		exchange: e,
	}
}

func (r *rateProvider) Get(ctx context.Context, in controller.Pair) (out *controller.Rate, err error) {
	market := market(in)
	orderBook, err := r.exchange.OrderBook(market, 1)
	if err != nil {
		return nil, err
	}

	return &controller.Rate{
		Pair:      in,
		Ask:       orderBook.Asks[0].Price,
		Bid:       orderBook.Bids[0].Price,
		Timestamp: orderBook.Timestamp,
	}, nil
}

func market(in controller.Pair) string {
	return strings.ToLower(in.TokenA + in.TokenB)
}
