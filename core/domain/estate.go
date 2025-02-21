package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrorEstatesNotFound = errors.New("estates not found")

type Estate struct {
	ID     uuid.UUID
	Width  int
	Length int
}

type EstateStats struct {
	Count  int
	Max    int
	Min    int
	Median int
}
