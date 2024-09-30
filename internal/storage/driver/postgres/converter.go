package postgres

import "github.com/irbgeo/usdt-rate/internal/controller"

func toPostgresRate(in controller.Rate) rate {
	return rate{
		TokenA:    in.TokenA,
		TokenB:    in.TokenB,
		Bid:       in.Bid,
		Ask:       in.Ask,
		Timestamp: in.Timestamp,
	}
}
