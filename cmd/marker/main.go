package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/supershabam/mgom"
	"gopkg.in/mgo.v2"
)

var (
	url  = flag.String("url", "", "mongodb url to connect to")
	ns   = flag.String("ns", "", "namespace to watch for sample inserts")
	over = flag.String("over", "", "comma-separated list of durations to roll up")
	tns  = flag.String("tns", "", "target namespace to write triggers into")
)

func parseOver(over string) ([]time.Duration, error) {
	durs := []time.Duration{}
	for _, part := range strings.Split(over, ",") {
		dur, err := time.ParseDuration(part)
		if err != nil {
			return []time.Duration{}, err
		}
		durs = append(durs, dur)
	}
	return durs, nil
}

func main() {
	flag.Parse()
	sess, err := mgo.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}
	sess.SetMode(mgo.Eventual, true)
	defer sess.Close()
	overDurs, err := parseOver(*over)
	if err != nil {
		log.Fatal(err)
	}
	marker := mgom.Marker{
		Sess:             sess,
		Namespace:        *ns,
		Over:             overDurs,
		TriggerNamespace: *tns,
	}
	if err := marker.Run(); err != nil {
		log.Fatal(err)
	}
}
