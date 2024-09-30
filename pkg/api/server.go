package api

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/irbgeo/usdt-rate/internal/controller"
	api "github.com/irbgeo/usdt-rate/pkg/api/proto"
)

type srv struct {
	svc          svc
	healthServer *health.Server
}

type svc interface {
	Rate(ctx context.Context, in controller.Pair) (*controller.Rate, error)
}

func ListenAndServe(
	port int,
	svc svc,
) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	healthServer := health.NewServer()

	srv := &srv{
		svc:          svc,
		healthServer: healthServer,
	}

	s := grpc.NewServer()

	api.RegisterRateServiceServer(s, srv)
	healthpb.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(s)

	return s.Serve(l)
}

func (s *srv) GetRate(ctx context.Context, in *api.RateRequest) (*api.RateResponse, error) {
	pair := controller.Pair{
		TokenA: in.TokenA,
		TokenB: in.TokenB,
	}

	rate, err := s.svc.Rate(ctx, pair)
	if err != nil {
		return nil, err
	}

	return &api.RateResponse{
		TokenA: rate.TokenA,
		TokenB: rate.TokenB,
		Ask:    rate.Ask,
		Bid:    rate.Bid,
	}, nil
}
