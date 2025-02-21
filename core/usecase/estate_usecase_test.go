package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/SawitProRecruitment/EstateService/core/interfaces"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_estateUsecase_CreateEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := interfaces.NewMockEstateRepository(ctrl)
	tests := []struct {
		name   string
		width  int
		length int
		mock   func()
		expect func() (*domain.Estate, error)
	}{
		{
			name:   "Success creating estate",
			width:  3,
			length: 2,
			mock: func() {
				droneRoutes := []domain.DroneRoute{
					{Route: 1, Plot: domain.Plot{Row: 1, Col: 1}, Altitude: 1},
					{Route: 2, Plot: domain.Plot{Row: 1, Col: 2}, Altitude: 1},
					{Route: 3, Plot: domain.Plot{Row: 2, Col: 2}, Altitude: 1},
					{Route: 4, Plot: domain.Plot{Row: 2, Col: 1}, Altitude: 1},
					{Route: 5, Plot: domain.Plot{Row: 3, Col: 1}, Altitude: 1},
					{Route: 6, Plot: domain.Plot{Row: 3, Col: 2}, Altitude: 1},
				}
				mockRepo.EXPECT().CreateEstateAndDroneRoute(gomock.Any(), gomock.Any(), droneRoutes).Return(nil)
			},
			expect: func() (*domain.Estate, error) {
				return &domain.Estate{Width: 3, Length: 2}, nil
			},
		},
		{
			name:   "Error creating estate",
			width:  2,
			length: 3,
			mock: func() {
				droneRoutes := []domain.DroneRoute{
					{Route: 1, Plot: domain.Plot{Row: 1, Col: 1}, Altitude: 1},
					{Route: 2, Plot: domain.Plot{Row: 1, Col: 2}, Altitude: 1},
					{Route: 3, Plot: domain.Plot{Row: 1, Col: 3}, Altitude: 1},
					{Route: 4, Plot: domain.Plot{Row: 2, Col: 3}, Altitude: 1},
					{Route: 5, Plot: domain.Plot{Row: 2, Col: 2}, Altitude: 1},
					{Route: 6, Plot: domain.Plot{Row: 2, Col: 1}, Altitude: 1},
				}
				mockRepo.EXPECT().CreateEstateAndDroneRoute(gomock.Any(), gomock.Any(), droneRoutes).Return(errors.New("failed to create estate"))
			},
			expect: func() (*domain.Estate, error) {
				return nil, errors.New("failed to create estate")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			e := NewEstateUsecase(mockRepo)
			got, err := e.CreateEstate(context.Background(), tt.width, tt.length)
			gotExpect, errExpect := tt.expect()
			if gotExpect != nil {
				assert.Equal(t, gotExpect.Width, got.Width)
				assert.Equal(t, gotExpect.Length, got.Length)
			} else {
				assert.Nil(t, got)
			}
			assert.Equal(t, errExpect, err)
		})
	}
}

