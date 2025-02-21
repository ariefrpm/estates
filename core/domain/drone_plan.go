package domain

type DroneRoute struct {
	Route    int
	Plot     Plot
	Altitude int
}

type DroneDistance struct {
	Distance int
	Rest     Plot
}

// DroneZigzagTraverse initialize drone route travel in zigzag based on estates width and length with altitude 1
// put it simply it will convert 2d array of estates into linear 1d array of drone route
func DroneZigzagTraverse(width int, length int) []DroneRoute {
	rows := width + 1
	if rows == 0 {
		return []DroneRoute{}
	}
	cols := length + 1
	result := []DroneRoute{}

	row, col := 1, 1

	i := 0
	for row < rows {
		// Move right
		for col < cols {
			i++

			result = append(result, DroneRoute{
				Route:    i,
				Plot:     Plot{Row: row, Col: col},
				Altitude: 1,
			})
			col++
		}
		col-- // Step back to stay in bounds
		row++ // Move to the next row

		if row >= rows {
			break
		}

		// Move left
		for col >= 1 {
			i++

			result = append(result, DroneRoute{
				Route:    i,
				Plot:     Plot{Row: row, Col: col},
				Altitude: 1,
			})
			col--
		}
		col++ // Step back to stay in bounds
		row++ // Move to the next row
	}

	return result
}

// DroneTotalDistance calculate drone vertical and horizontal movement distance to travel estate with tree
func DroneTotalDistance(maxDistance *int, droneRoutes []DroneRoute) int {
	horizontal := 0
	vertical := 0

	for i, route := range droneRoutes {
		// horizontal movement
		if i < len(droneRoutes)-1 {
			horizontal += DistanceBetweenPlot
		}

		// vertical movement
		// add takeoff altitude
		if i == 0 {
			vertical += route.Altitude
		}

		// calc diff between current altitude to next altitude
		currAltitude := route.Altitude
		nextAltitude := 0
		if i+1 < len(droneRoutes) {
			nextAltitude = droneRoutes[i+1].Altitude
		}
		diff := nextAltitude - currAltitude

		// diff must be positif
		if diff < 0 {
			diff *= -1
		}
		vertical += diff
	}

	return horizontal + vertical

}
