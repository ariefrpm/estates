package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var sqlOpen = sql.Open

type postgres struct {
	DB *sql.DB
}

type PostgresOptions func(*postgres)

func WithPostgresPool(maxOpenConns, maxIdleConns int, connMaxLifetimeSec int) PostgresOptions {
	return func(p *postgres) {
		p.DB.SetMaxOpenConns(maxOpenConns)
		p.DB.SetMaxIdleConns(maxIdleConns)
		p.DB.SetConnMaxLifetime(time.Duration(connMaxLifetimeSec) * time.Second)
	}
}

func NewRepository(dsn string, opts ...PostgresOptions) *postgres {
	db, err := sqlOpen("postgres", dsn)
	if err != nil {
		panic(err)
	}
	postgres := &postgres{
		DB: db,
	}
	for _, opt := range opts {
		opt(postgres)
	}
	return postgres
}
