package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/irbgeo/usdt-rate/internal/controller"
	api "github.com/irbgeo/usdt-rate/pkg/api/proto"
)

type srv struct {
	grpcSrv      *grpc.Server
	svc          svc
	healthServer *health.Server
}

type svc interface {
	Rate(ctx context.Context, in controller.Pair) (*controller.Rate, error)
}

func NewServer(
	svc svc,
) *srv {
	return &srv{
		svc: svc,
	}
}

func (s *srv) ListenAndServe(
	port int,
) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.healthServer = health.NewServer()
	s.grpcSrv = grpc.NewServer()

	api.RegisterRateServiceServer(s.grpcSrv, s)
	healthpb.RegisterHealthServer(s.grpcSrv, s.healthServer)
	s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(s.grpcSrv)

	return s.grpcSrv.Serve(l)
}

func (s *srv) Shutdown() {
	s.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	stopped := make(chan struct{})
	go func() {
		s.grpcSrv.GracefulStop()
		close(stopped)
	}()

	t := time.NewTimer(30 * time.Second)
	select {
	case <-t.C:
		s.grpcSrv.Stop()
	case <-stopped:
		t.Stop()
	}

	s.healthServer.Shutdown()
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
