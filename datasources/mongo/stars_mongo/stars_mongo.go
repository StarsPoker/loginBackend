package stars_mongo

import (
	"fmt"
	"os"
	"log"
	"gopkg.in/mgo.v2"
)

var (
	globalSession *mgo.Session
)

const (
	mongo_db_username = "mongo_db_username"
	mongo_db_password = "mongo_db_password"
	mongo_db_host     = "mongo_db_host"
	mongo_db_port     = "mongo_db_port"
)

var (
	host     = os.Getenv(mongo_db_host)
	port     = os.Getenv(mongo_db_port)
	username = os.Getenv(mongo_db_username)
	password = os.Getenv(mongo_db_password)
)

func GetSession() (*mgo.Session, error) {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "27017"
	}

	if globalSession == nil {
		var err error
		globalSession, err = mgo.Dial(username + ":" + password + "@" + host + ":" + port)
		if err != nil {
			return nil, err
		}

		globalSession.SetMode(mgo.Monotonic, true)
	}
	return globalSession.Copy(), nil
}

func init() {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "27017"
	}
	var err error

	globalSession, err = mgo.Dial(username + ":" + password + "@" + host + ":" + port)
	if err != nil {
		fmt.Println("Erro ao iniciar o mongo: ", err)
		//		panic(err)
	}

	globalSession.SetMode(mgo.Monotonic, true)
	log.Println("Mongo successfully configured.")
}
