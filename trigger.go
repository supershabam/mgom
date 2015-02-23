package mgom

import "time"

// Trigger is queued into a mongo collection to signal a rollup to create at a
// later time
type Trigger struct {
	Name  string
	At    time.Time
	Over  time.Duration
	After time.Time
}
