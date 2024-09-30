package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/irbgeo/usdt-rate/internal/controller"
	rateprovider "github.com/irbgeo/usdt-rate/internal/rate-provider"
	"github.com/irbgeo/usdt-rate/internal/rate-provider/exchange/garantex"
	"github.com/irbgeo/usdt-rate/internal/storage"
	"github.com/irbgeo/usdt-rate/internal/storage/driver/postgres"
	"github.com/irbgeo/usdt-rate/internal/utils/logging"
	"github.com/irbgeo/usdt-rate/pkg/api"
)

var (
	LoggerLevel = "debug"
)

func main() {
	logging.Init(LoggerLevel)

	cfg, err := readConfig()
	if err != nil {
		logging.Error(err, "msg", "failed to read config")
	}

	logging.Info("read config", "values", cfg)

	postgresOpts := postgres.StartOpts{
		Host:     cfg.db.Host,
		Port:     cfg.db.Port,
		Username: cfg.db.Username,
		Password: cfg.db.Password,
		Name:     cfg.db.Name,
	}

	postgresDriver, err := postgres.NewDriver(postgresOpts)
	if err != nil {
		logging.Error(err, "msg", "failed to create postgres driver")
		os.Exit(1)
	}

	stor := storage.New(postgresDriver, nil)

	garantexCli := garantex.NewClient()

	rateProvider := rateprovider.New(garantexCli, nil)

	ctrl := controller.NewService(rateProvider, stor, nil)

	go func() {
		err := api.ListenAndServe(cfg.api.Port, ctrl)
		if err != nil {
			logging.Error(err, "msg", "failed to start api")
			os.Exit(1)
		}
	}()

	slog.Info("I'm turned on")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	slog.Info("Goodbye!")
}
