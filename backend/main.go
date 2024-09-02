package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func hello_go(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enableCors(&w)
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}

func main() {
	ctx := context.Background()

	config, client := init_github_connection(ctx)

	// Set up the route and handler
	http.HandleFunc("/config/hello_go", hello_go)
	http.HandleFunc("/config/get_contributors", get_contributors(ctx, client, config.Owner, config.Repo))
	http.HandleFunc("/get_pr_list", get_pr_list(ctx, client, config.Owner, config.Repo))

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
