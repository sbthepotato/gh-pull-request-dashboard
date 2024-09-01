package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func test_json_handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)

	// Create a response object
	response := Response{
		Message: "Hello, from the golang backend " + time.Now().String(),
		Status:  http.StatusOK,
	}

	// Encode the response as JSON and send it
	json.NewEncoder(w).Encode(response)
}

func main() {
	ctx := context.Background()

	config, client := init_github_connection(ctx)

	// Set up the route and handler
	http.HandleFunc("/config/hello_go", test_json_handler)
	http.HandleFunc("/config/get_contributors", get_contributors(ctx, client, config.Owner, config.Repo))
	http.HandleFunc("/config/listPRS", get_pr_list(ctx, client, config.Owner, config.Repo))

	//http.ListenAndServe(":8080", nil)

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
