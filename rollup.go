package mgom

import "time"

// Rollup is an object stored into a rollup collection
type Rollup struct {
	Name string
	At   time.Time
	Over time.Duration
	Mode float64
	Sum  float64
	Min  float64
	Max  float64
}
