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
