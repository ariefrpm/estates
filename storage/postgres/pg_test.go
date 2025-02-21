package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestWithPostgresPool(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pg := &postgres{
		DB: db,
	}

	option := WithPostgresPool(10, 5, 60)
	option(pg)

	assert.NotNil(t, pg.DB)
}

func TestNewRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dsn := "mock-dsn"

	// Mocking the sql.Open call
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		assert.Equal(t, "postgres", driverName)
		assert.Equal(t, dsn, dataSourceName)
		return db, nil
	}
	defer func() { sqlOpen = sql.Open }() // Restore original sql.Open after test

	repo := NewRepository(dsn)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)

	mock.ExpectClose()
	repo.DB.Close()
	assert.NoError(t, mock.ExpectationsWereMet())
}
