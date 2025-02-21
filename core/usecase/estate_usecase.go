package usecase

import (
	"context"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/SawitProRecruitment/EstateService/core/interfaces"
	"github.com/google/uuid"
)

type estateUsecase struct {
	estateRepository interfaces.EstateRepository
}

func NewEstateUsecase(repo interfaces.EstateRepository) *estateUsecase {
	return &estateUsecase{
		estateRepository: repo,
	}
}

// CreateEstate create estate and initialize drone routes that are covering all state with altitude 1
func (e *estateUsecase) CreateEstate(ctx context.Context, width int, length int) (*domain.Estate, error) {
	estate := &domain.Estate{
		ID:     uuid.New(),
		Width:  width,
		Length: length,
	}

	droneRoutes := domain.DroneZigzagTraverse(width, length)
	err := e.estateRepository.CreateEstateAndDroneRoute(ctx, estate, droneRoutes)
	if err != nil {
		return nil, err
	}

	return estate, nil
}

// CreateTree create tree and adjust drone routes altitude to cover tree height
func (e *estateUsecase) CreateTree(ctx context.Context, estateID uuid.UUID, plot domain.Plot, height int) (*domain.Tree, error) {
	tree := &domain.Tree{
		ID:     uuid.New(),
		Plot:   plot,
		Height: height,
	}

	estate, existingTree, err := e.estateRepository.GetEstateAndTree(ctx, estateID, plot)
	if err != nil {
		return nil, err
	}

	if estate == nil {
		return nil, domain.ErrorEstatesNotFound
	}

	if existingTree != nil && existingTree.ID != uuid.Nil {
		return nil, domain.ErrorTreeAlreadyExists
	}

	if !tree.IsValidTreePlot(estate) {
		return nil, domain.ErrorTreePlotOutOfBound
	}

	err = e.estateRepository.CreateTreeAndUpdateDroneRoute(ctx, estateID, tree.Height+1, tree)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// GetEstateStats get estate stats of count, min, max and median of all tree
func (e *estateUsecase) GetEstateStats(ctx context.Context, estateID uuid.UUID) (*domain.EstateStats, error) {
	estate, stats, err := e.estateRepository.GetEstateAndStats(ctx, estateID)
	if err != nil {
		return nil, err
	}

	if estate == nil {
		return nil, domain.ErrorEstatesNotFound
	}

	if stats == nil {
		stats = &domain.EstateStats{}
	}

	return stats, nil
}

// GetDroneDistance get drone total distance to cover all estates plot
func (e *estateUsecase) GetDroneDistance(ctx context.Context, estateID uuid.UUID, maxDistance *int) (*domain.DroneDistance, error) {
	droneRoutes, err := e.estateRepository.GetDroneRoutes(ctx, estateID)
	if err != nil {
		return nil, err
	}

	if droneRoutes == nil {
		return nil, domain.ErrorEstatesNotFound
	}

	distance := domain.DroneTotalDistance(maxDistance, droneRoutes)

	return &domain.DroneDistance{
		Distance: distance,
	}, nil
}
