package domain

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var ErrorTreeAlreadyExists = errors.New("tree already exists")
var ErrorTreePlotOutOfBound = errors.New("tree plot out of bound")

type Tree struct {
	ID     uuid.UUID
	Plot   Plot
	Height int
}

func (t *Tree) IsValidTreePlot(estate *Estate) bool {
	slog.Info("estates", "width", estate.Width, "length", estate.Length, "id", estate.ID)
	slog.Info("tree", "row", t.Plot.Row, "col", t.Plot.Col)
	return t.Plot.Col >= 1 && t.Plot.Row >= 1 && t.Plot.Col <= estate.Length && t.Plot.Row <= estate.Width
}
