package mgom

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Marker encapsulates the logic for running a marker
type Marker struct {
	Sess             *mgo.Session
	Namespace        string
	Over             []time.Duration
	TriggerNamespace string
}

// Run executes a marker until an error is encountered
func (m Marker) Run() error {
	oplog, err := LatestOplog(m.Sess)
	if err != nil {
		return err
	}
	query := bson.M{"ts": bson.M{"$gt": oplog.Timestamp}, "ns": m.Namespace, "op": "i"}
	var sample Sample
	targetParts := strings.SplitN(m.TriggerNamespace, ".", 2)
	targetC := m.Sess.DB(targetParts[0]).C(targetParts[1])
	oplogch, errch := OplogCh(m.Sess, query)
	for oplog := range oplogch {
		err = oplog.Object.Unmarshal(&sample)
		if err != nil {
			log.Printf("continuing from unmarshalling error: %s", err)
			continue
		}
		for _, over := range m.Over {
			at := sample.At.Round(over)
			after := at.Add(over)
			trigger := Trigger{
				Name:  sample.Name,
				At:    at,
				Over:  over,
				After: after,
			}
			if _, err := targetC.Upsert(bson.M{"name": sample.Name, "at": at, "over": over}, trigger); err != nil {
				return err
			}
		}
	}
	err = <-errch
	if err != nil {
		return err
	}
	return nil
}
