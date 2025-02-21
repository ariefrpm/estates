package domain

import (
	"errors"

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
	return t.Plot.Col >= 1 && t.Plot.Row >= 1 && t.Plot.Col <= estate.Length && t.Plot.Row <= estate.Width
}
