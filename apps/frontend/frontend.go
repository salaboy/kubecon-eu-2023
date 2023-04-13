package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type MyValues struct {
	Values []string
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest("DELETE", "http://localhost:3500/v1.0/invoke/write-app/method/", nil)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	defer resp.Body.Close()

	log.Println("Result: ")
	log.Println(resp.StatusCode)

	respondWithJSON(w, http.StatusOK, resp.StatusCode)
}
func writeHandler(w http.ResponseWriter, r *http.Request) {

	postBody, _ := json.Marshal(map[string]string{})
	body := bytes.NewBuffer(postBody)

	value := r.URL.Query().Get("message")

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:3500/v1.0/invoke/write-app/method/?message="+value, "application/json", body)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	log.Println("Result: ")
	log.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}

	respondWithJSON(w, http.StatusOK, resp.StatusCode)
}

func subscriptionsHandler(w http.ResponseWriter, r *http.Request) {

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Get("http://localhost:3500/v1.0/invoke/subscriber-app/method/notifications")
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	log.Println("Result: ")
	log.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		myValues := []string{}
		json.Unmarshal(bodyBytes, &myValues)

		respondWithJSON(w, http.StatusOK, myValues)

	} else {
		respondWithJSON(w, http.StatusOK, "no values yet.")
	}

}

func readHandler(w http.ResponseWriter, r *http.Request) {

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Get("http://localhost:3500/v1.0/invoke/read-app/method/")
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	log.Println("Result: ")
	log.Println(resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		myValues := MyValues{}
		json.Unmarshal(bodyBytes, &myValues)

		respondWithJSON(w, http.StatusOK, myValues.Values)

	} else {
		respondWithJSON(w, http.StatusOK, "no values yet.")
	}

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route

	r.HandleFunc("/write", writeHandler).Methods("POST")
	r.HandleFunc("/delete", deleteHandler).Methods("POST")
	r.HandleFunc("/read", readHandler).Methods("GET")
	r.HandleFunc("/subscriptions", subscriptionsHandler).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))
	http.Handle("/", r)

	log.Printf("Dapr+Wazero Frontend App Started in port 8080!")
	// Start the server; this is a blocking call
	log.Fatal(http.ListenAndServe(":"+appPort, nil))

}
