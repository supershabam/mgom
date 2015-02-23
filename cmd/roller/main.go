package main

import (
	"flag"
	"log"
	"time"

	"github.com/supershabam/mgom"
	"gopkg.in/mgo.v2"
)

var (
	url = flag.String("url", "", "mongodb url to connect to")
	ns  = flag.String("ns", "", "namespace to watch for sample inserts")
	tns = flag.String("tns", "", "target namespace to write triggers into")
)

func main() {
	flag.Parse()
	sess, err := mgo.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}
	sess.SetMode(mgo.Eventual, true)
	defer sess.Close()

	roller := mgom.Roller{
		Sess:             sess,
		Namespace:        *ns,
		TriggerNamespace: *tns,
		Period:           time.Second,
	}
	if err := roller.Run(); err != nil {
		log.Fatal(err)
	}
}
