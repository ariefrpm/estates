package interfaces

import (
	"context"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/google/uuid"
)

type EstateUsecase interface {
	CreateEstate(ctx context.Context, width int, length int) (*domain.Estate, error)
	CreateTree(ctx context.Context, estateID uuid.UUID, plot domain.Plot, height int) (*domain.Tree, error)
	GetEstateStats(ctx context.Context, estateID uuid.UUID) (*domain.EstateStats, error)
	GetDroneDistance(ctx context.Context, estateID uuid.UUID, maxDistance *int) (*domain.DroneDistance, error)
}

type EstateRepository interface {
	CreateEstateAndDroneRoute(ctx context.Context, estate *domain.Estate, droneRoutes []domain.DroneRoute) error
	CreateTreeAndUpdateDroneRoute(ctx context.Context, estateID uuid.UUID, droneRouteAltitude int, tree *domain.Tree) error
	GetEstateAndStats(ctx context.Context, estateID uuid.UUID) (*domain.Estate, *domain.EstateStats, error)
	GetDroneRoutes(ctx context.Context, estateID uuid.UUID) ([]domain.DroneRoute, error)
	GetEstateAndTree(ctx context.Context, estateID uuid.UUID, plot domain.Plot) (*domain.Estate, *domain.Tree, error)
}
