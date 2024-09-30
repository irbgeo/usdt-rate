package postgres

import (
	"context"
	_ "embed"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/irbgeo/usdt-rate/internal/controller"
)

var (
	//go:embed query/insert_rate.sql
	insertRateQuery string
)

type postgres struct {
	db *sqlx.DB
}

func NewDriver(opts StartOpts) (*postgres, error) {
	db, err := sqlx.Connect("postgres", "host="+opts.Host+" port="+strconv.Itoa(opts.Port)+" user="+opts.Username+" password="+opts.Password+" dbname="+opts.Name+" sslmode=disable")
	if err != nil {
		return nil, err
	}

	d := &postgres{
		db: db,
	}

	err = d.Migrate()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *postgres) Migrate() error {
	driver, err := pg.WithInstance(s.db.DB, &pg.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (s *postgres) InsertRate(ctx context.Context, in controller.Rate) error {
	_, err := s.db.NamedExecContext(ctx, insertRateQuery, toPostgresRate(in))
	return err
}
