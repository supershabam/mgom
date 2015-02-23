package mgom

import (
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
	Object       bson.Raw            `bson:"o"`
	QueryObject  bson.Raw            `bson:"o2"`
}

// LatestOplog returns the most recent oplog from the database
func LatestOplog(sess *mgo.Session) (Oplog, error) {
	var oplog Oplog
	err := sess.DB("local").C("oplog.rs").Find(nil).Sort("-$natural").One(&oplog)
	return oplog, err
}

// OplogCh returns a channel of oplogs that match the given query as they come
// into the database. If there is an error along the way, the channel is terminated
// and there will be a non-nil error ready to take from the error channel
//
// TODO let caller have a way to stop iteration
func OplogCh(sess *mgo.Session, query bson.M) (<-chan Oplog, <-chan error) {
	out := make(chan Oplog)
	errc := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			errc <- err
			close(errc)
		}()
		defer close(out)
		iter := sess.DB("local").
			C("oplog.rs").
			Find(query).
			Sort("$natural").
			LogReplay().
			Tail(-1) // tail forever
		var oplog Oplog
		for iter.Next(&oplog) {
			out <- oplog
		}
		err = iter.Err()
		if err != nil {
			return
		}
		err = iter.Close()
	}()
	return out, errc
}
