package messages

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/SteveAzz/stream-api/cassandra"
	"github.com/SteveAzz/stream-api/stream"
	"github.com/SteveAzz/stream-api/users"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// Get -- handle GET requests to /messages/ to fetch all messages
// params:
// w - Response writer for building JSON payload
// r - Request reader for fetch data from url or form data
func Get(w http.ResponseWriter, r *http.Request) {
	var messageList []Message
	var enrichMessages []Message
	var userList []gocql.UUID
	var err error
	m := map[string]interface{}{}

	globalMessage, err := stream.Client.FlatFeed("message", "global")

	// Fetch stream
	if err == nil {
		activies, err := globalMessage.Activities(nil)
		if err == nil {
			for _, activity := range activies.Activities {
				fmt.Println(activity)
				userID, _ := gocql.ParseUUID(activity.Actor)
				messageID, _ := gocql.ParseUUID(activity.Object)
				messageList = append(messageList, Message{
					ID:      messageID,
					UserID:  userID,
					Message: activity.MetaData["message"],
				})
				userList = append(userList, userID)
			}
		}
	}

	// If Stream fails, pull form database
	if err != nil {
		fmt.Println("Fetching activies from database")
		query := "SELECT id,user_id,message FROM messages"
		iterable := cassandra.Session.Query(query).Iter()

		for iterable.MapScan(m) {
			userID := m["UserID"].(gocql.UUID)
			messageList = append(messageList, Message{
				ID:      m["id"].(gocql.UUID),
				UserID:  userID,
				Message: m["message"].(string),
			})
			userList = append(userList, userID)
			m = map[string]interface{}{}
		}
	}

	names := users.Enrich(userList)

	for _, message := range messageList {
		message.UserFullName = names[message.UserID.String()]
		enrichMessages = append(enrichMessages, message)
	}

	fmt.Println("message list after enrichment", enrichMessages)
	json.NewEncoder(w).Encode(AllMessagesResponse{Messages: enrichMessages})
}

// GetOne -- handles GET request to /messages/{message_uuid} to fetch a message
// params:
// w - response write for building JSON payload
// r - request reader for fetching data from url and form data
func GetOne(w http.ResponseWriter, r *http.Request) {
	var message Message
	var errs []string
	found := false

	vars := mux.Vars(r)
	id := vars["message_uuid"]

	uuid, err := gocql.ParseUUID(id)

	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,user_id,message FROM messages WHERE id=? LIMIT 1"

		iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()

		for iterable.MapScan(m) {
			found = true
			userID := m["user_id"].(gocql.UUID)
			names := users.Enrich([]gocql.UUID{userID})
			fmt.Println("names", names)
			message = Message{
				ID:           userID,
				UserID:       m["userID"].(gocql.UUID),
				UserFullName: names[userID.String()],
				Message:      m["message"].(string),
			}
		}
		if !found {
			errs = append(errs, "Messages not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetMessageResponse{Message: message})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}

}
