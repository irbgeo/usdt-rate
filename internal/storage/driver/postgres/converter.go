package postgres

import (
	"time"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

func toPostgresRate(in controller.Rate) rate {
	return rate{
		TokenA:    in.TokenA,
		TokenB:    in.TokenB,
		Bid:       in.Bid,
		Ask:       in.Ask,
		Timestamp: time.Unix(in.Timestamp, 0),
	}
}
