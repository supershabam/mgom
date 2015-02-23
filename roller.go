package mgom

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gonum/stat"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Roller polls for triggers and calculated rollups
type Roller struct {
	Sess             *mgo.Session
	Namespace        string
	TriggerNamespace string
	Period           time.Duration
}

// Run executes a roller until an error is encountered
func (r Roller) Run() error {
	t := time.NewTicker(r.Period)
	defer t.Stop()
	tns := strings.SplitN(r.TriggerNamespace, ".", 2)
	c := r.Sess.DB(tns[0]).C(tns[1])
	var trigger Trigger
	for {
		<-t.C
		_, err := c.Find(bson.M{"after": bson.M{"$lt": time.Now()}}).
			Sort("$natural").
			Apply(mgo.Change{Remove: true}, &trigger)
		if err != nil && err == mgo.ErrNotFound {
			continue
		}
		if err != nil {
			return err
		}
		err = r.Rollup(trigger)
		if err != nil {
			return err
		}
	}
}

// Rollup writes a rollup result based on the current samples in the database
func (r Roller) Rollup(trigger Trigger) error {
	log.Printf("handling trigger: %+v", trigger)
	datapoints := []Datapoint{}
	tns := strings.SplitN(r.Namespace, ".", 2)
	source := r.Sess.DB(tns[0]).C(tns[1])
	dest := r.Sess.DB(tns[0]).C(fmt.Sprintf("%s_%s", tns[1], trigger.Over))
	end := trigger.At.Add(trigger.Over)
	err := source.Find(bson.M{"name": trigger.Name, "at": bson.M{"$gte": trigger.At, "$lt": end}}).
		Select(bson.M{"value": true}).
		All(&datapoints)
	if err != nil {
		return err
	}
	if len(datapoints) == 0 {
		return nil
	}
	x := DatapointsToFloats(datapoints)
	sort.Float64s(x)
	// http://en.wikipedia.org/wiki/Seven-number_summary
	rollup := Rollup{
		Name: trigger.Name,
		At:   trigger.At,
		P2:   stat.Quantile(0.02, stat.Empirical, x, nil),
		P9:   stat.Quantile(0.09, stat.Empirical, x, nil),
		P25:  stat.Quantile(0.25, stat.Empirical, x, nil),
		P50:  stat.Quantile(0.50, stat.Empirical, x, nil),
		P75:  stat.Quantile(0.75, stat.Empirical, x, nil),
		P91:  stat.Quantile(0.91, stat.Empirical, x, nil),
		P98:  stat.Quantile(0.98, stat.Empirical, x, nil),
	}
	_, err = dest.Upsert(bson.M{"name": rollup.Name, "at": rollup.At}, rollup)
	if err != nil {
		return err
	}
	return nil
}
