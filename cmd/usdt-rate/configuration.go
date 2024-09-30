package main

import (
	"errors"
	"flag"

	"github.com/kelseyhightower/envconfig"

	rateerr "github.com/irbgeo/usdt-rate/internal/utils/rate-error"
)

var (
	errDBUsernameEmpty = errors.New("DB_USERNAME is empty")
	errDBNameEmpty     = errors.New("DB_NAME is empty")
)

type configuration struct {
	api *apiConfiguration
	db  *dbConfiguration
}

type dbConfiguration struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     int    `envconfig:"PORT" default:"5432"`
	Username string `envconfig:"USERNAME" default:""`
	Password string `envconfig:"PASSWORD" default:""`
	Name     string `envconfig:"NAME" default:""`
}

type apiConfiguration struct {
	Port int `envconfig:"PORT" default:"8080"`
}

func readConfig() (*configuration, error) {
	cfg := &configuration{}

	dbCfg, err := readDBConfig()
	if err != nil {
		return nil, rateerr.New(err, "cfg", "db")
	}
	cfg.db = dbCfg

	apiCfg, err := readAPIConfig()
	if err != nil {
		return nil, rateerr.New(err, "cfg", "api")
	}
	cfg.api = apiCfg

	return cfg, nil
}

func readDBConfig() (*dbConfiguration, error) {
	config := &dbConfiguration{}
	err := envconfig.Process("DB", config)
	if err != nil {
		return nil, err
	}

	// Define command-line flags
	dbHost := flag.String("db-host", "localhost", "Database host")
	dbPort := flag.Int("db-port", 5432, "Database port")
	dbUsername := flag.String("db-username", "", "Database username")
	dbPassword := flag.String("db-password", "", "Database password")
	dbName := flag.String("db-name", "", "Database name")

	// Parse command-line flags
	flag.Parse()

	// Override configuration with command-line flags
	if *dbHost != "" {
		config.Host = *dbHost
	}

	if *dbPort != 0 {
		config.Port = *dbPort
	}

	if *dbUsername != "" {
		config.Username = *dbUsername
	}

	if *dbPassword != "" {
		config.Password = *dbPassword
	}

	if *dbName != "" {
		config.Name = *dbName
	}

	return config, nil
}

func validateDBConfig(config *dbConfiguration) error {
	if config.Username == "" {
		return errDBUsernameEmpty
	}
	if config.Name == "" {
		return errDBNameEmpty
	}
	return nil
}

func readAPIConfig() (*apiConfiguration, error) {
	config := &apiConfiguration{}
	err := envconfig.Process("API", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
