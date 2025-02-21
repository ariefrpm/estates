package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDroneZigzagTraverse(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		length   int
		expected []DroneRoute
	}{
		{
			name:   "3x3 grid",
			width:  3,
			length: 3,
			expected: []DroneRoute{
				{Route: 1, Plot: Plot{Row: 1, Col: 1}, Altitude: 1},
				{Route: 2, Plot: Plot{Row: 1, Col: 2}, Altitude: 1},
				{Route: 3, Plot: Plot{Row: 1, Col: 3}, Altitude: 1},
				{Route: 4, Plot: Plot{Row: 2, Col: 3}, Altitude: 1},
				{Route: 5, Plot: Plot{Row: 2, Col: 2}, Altitude: 1},
				{Route: 6, Plot: Plot{Row: 2, Col: 1}, Altitude: 1},
				{Route: 7, Plot: Plot{Row: 3, Col: 1}, Altitude: 1},
				{Route: 8, Plot: Plot{Row: 3, Col: 2}, Altitude: 1},
				{Route: 9, Plot: Plot{Row: 3, Col: 3}, Altitude: 1},
			},
		},
		{
			name:   "1x1 grid",
			width:  1,
			length: 1,
			expected: []DroneRoute{
				{Route: 1, Plot: Plot{Row: 1, Col: 1}, Altitude: 1},
			},
		},
		{
			name:     "0x0 grid",
			width:    0,
			length:   0,
			expected: []DroneRoute{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DroneZigzagTraverse(tt.width, tt.length)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDroneTotalDistance(t *testing.T) {
	tests := []struct {
		name             string
		droneRoutes      []DroneRoute
		expectedDistance int
	}{
		{
			name: "Route similar with examples",
			droneRoutes: []DroneRoute{
				{Altitude: 1}, {Altitude: 6}, {Altitude: 4}, {Altitude: 5}, {Altitude: 1},
			},
			expectedDistance: (4 * DistanceBetweenPlot) + (6 + 2 + 1 + 5),
		},
		{
			name: "Single point route",
			droneRoutes: []DroneRoute{
				{Altitude: 10},
			},
			expectedDistance: 10 + 10,
		},
		{
			name: "Two points route",
			droneRoutes: []DroneRoute{
				{Altitude: 3}, {Altitude: 8},
			},
			expectedDistance: (1 * DistanceBetweenPlot) + 3 + 5 + 8,
		},
		{
			name: "Flat route",
			droneRoutes: []DroneRoute{
				{Altitude: 5}, {Altitude: 5}, {Altitude: 5}, {Altitude: 5},
			},
			expectedDistance: (3 * DistanceBetweenPlot) + 5 + 5,
		},
		{
			name: "Descending route",
			droneRoutes: []DroneRoute{
				{Altitude: 10}, {Altitude: 8}, {Altitude: 6}, {Altitude: 2},
			},
			expectedDistance: (3 * DistanceBetweenPlot) + 10 + 2 + 2 + 4 + 2,
		},
		{
			name: "Route with empty tree in the middle",
			droneRoutes: []DroneRoute{
				{Altitude: 5}, {Altitude: 5}, {Altitude: 1}, {Altitude: 5}, {Altitude: 5},
			},
			expectedDistance: (4 * DistanceBetweenPlot) + 5 + 4 + 4 + 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DroneTotalDistance(nil, tt.droneRoutes)
			assert.Equal(t, tt.expectedDistance, result)
		})
	}
}
