package messages

import (
	"net/http"

	"fmt"

	"encoding/json"

	getstream "github.com/GetStream/stream-go"
	"github.com/SteveAzz/stream-api/cassandra"
	"github.com/SteveAzz/stream-api/stream"
	"github.com/gocql/gocql"
)

// Post - handles POST request to /messages/new to create a new message
// params:
// w - response write for building JSON payload response
// r - request reader to fetch form data or url params
func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var errStr, userIDstr, message string

	if userIDstr, errStr = processFormField(r, "userID"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}
	userID, err := gocql.ParseUUID(userIDstr)
	if err != nil {
		errs = append(errs, "Parameter 'userID' not a UUID")
	}

	if message, errStr = processFormField(r, "message"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}

	gocqlUUID := gocql.TimeUUID()

	created := false
	if len(errs) == 0 {
		if err := cassandra.Session.Query(`
			INSERT INTO messages (id, user_id, message) VALUES (?, ?, ?)`,
			gocqlUUID, userID, message).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	if created {
		// Send message to stream
		globalMessages, err := stream.Client.FlatFeed("messages", "global")

		if err == nil {
			globalMessages.AddActivity(&getstream.Activity{
				Actor:  fmt.Sprintf("user:%s", userID.String()),
				Verb:   "post",
				Object: fmt.Sprintf("object:%s", gocqlUUID.String()),
				MetaData: map[string]string{
					// add as many custom keys/values here as you like
					"message": message,
				},
			})

			json.NewEncoder(w).Encode(NewMessageResponse{ID: gocqlUUID})
		}
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}

}
