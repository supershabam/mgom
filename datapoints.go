package mgom

// Datapoint exists to ease parsing out of MongoDB
type Datapoint struct {
	Value float64 `bson:"value"`
}

// DatapointsToFloats lets us work with a natural slice of floats
func DatapointsToFloats(dps []Datapoint) []float64 {
	floats := make([]float64, 0, len(dps))
	for _, dp := range dps {
		floats = append(floats, dp.Value)
	}
	return floats
}
