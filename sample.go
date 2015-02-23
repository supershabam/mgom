package mgom

import "time"

// Sample is a single metric datapoint
type Sample struct {
	Name  string    `bson:"name"`
	Value float64   `bson:"value"`
	At    time.Time `bson:"at"`
}
