package storage

import (
	"context"
	"fmt"

	// Import pg driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}
