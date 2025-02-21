package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/SawitProRecruitment/EstateService/core/interfaces"
	"github.com/SawitProRecruitment/EstateService/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_PostEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := interfaces.NewMockEstateUsecase(ctrl)
	srv := &Server{
		estateUsecase: mockUsecase,
	}
	e := echo.New()

	tests := []struct {
		name         string
		requestBody  []byte
		mockFunc     func()
		expectStatus int
	}{
		{
			name:        "Success",
			requestBody: []byte(`{"width": 10, "length": 20}`),
			mockFunc: func() {
				mockUsecase.EXPECT().CreateEstate(gomock.Any(), 10, 20).Return(&domain.Estate{ID: uuid.New()}, nil)
			},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "Invalid request body",
			requestBody:  []byte(`{invalid-json}`),
			mockFunc:     func() {},
			expectStatus: http.StatusBadRequest,
		},
		{
			name:        "Internal server error",
			requestBody: []byte(`{"width": 10, "length": 20}`),
			mockFunc: func() {
				mockUsecase.EXPECT().CreateEstate(gomock.Any(), 10, 20).Return(nil, errors.New("unexpected error"))
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/estate", io.NopCloser(bytes.NewReader(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			tt.mockFunc()
			err := srv.PostEstate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectStatus, rec.Code)
		})
	}
}

func TestServer_PostEstateIdTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := interfaces.NewMockEstateUsecase(ctrl)
	server := &Server{estateUsecase: mockUsecase}
	e := echo.New()

	estateID := uuid.New()
	tests := []struct {
		name         string
		requestBody  []byte
		mockFunc     func()
		expectStatus int
	}{
		{
			name:        "Success",
			requestBody: []byte(`{"x": 1, "y": 2, "height": 10}`),
			mockFunc: func() {
				mockUsecase.EXPECT().CreateTree(gomock.Any(), estateID, domain.Plot{Row: 2, Col: 1}, 10).Return(&domain.Tree{ID: uuid.New()}, nil)
			},
			expectStatus: http.StatusCreated,
		},
		{
			name:        "Usecase error",
			requestBody: []byte(`{"x": 1, "y": 2, "height": 10}`),
			mockFunc: func() {
				mockUsecase.EXPECT().CreateTree(gomock.Any(), estateID, domain.Plot{Row: 2, Col: 1}, 10).Return(nil, errors.New("usecase error"))
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/estate/:id/tree", io.NopCloser(bytes.NewReader(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetPath("/estate/:id/tree")
			ctx.SetParamNames("id")
			ctx.SetParamValues(estateID.String())

			tt.mockFunc()

			assert.NoError(t, server.PostEstateIdTree(ctx, estateID))
			assert.Equal(t, tt.expectStatus, rec.Code)
		})
	}
}

func TestServer_GetEstateIdStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := interfaces.NewMockEstateUsecase(ctrl)
	srv := &Server{
		estateUsecase: mockUsecase,
	}
	e := echo.New()

	tests := []struct {
		name         string
		estateID     uuid.UUID
		prepareMock  func()
		expectStatus int
	}{
		{
			name:     "Success",
			estateID: uuid.New(),
			prepareMock: func() {
				mockUsecase.EXPECT().GetEstateStats(gomock.Any(), gomock.Any()).Return(&domain.EstateStats{
					Count:  10,
					Max:    100,
					Min:    1,
					Median: 50,
				}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:     "Internal server error",
			estateID: uuid.New(),
			prepareMock: func() {
				mockUsecase.EXPECT().GetEstateStats(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/estate/:id/stats", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetPath("/estate/:id/drone-plan")
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.estateID.String())

			tt.prepareMock()
			err := srv.GetEstateIdStats(ctx, tt.estateID)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectStatus, rec.Code)
		})
	}
}

func TestServer_GetEstateIdDronePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := interfaces.NewMockEstateUsecase(ctrl)
	server := &Server{estateUsecase: mockUsecase}
	e := echo.New()

	estateID := uuid.New()
	tests := []struct {
		name         string
		mockFunc     func()
		expectStatus int
	}{
		{
			name: "Success",
			mockFunc: func() {
				mockUsecase.EXPECT().GetDroneDistance(gomock.Any(), estateID, gomock.Any()).Return(&domain.DroneDistance{Distance: 42}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name: "Usecase error",
			mockFunc: func() {
				mockUsecase.EXPECT().GetDroneDistance(gomock.Any(), estateID, gomock.Any()).Return(nil, errors.New("usecase error"))
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/estate/:id/drone-plan", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			ctx.SetPath("/estate/:id/drone-plan")
			ctx.SetParamNames("id")
			ctx.SetParamValues(estateID.String())

			tt.mockFunc()
			assert.NoError(t, server.GetEstateIdDronePlan(ctx, estateID, generated.GetEstateIdDronePlanParams{}))
			assert.Equal(t, tt.expectStatus, rec.Code)
		})
	}
}
