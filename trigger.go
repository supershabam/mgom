package monup

import "time"

// Trigger is queued into a mongo collection to signal a rollup to create at a
// later time. It's essentially turning MongoDB into a queue, which isn't the
// best task for MongoDB.
type Trigger struct {
	Name  string
	At    time.Time
	Over  time.Duration
	After time.Time
}
