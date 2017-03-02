package cassandra

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func init() {
	connect()
}

func connect() {
	var err error

	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "streamdemoapi"
	Session, err = cluster.CreateSession()

	if err != nil {
		// Must be because the db is not ready yet
		time.Sleep(10000 * time.Millisecond)
		connect()
	}

	fmt.Println("cassandra init done!")
}
