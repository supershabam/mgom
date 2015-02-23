# monup

[EXPERIMENT] mongo metric rollups by tailing the oplog

# summary

The `marker` process tails the MongoDB oplog to record when a new metric dirties a rollup of that metric and needs to be computed (and also when the earliest date it can be computed).

The `rollup` process polls the triggered marks to compute a rollup, on what metric, and over what time range.

## objects

### sample

```
{
  name: "some.metric.name",
  value: 42.0,
  at: new Date()
}
```

### trigger

```
{
  name: "some.metric.name",
  at: new Date(), // start of metric window, floored to time in line with window
  over: "5m", // time over which to aggregate
  after: new Date() // when the trigger should be processed (a future date when all data is available)
}
```

There should be a unique index on triggers ensured by marker.

```
{name: 1, over: 1, start: 1}
```

### rollup

// in metrics.samples_5m
// or metrics.samples_30m

```
{
  name: "some.metric.name"
  at: new Date() // start of metric window
  mode: 34.2,
  sum: 98.4,
  max: 58,
  min: 48,
  mean: 50,
  count: 2,
  stddev: 3.8,
  variance: 20,
  p001: 42.8,
  p01: 42.8
  p05: 42.8
  p25: 42.8
  p50: 42.8
  p75: 58.2
  p95: 58.2
  p99: 58.2
  p999: 58.2
}
```
## marker

marker watches for metric object inserts on a specified namespace and ensures a
rollups are queued for that metric

### arguments

* url - mongodb url to connect to
* ns - namespace to watch for new objects
* over - comma separated list of rollups to trigger
* tns - namespace to write rollup triggers to

`marker -url="$mongodb_url" -ns="metrics.samples" -over="5m,30m" -tns="metrics.triggers"`

## roller

roller watches for rollup triggers, calculates the rollup and writes the result into
a target collection with the range appended to the collection e.g. "5m"

### arguments

* url - mongodb url to connect to
* tns - trigger namespace to poll
* ns - source namespace to read metrics from and name to append rollup window to for writing metrics
