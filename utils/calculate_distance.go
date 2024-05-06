package utils

import "math"

// CalculateDistance calculate the distance between 2 coordinates.
func CalculateDistance(lon1, lat1, lon2, lat2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // Distance in kilometers
}
