package monup

import "time"

// Sample is a single datapoint. Any process can write to MongoDB following
// this format and the mark/roll process will handle the data.
type Sample struct {
	Name  string    `bson:"name"`
	Value float64   `bson:"value"`
	At    time.Time `bson:"at"`
}
