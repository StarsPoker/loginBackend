package stars_mongo

import (
	"os"

	"gopkg.in/mgo.v2"
)

var (
	globalSession *mgo.Session
)

const (
	mongo_db_username = "mongo_db_username"
	mongo_db_password = "mongo_db_password"
)

var (
	username = os.Getenv(mongo_db_username)
	password = os.Getenv(mongo_db_password)
)

func GetSession() (*mgo.Session, error) {
	if globalSession == nil {
		var err error
		globalSession, err = mgo.Dial(username + ":" + password + "@127.0.0.1:27017")
		if err != nil {
			return nil, err
		}

		globalSession.SetMode(mgo.Monotonic, true)
	}
	return globalSession.Copy(), nil
}

func init() {
	var err error
	globalSession, err = mgo.Dial(username + ":" + password + "@127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	globalSession.SetMode(mgo.Monotonic, true)
}
