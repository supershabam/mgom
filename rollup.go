package mgom

import "time"

// Rollup is an object stored into a rollup collection
type Rollup struct {
	Name string    `bson:"name"`
	At   time.Time `bson:"at"`
	P2   float64   `bson:"p2"`
	P9   float64   `bson:"p9"`
	P25  float64   `bson:"p25"`
	P50  float64   `bson:"p50"`
	P75  float64   `bson:"p75"`
	P91  float64   `bson:"p91"`
	P98  float64   `bson:"p98"`
}
