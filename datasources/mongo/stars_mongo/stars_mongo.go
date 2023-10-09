package stars_mongo

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	globalSession *mgo.Session
)

const (
	mongo_connection_uri = "mongo_connection_uri"
)

var (
	connection_uri = os.Getenv(mongo_connection_uri)
)

func GetSession() (*mgo.Session, error) {
	if globalSession == nil {
		var err error
		globalSession, err = makeMongoSession()
		if err != nil {
			return nil, err
		}

		globalSession.SetMode(mgo.Monotonic, true)
	}
	return globalSession.Copy(), nil
}

func makeMongoSession() (*mgo.Session, error) {

	// Our connection string - would be an environment variable most likely
	// Note: we changed the host to the first shard
	if strings.Contains(connection_uri, "mongodb+srv") {
		parts, err := url.Parse(connection_uri)

		if err != nil {
			return nil, err
		}
		password, isSet := parts.User.Password()
		if !isSet {
			return nil, errors.New("Error parsing Mongo password (env var)")
		}

		// Build our list of hosts (currently set to 3)
		mongoHost := []string{
			parts.Host,
			strings.ReplaceAll(parts.Host, "00-00", "00-01"),
			strings.ReplaceAll(parts.Host, "00-00", "00-02"),
		}

		dialInfo := mgo.DialInfo{
			Addrs:    mongoHost,
			Timeout:  15 * time.Second,
			Database: strings.ReplaceAll(parts.Path, "/", ""),
			Username: parts.User.Username(),
			Password: password,
			Source:   "admin", // auth db
		}
		// Connect with TLS
		tlsConfig := &tls.Config{}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig) // add TLS config
			return conn, err
		}
		return mgo.DialWithInfo(&dialInfo)
	} else {
		return mgo.Dial(connection_uri)
	}
}

func init() {
	var err error
	globalSession, err = makeMongoSession()

	if err != nil {
		fmt.Println("Dei erro ao iniciar a conexão", err)
		panic(err)
	}
	fmt.Println("Conexão (mongo) iniciada com sucesso")

	globalSession.SetMode(mgo.Monotonic, true)
}
