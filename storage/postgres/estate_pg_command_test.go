package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_postgres_CreateEstateAndDroneRoute(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	pg := &postgres{DB: mockDB}

	estate := &domain.Estate{
		ID:     uuid.New(),
		Width:  100,
		Length: 200,
	}
	droneRoutes := []domain.DroneRoute{
		{Route: 1, Plot: domain.Plot{Row: 1, Col: 1}, Altitude: 1},
	}

	tests := []struct {
		name      string
		mockFunc  func()
		wantError bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO estates").WithArgs(estate.ID, estate.Width, estate.Length).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("INSERT INTO drone_routes").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		{
			name: "BeginTx error",
			mockFunc: func() {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))
			},
			wantError: true,
		},
		{
			name: "ExecContext error",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO estates").WithArgs(estate.ID, estate.Width, estate.Length).WillReturnError(errors.New("failed to execute query"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := pg.CreateEstateAndDroneRoute(ctx, estate, droneRoutes)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_postgres_CreateTreeAndUpdateDroneRoute(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	pg := &postgres{DB: mockDB}

	estateID := uuid.New()
	tree := &domain.Tree{
		ID:     uuid.New(),
		Plot:   domain.Plot{Row: 1, Col: 1},
		Height: 10,
	}

	tests := []struct {
		name      string
		mockFunc  func()
		wantError bool
	}{
		{
			name: "Success",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO trees").WithArgs(tree.ID, estateID, tree.Plot.Row, tree.Plot.Col, tree.Height).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE drone_routes").WithArgs(sqlmock.AnyArg(), estateID, tree.Plot.Row, tree.Plot.Col).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantError: false,
		},
		{
			name: "BeginTx error",
			mockFunc: func() {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))
			},
			wantError: true,
		},
		{
			name: "ExecContext error",
			mockFunc: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO trees").WithArgs(tree.ID, estateID, tree.Plot.Row, tree.Plot.Col, tree.Height).WillReturnError(errors.New("failed to execute query"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			err := pg.CreateTreeAndUpdateDroneRoute(ctx, estateID, 1, tree)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
