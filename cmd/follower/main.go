package main

import (
	"flag"
	"log"
	"time"

	"github.com/supershabam/mgom"
	"gopkg.in/mgo.v2"
)

var (
	url = flag.String("url", "", "mongodb url")
)

// Sample is a metric datapoint
type Sample struct {
	Name  string    `bson:"name"`
	Value float64   `bson:"value"`
	At    time.Time `bson:"at"`
}

func main() {
	flag.Parse()
	session, err := mgo.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	last, err := mgom.Last(session)
	if err != nil {
		log.Fatal(err)
	}
	opch, errc := mgom.Inserts(session, last, "test.jerks")
	var sample Sample
	for op := range opch {
		err := op.Object.Unmarshal(&sample)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("op: %+v", op)
		log.Printf("sample: %+v", sample)
	}
	err = <-errc
	if err != nil {
		log.Fatal(err)
	}
}
