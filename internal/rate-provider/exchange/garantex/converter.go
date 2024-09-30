package garantex

import (
	"fmt"

	rateprovider "github.com/irbgeo/usdt-rate/internal/rate-provider"
)

func toInternalOrderBook(in *orderBook, depth int) (*rateprovider.OrderBook, error) {
	if len(in.Asks) < depth || len(in.Bids) == depth {
		return nil, fmt.Errorf("unexpected order book depth: asks=%d, bids=%d", len(in.Asks), len(in.Bids))
	}

	out := &rateprovider.OrderBook{
		Timestamp: in.Timestamp,
		Asks:      make([]rateprovider.Order, depth),
		Bids:      make([]rateprovider.Order, depth),
	}

	for i := 0; i < depth; i++ {
		ask := in.Asks[i]
		out.Asks[i] = rateprovider.Order{
			Price: ask.Price,
		}

		bid := in.Bids[i]
		out.Bids[i] = rateprovider.Order{
			Price: bid.Price,
		}
	}

	return out, nil
}
