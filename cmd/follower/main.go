package main

import (
	"flag"
	"log"

	"github.com/supershabam/mgom"
	"gopkg.in/mgo.v2"
)

var (
	url = flag.String("url", "", "mongodb url")
)

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
	for op := range opch {
		log.Printf("op: %+v", op)
	}
	err = <-errc
	if err != nil {
		log.Fatal(err)
	}
}
