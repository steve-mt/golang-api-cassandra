package users

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/SteveAzz/stream-api/cassandra"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// Get all users inside of our persistant storage
func Get(w http.ResponseWriter, r *http.Request) {
	var userList []User
	m := map[string]interface{}{}

	query := "SELECT id,age,firstname,lastname,city,email FROM users"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userList = append(userList, User{
			ID:        m["id"].(gocql.UUID),
			Age:       m["age"].(int),
			FirstName: m["firstname"].(string),
			LastName:  m["lastname"].(string),
			Email:     m["email"].(string),
			City:      m["city"].(string),
		})

		m = map[string]interface{}{}
	}

	json.NewEncoder(w).Encode(AllUsersResponse{Users: userList})
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	var user User
	var errs []string
	found := false

	vars := mux.Vars(r)
	id := vars["user_uuid"]

	uuid, err := gocql.ParseUUID(id)

	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,age,firstname,lastname,city,email FROM users WHERE id=? LIMIT 1"

		iterable := cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()

		for iterable.MapScan(m) {
			found = true

			user = User{
				ID:        m["id"].(gocql.UUID),
				Age:       m["age"].(int),
				FirstName: m["firstname"].(string),
				LastName:  m["lastname"].(string),
				Email:     m["email"].(string),
				City:      m["city"].(string),
			}
		}

		if !found {
			errs = append(errs, "User not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetUserResponse{User: user})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

func Enrich(uuids []gocql.UUID) map[string]string {
	if len(uuids) == 0 {
		return map[string]string{}
	}

	names := map[string]string{}
	m := map[string]interface{}{}

	query := "SELECT id,firstname,lastname FROM users WHERE id IN ?"

	iterable := cassandra.Session.Query(query, uuids).Iter()

	for iterable.MapScan(m) {
		fmt.Println("m", m)

		user_id := m["id"].(gocql.UUID)
		names[user_id.String()] = fmt.Sprintf("%s %s", m["firstname"].(string), m["lastname"].(string))
		m = map[string]interface{}{}
	}

	return names
}
