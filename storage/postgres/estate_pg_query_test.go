package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_postgres_GetEstateStats(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	pg := &postgres{DB: mockDB}
	estateID := uuid.New()

	tests := []struct {
		name      string
		mockFunc  func()
		wantError bool
		estate    *domain.Estate
		stats     *domain.EstateStats
	}{
		{
			name: "Success",
			mockFunc: func() {
				mock.ExpectQuery(`
                        SELECT e.id, e.width, e.length, m.tree_count, m.max_height, m.min_height, m.median_height
                        FROM estates e LEFT JOIN estate_stats_mv m ON m.estate_id = e.id 
                        WHERE e.id = \$1
                    `).
					WithArgs(estateID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "width", "length", "tree_count", "max_height", "min_height", "median_height"}).
						AddRow(estateID, 10, 10, 10, 100, 1, 50))
			},
			wantError: false,
			estate:    &domain.Estate{ID: estateID, Width: 10, Length: 10},
			stats:     &domain.EstateStats{Count: 10, Max: 100, Min: 1, Median: 50},
		},
		{
			name: "Query error",
			mockFunc: func() {
				mock.ExpectQuery(`
                        SELECT e.id, e.width, e.length, m.tree_count, m.max_height, m.min_height, m.median_height
                        FROM estates e LEFT JOIN estate_stats_mv m ON m.estate_id = e.id 
                        WHERE e.id = \$1
                    `).
					WithArgs(estateID).
					WillReturnError(errors.New("query error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			estate, stats, err := pg.GetEstateAndStats(ctx, estateID)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.estate, estate)
				assert.Equal(t, tt.stats, stats)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_postgres_GetDroneRoutes(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	pg := &postgres{DB: mockDB}
	estateID := uuid.New()

	tests := []struct {
		name      string
		mockFunc  func()
		wantError bool
		routes    []domain.DroneRoute
	}{
		{
			name: "Success",
			mockFunc: func() {
				mock.ExpectQuery("SELECT route, row, col, altitude FROM drone_routes").
					WithArgs(estateID).
					WillReturnRows(sqlmock.NewRows([]string{"route", "row", "col", "altitude"}).
						AddRow(1, 1, 1, 10))
			},
			wantError: false,
			routes:    []domain.DroneRoute{{Route: 1, Plot: domain.Plot{Row: 1, Col: 1}, Altitude: 10}},
		},
		{
			name: "Query error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT route, row, col, altitude FROM drone_routes").
					WithArgs(estateID).
					WillReturnError(errors.New("query error"))
			},
			wantError: true,
		},
		{
			name: "Row scan error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT route, row, col, altitude FROM drone_routes").
					WithArgs(estateID).
					WillReturnRows(sqlmock.NewRows([]string{"route", "row", "col", "altitude"}).
						AddRow(1, "invalid", 1, 10))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			routes, err := pg.GetDroneRoutes(ctx, estateID)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.routes, routes)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_postgres_GetTree(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	estateID := uuid.New()
	plot := domain.Plot{Row: 1, Col: 1}

	tree := domain.Tree{
		ID:     uuid.New(),
		Plot:   plot,
		Height: 10,
	}
	estate := domain.Estate{
		ID:     estateID,
		Width:  100,
		Length: 200,
	}

	repo := &postgres{DB: db}

	tests := []struct {
		name         string
		mockSetup    func()
		expectEstate *domain.Estate
		expectTree   *domain.Tree
		expectErr    error
	}{
		{
			name: "Success",
			mockSetup: func() {
				mock.ExpectQuery(`
                        SELECT e.id as estate_id, e.width, e.length, t.id as tree_id, t.row, t.col, t.height 
                        FROM estates e LEFT JOIN trees t ON t.estate_id = e.id AND t.row = \$2 AND t.col = \$3 
                        WHERE e.id = \$1 
                    `).
					WithArgs(estateID, plot.Row, plot.Col).
					WillReturnRows(sqlmock.NewRows([]string{"eid", "width", "length", "tid", "row", "col", "height"}).
						AddRow(estate.ID, estate.Width, estate.Length, tree.ID, tree.Plot.Row, tree.Plot.Col, tree.Height))
			},
			expectEstate: &estate,
			expectTree:   &tree,
			expectErr:    nil,
		},
		{
			name: "Query Error",
			mockSetup: func() {
				mock.ExpectQuery(`
                        SELECT e.id as estate_id, e.width, e.length, t.id as tree_id, t.row, t.col, t.height 
                        FROM estates e LEFT JOIN trees t ON t.estate_id = e.id AND t.row = \$2 AND t.col = \$3 
                        WHERE e.id = \$1 
                    `).
					WithArgs(estateID, plot.Row, plot.Col).
					WillReturnError(errors.New("database error"))
			},
			expectEstate: nil,
			expectTree:   nil,
			expectErr:    errors.New("database error"),
		},
		{
			name: "Not Found",
			mockSetup: func() {
				mock.ExpectQuery(`
                        SELECT e.id as estate_id, e.width, e.length, t.id as tree_id, t.row, t.col, t.height 
                        FROM estates e LEFT JOIN trees t ON t.estate_id = e.id AND t.row = \$2 AND t.col = \$3
                        WHERE e.id = \$1 
                    `).
					WithArgs(estateID, plot.Row, plot.Col).
					WillReturnError(sql.ErrNoRows)
			},
			expectEstate: nil,
			expectTree:   nil,
			expectErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			estate, tree, err := repo.GetEstateAndTree(context.Background(), estateID, plot)
			assert.Equal(t, tt.expectEstate, estate)
			assert.Equal(t, tt.expectTree, tree)
			assert.Equal(t, tt.expectErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
