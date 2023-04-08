package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/gorilla/mux"
)

type Result struct {
	Data string `json:"data"`
}

var sub = &common.Subscription{
	PubsubName: "notifications-pubsub",
	Topic:      "notifications",
	Route:      "/notifications",
}

var notifications []string

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	r.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var result Result
		err = json.Unmarshal(data, &result)
		if err != nil {
			log.Fatal(err.Error())
		}
		notification := "Received notification: " + string(result.Data)
		fmt.Println(notification)
		notifications = append(notifications, notification)
		obj, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = w.Write(obj)
		if err != nil {
			log.Fatal(err.Error())
		}
	}).Methods("POST")

	r.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		obj, err := json.Marshal(notifications)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = w.Write(obj)
		if err != nil {
			log.Fatal(err.Error())
		}
	}).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	log.Printf("Starting Subscriber in Port: %s", appPort)
	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}
