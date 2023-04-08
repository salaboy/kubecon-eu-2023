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

func writeHandler(w http.ResponseWriter, r *http.Request) {

	postBody, _ := json.Marshal(map[string]string{})
	body := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://localhost:3500/v1.0/invoke/write-app/method/?message=holamamasita!", "application/json", body)
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

// func readHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()
// 	daprClient, err := dapr.NewClient()
// 	if err != nil {
// 		panic(err)
// 	}

// 	respondWithJSON(w, http.StatusOK, myValues)
// }

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

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/write", writeHandler).Methods("POST")
	// r.HandleFunc("/read", readHandler).Methods("GET")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	log.Printf("Dapr+Wazero Frontend App Started in port 8080!")
	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}
