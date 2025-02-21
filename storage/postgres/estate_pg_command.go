package postgres

import (
	"context"
	"fmt"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/google/uuid"
)

// CreateTreeAndUpdateDroneRoute Create estate and initialize drone route with empty tree, altitude is 1
func (p *postgres) CreateEstateAndDroneRoute(ctx context.Context, estate *domain.Estate, droneRoutes []domain.DroneRoute) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO estates (id, width, length) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, query, estate.ID, estate.Width, estate.Length)
	if err != nil {
		return err
	}

	// Bulk insert for drone_routes
	query = `INSERT INTO drone_routes (estate_id, route, row, col, altitude) VALUES `
	args := []interface{}{}
	argPos := 1
	for _, route := range droneRoutes {
		// constructed with placeholder, still safe from sql injections
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", argPos, argPos+1, argPos+2, argPos+3, argPos+4)
		args = append(args, estate.ID, route.Route, route.Plot.Row, route.Plot.Col, route.Altitude)
		argPos += 5
	}

	// Trim the trailing comma
	query = query[:len(query)-1]

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CreateTreeAndUpdateDroneRoute Create tree and update drone route altitude because that plot will be planted by tree
func (p *postgres) CreateTreeAndUpdateDroneRoute(ctx context.Context, estateID uuid.UUID, droneRouteAltitude int, tree *domain.Tree) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO trees (id, estate_id, row, col, height) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, query, tree.ID, estateID, tree.Plot.Row, tree.Plot.Col, tree.Height)
	if err != nil {
		return err
	}

	query = `UPDATE drone_routes SET altitude = $1 WHERE estate_id = $2 AND row = $3 AND col = $4`
	_, err = tx.ExecContext(ctx, query, droneRouteAltitude, estateID, tree.Plot.Row, tree.Plot.Col)
	if err != nil {
		return err
	}

	return tx.Commit()
}