func Test_estateUsecase_CreateTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := interfaces.NewMockEstateRepository(ctrl)
	tests := []struct {
		name     string
		estateID uuid.UUID
		plot     domain.Plot
		height   int
		mock     func()
		expect   func() (*domain.Tree, error)
	}{
		{
			name:     "Success creating tree",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Estate{
					ID:     [16]byte{12},
					Width:  10,
					Length: 10,
				}, nil, nil)
				mockRepo.EXPECT().CreateTreeAndUpdateDroneRoute(gomock.Any(), gomock.Any(), 21, gomock.Any()).Return(nil)
			},
			expect: func() (*domain.Tree, error) {
				return &domain.Tree{
					ID:     [16]byte{},
					Plot:   domain.Plot{Row: 2, Col: 3},
					Height: 20,
				}, nil
			},
		},
		{
			name:     "Estate is empty",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, nil)
			},
			expect: func() (*domain.Tree, error) {
				return nil, domain.ErrorEstatesNotFound
			},
		},
		{
			name:     "Plot already have tree",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Estate{}, &domain.Tree{ID: [16]byte{1}}, nil)
			},
			expect: func() (*domain.Tree, error) {
				return nil, domain.ErrorTreeAlreadyExists
			},
		},
		{
			name:     "Tree out of bound",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Estate{
					ID:     [16]byte{12},
					Width:  1,
					Length: 1,
				}, nil, nil)
			},
			expect: func() (*domain.Tree, error) {
				return nil, domain.ErrorTreePlotOutOfBound
			},
		},
		{
			name:     "Error creating tree, failed get estates and tree",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("failed get estate"))
			},
			expect: func() (*domain.Tree, error) {
				return nil, errors.New("failed get estate")
			},
		},
		{
			name:     "Error creating tree",
			estateID: [16]byte{12},
			plot:     domain.Plot{Row: 2, Col: 3},
			height:   20,
			mock: func() {
				mockRepo.EXPECT().GetEstateAndTree(gomock.Any(), gomock.Any(), gomock.Any()).Return(&domain.Estate{
					ID:     [16]byte{12},
					Width:  10,
					Length: 10,
				}, nil, nil)
				mockRepo.EXPECT().CreateTreeAndUpdateDroneRoute(gomock.Any(), gomock.Any(), 21, gomock.Any()).Return(errors.New("failed to create tree"))
			},
			expect: func() (*domain.Tree, error) {
				return nil, errors.New("failed to create tree")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			e := NewEstateUsecase(mockRepo)
			got, err := e.CreateTree(context.Background(), tt.estateID, tt.plot, tt.height)
			gotExpect, errExpect := tt.expect()
			if gotExpect != nil {
				assert.Equal(t, gotExpect.Height, got.Height)
				assert.Equal(t, gotExpect.Plot, got.Plot)
			} else {
				assert.Nil(t, got)
			}
			assert.Equal(t, errExpect, err)
		})
	}
}

func Test_estateUsecase_GetEstateStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := interfaces.NewMockEstateRepository(ctrl)
	tests := []struct {
		name     string
		estateID uuid.UUID
		mock     func()
		expect   func() (*domain.EstateStats, error)
	}{
		{
			name:     "Success creating tree",
			estateID: [16]byte{12},
			mock: func() {
				mockRepo.EXPECT().GetEstateAndStats(gomock.Any(), [16]byte{12}).Return(&domain.Estate{}, &domain.EstateStats{
					Count:  3,
					Max:    21,
					Min:    4,
					Median: 8,
				}, nil)
			},
			expect: func() (*domain.EstateStats, error) {
				return &domain.EstateStats{
					Count:  3,
					Max:    21,
					Min:    4,
					Median: 8,
				}, nil
			},
		},
		{
			name:     "Error creating tree",
			estateID: [16]byte{12},
			mock: func() {
				mockRepo.EXPECT().GetEstateAndStats(gomock.Any(), [16]byte{12}).Return(nil, nil, errors.New("failed to get stats"))
			},
			expect: func() (*domain.EstateStats, error) {
				return nil, errors.New("failed to get stats")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			e := NewEstateUsecase(mockRepo)
			got, err := e.GetEstateStats(context.Background(), tt.estateID)
			gotExpect, errExpect := tt.expect()
			assert.Equal(t, gotExpect, got)
			assert.Equal(t, errExpect, err)
		})
	}
}

func Test_estateUsecase_GetDroneDistance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := interfaces.NewMockEstateRepository(ctrl)
	u := NewEstateUsecase(repo)

	ctx := context.Background()
	id := uuid.New()

	tests := []struct {
		name   string
		mock   func()
		expect func() (*domain.DroneDistance, error)
	}{
		{
			name: "success - multiple routes",
			mock: func() {
				repo.EXPECT().GetDroneRoutes(ctx, id).Return([]domain.DroneRoute{
					{Altitude: 1}, {Altitude: 6}, {Altitude: 4}, {Altitude: 5}, {Altitude: 1},
				}, nil)
			},
			expect: func() (*domain.DroneDistance, error) {
				return &domain.DroneDistance{
					Distance: domain.DistanceBetweenPlot*4 + 6 + 2 + 1 + 5, // Horizontal + vertical including takeoff and landing,
				}, nil
			},
		},
		{
			name: "failure - repo error",
			mock: func() {
				repo.EXPECT().GetDroneRoutes(ctx, id).Return(nil, errors.New("repo error"))
			},
			expect: func() (*domain.DroneDistance, error) {
				return nil, errors.New("repo error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := u.GetDroneDistance(ctx, id, nil)
			gotExpect, errExpect := tt.expect()
			assert.Equal(t, gotExpect, got)
			assert.Equal(t, errExpect, err)
		})
	}
}
