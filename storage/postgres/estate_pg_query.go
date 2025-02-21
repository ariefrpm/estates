package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/google/uuid"
)

// GetEstateAndStats retrieves estate and estate statistics for all trees, using a LEFT JOIN on the estate table,
// allowing estates to be returned even if they have no associated tree stats.
// this is for minimizing query to db when doing validations both on estate and tree stats existense
func (p *postgres) GetEstateAndStats(ctx context.Context, estateID uuid.UUID) (*domain.Estate, *domain.EstateStats, error) {
	query := `
        SELECT e.id, e.width, e.length, m.tree_count, m.max_height, m.min_height, m.median_height
        FROM estates e LEFT JOIN estate_stats_mv m ON m.estate_id = e.id 
        WHERE e.id = $1
    `

	var estate domain.Estate

	//stats can be empty
	var statsCount *int
	var statsMax *int
	var statsMin *int
	var statsMedian *int

	err := p.DB.QueryRowContext(ctx, query, estateID).Scan(
		&estate.ID, &estate.Width, &estate.Length, &statsCount, &statsMax, &statsMin, &statsMedian,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	if statsCount == nil || statsMax == nil || statsMin == nil || statsMedian == nil {
		return &estate, nil, nil
	}

	stats := domain.EstateStats{
		Count:  *statsCount,
		Max:    *statsMax,
		Min:    *statsMin,
		Median: *statsMedian,
	}
	return &estate, &stats, nil
}

// GetDroneRoutes retrieves a list of drone routes within estates.
// The routes and altitude is precomputed when creating estates and trees.
// While this doesn't necessarily improve read performance (as we still need to sum the total distances),
// it makes the data structure cleaner and more intuitive.
func (p *postgres) GetDroneRoutes(ctx context.Context, estateID uuid.UUID) ([]domain.DroneRoute, error) {
	query := `
        SELECT route, row, col, altitude 
        FROM drone_routes 
        WHERE estate_id = $1
        ORDER BY route
    `

	rows, err := p.DB.QueryContext(ctx, query, estateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var routes []domain.DroneRoute
	for rows.Next() {
		var route domain.DroneRoute
		err := rows.Scan(&route.Route, &route.Plot.Row, &route.Plot.Col, &route.Altitude)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

// GetEstateAndTree retrieves tree along with its estate, using a LEFT JOIN on the estate table,
// allowing estates to be returned even if they have no associated trees.
// this is for minimizing query to db when doing validations both on estate and tree existense
func (p *postgres) GetEstateAndTree(ctx context.Context, estateID uuid.UUID, plot domain.Plot) (*domain.Estate, *domain.Tree, error) {
	query := `
        SELECT e.id as estate_id, e.width, e.length, t.id as tree_id, t.row, t.col, t.height 
        FROM estates e LEFT JOIN trees t ON t.estate_id = e.id AND t.row = $2 AND t.col = $3 
        WHERE e.id = $1
    `

	var estate domain.Estate

	//tree can be nil
	var treeID uuid.UUID
	var treeRow *int
	var treeCol *int
	var treeHeight *int

	err := p.DB.QueryRowContext(ctx, query, estateID, plot.Row, plot.Col).Scan(
		&estate.ID, &estate.Width, &estate.Length, &treeID, &treeRow, &treeCol, &treeHeight,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	if treeID == uuid.Nil || treeRow == nil || treeCol == nil || treeHeight == nil {
		return &estate, nil, nil
	}

	tree := domain.Tree{
		ID:     treeID,
		Plot:   domain.Plot{Row: *treeRow, Col: *treeCol},
		Height: *treeHeight,
	}
	return &estate, &tree, nil
}
