# monup

[EXPERIMENT] mongo metric rollups by tailing the oplog

https://godoc.org/github.com/supershabam/monup

## usage

```shell
#!/usr/bin/env bash
# run the marker service against localhost
# read incoming metrics from the samples collection in the metrics database
# write triggers to the triggers collection in the metrics database
# trigger rollups for datapoints in a 5 minute window, 15 minute window, and 1 hour window

marker -url="mongodb://localhost" -ns="metrics.samples" -tns="metrics.triggers" -over="5m,15m,1h"
```

```shell
#!/usr/bin/env bash
# run the roller service against localhost
# poll for triggers from the triggers collection in the metrics database
# get raw datapoints from the samples database in the metrics collection
# also write datapoints back out to the "samples_#{window}" collection e.g. samples_5m

roller -url="mongodb://localhost" -ns="metrics.samples" -tns="metrics.triggers"
```

## what you get

Now, with the marker and roller services running, whenever you write a datapoint document into the samples collection, you'll get a rollup computed.

```
// insert datapoint by any application
db.samples.insert({
  name: 'http.request_ms',
  at: new Date(),
  value: 233
})

// rollups are automatically created
db.samples_5m.find({name: 'http.request_ms'})
db.samples_15m.find({name: 'http.request_ms'})
db.samples_1h.find({name: 'http.request_ms'})
```

## MongoDB collection setup

`metrics.samples` should be a capped collection with enough space to hold all your different metric names for 3x the length of your longest rollup window.

`metrics.samples` should have an index on `{name: 1, at: 1}`

`metrics.triggers` should have an index on `{name: 1, at: 1, over: 1}`
