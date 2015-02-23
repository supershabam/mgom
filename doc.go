// Package monup is an experiment in augmenting MongoDB data by tailing
// the oplog. In this case, we're computing statistiacal rollups of data.
//
// There are two main processes to run: `marker` and `roller`. Marker tails
// the oplog so that new datapoint inserts mark the rollup values they
// effect to be computed. Roller computes those rollups.
//
// This is an experiment and probably shouldn't be used in production.
package monup
