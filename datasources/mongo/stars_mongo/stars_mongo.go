package stars_mongo

import (
	"gopkg.in/mgo.v2"
)

var (
	globalSession *mgo.Session
)

func GetSession() (*mgo.Session, error) {
	if globalSession == nil {
		var err error
		globalSession, err = mgo.Dial(":27017")
		if err != nil {
			return nil, err
		}

		globalSession.SetMode(mgo.Monotonic, true)
	}
	return globalSession.Copy(), nil
}

func init() {
	var err error
	globalSession, err = mgo.Dial(":27017")
	if err != nil {
		panic(err)
	}

	globalSession.SetMode(mgo.Monotonic, true)
}
