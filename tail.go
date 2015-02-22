package mgom

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Oplog is a document returned by tailing the replication log
type Oplog struct {
	Timestamp    bson.MongoTimestamp `bson:"ts"`
	HistoryID    int64               `bson:"h"`
	MongoVersion int                 `bson:"v"`
	Operation    string              `bson:"op"`
	Namespace    string              `bson:"ns"`
	Object       bson.RawD           `bson:"o"`
	QueryObject  bson.RawD           `bson:"o2"`
}

// Last returns the timestamp of the last seen oplog at the time of making this
// query
func Last(session *mgo.Session) (bson.MongoTimestamp, error) {
	var member Oplog
	err := session.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&member)
	return member.Timestamp, err
}

// Inserts returns a channel of insert oplogs that affect the provided namespace
// after the provided timestamp.
func Inserts(session *mgo.Session, after bson.MongoTimestamp, namespace string) (<-chan Oplog, <-chan error) {
	out := make(chan Oplog)
	errc := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			errc <- err
			close(errc)
		}()
		defer close(out)
		iter := session.DB("local").
			C("oplog.rs").
			Find(bson.M{"ts": bson.M{"$gt": after}, "ns": namespace, "op": "i"}).
			Sort("$natural").
			LogReplay().
			Tail(time.Hour)
		var ol Oplog
		for iter.Next(&ol) {
			out <- ol
		}
		err = iter.Err()
		if err != nil {
			return
		}
		err = iter.Close()
	}()
	return out, errc
}
