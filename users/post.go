package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SteveAzz/stream-api/cassandra"
	"github.com/gocql/gocql"
)

// Post save the user into the persistant sotrage
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUuid gocql.UUID

	user, errs := FormToUser(r)

	created := false

	if len(errs) == 0 {
		fmt.Println("creating a new user")

		// generate a UUID for the user
		gocqlUuid = gocql.TimeUUID()

		// write data to Cassandra
		if err := cassandra.Session.Query(
			`INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUuid,
			user.FirstName,
			user.LastName,
			user.Email,
			user.City,
			user.Age,
		).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	if created {
		fmt.Println("user_id", gocqlUuid)
		json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUuid})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
