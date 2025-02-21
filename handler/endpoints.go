package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SawitProRecruitment/EstateService/core/domain"
	"github.com/SawitProRecruitment/EstateService/generated"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Create an estate
// (POST /estate)
func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.CreateEstateRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request"})
	}

	estate, err := s.estateUsecase.CreateEstate(ctx.Request().Context(), req.Width, req.Length)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Internal server error"})
	}

	return ctx.JSON(http.StatusCreated, generated.CreateEstateResponse{
		Id: &estate.ID,
	})
}

// Add a tree to an estate
// (POST /estate/{id}/tree)
func (s *Server) PostEstateIdTree(ctx echo.Context, id uuid.UUID) error {
	var req generated.CreateTreeRequest

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "Invalid request"})
	}

	tree, err := s.estateUsecase.CreateTree(ctx.Request().Context(), id, domain.Plot{Row: req.Y, Col: req.X}, req.Height)
	if err != nil {
		slog.Error("error", "message", err.Error())
		switch {
		case errors.Is(err, domain.ErrorEstatesNotFound):
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: err.Error()})
		case errors.Is(err, domain.ErrorTreeAlreadyExists), errors.Is(err, domain.ErrorTreePlotOutOfBound):
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
		default:
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Internal server error"})
		}
	}

	return ctx.JSON(http.StatusCreated, generated.CreateTreeResponse{
		Id: &tree.ID,
	})
}

// Get stats for trees in an estate
// (GET /estate/{id}/stats)
func (s *Server) GetEstateIdStats(ctx echo.Context, id uuid.UUID) error {
	stats, err := s.estateUsecase.GetEstateStats(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorEstatesNotFound) {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, generated.GetEstateTreeStatsResponse{
		Count:  &stats.Count,
		Max:    &stats.Max,
		Min:    &stats.Min,
		Median: &stats.Median,
	})
}

// Get the sum distance of the drone monitoring travel in the estate
// (GET /estate/{id}/drone-plan)
func (s *Server) GetEstateIdDronePlan(ctx echo.Context, id uuid.UUID, params generated.GetEstateIdDronePlanParams) error {
	droneDistance, err := s.estateUsecase.GetDroneDistance(ctx.Request().Context(), id, params.MaxDistance)
	if err != nil {
		if errors.Is(err, domain.ErrorEstatesNotFound) {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, generated.GetEstateDronePlanResponse{
		Distance: &droneDistance.Distance,
	})
}
