package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/SteveAzz/stream-api/cassandra"
	"github.com/SteveAzz/stream-api/messages"
	"github.com/SteveAzz/stream-api/stream"
	"github.com/SteveAzz/stream-api/users"
	"github.com/gorilla/mux"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	err := stream.Connect(
		"b2wtmghukzrw",
		"ypyh5e26bvvwq4h8zjusvzdqsp3b2ejhu7ybf26crxvn5yycabvgxnpc6h7ahsqf",
		"us-east")

	if err != nil {
		log.Fatal("Cloud not connect to stream, abort")
	}

	CassandraSession := cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)

	router.HandleFunc("/users", users.Get)
	router.HandleFunc("/users/new", users.Post)
	router.HandleFunc("/users/{user_uuid}", users.GetOne)

	router.HandleFunc("/messages", messages.Get)
	router.HandleFunc("/messages/new", messages.Post)
	router.HandleFunc("/messages/{message_uuid}", messages.GetOne)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
