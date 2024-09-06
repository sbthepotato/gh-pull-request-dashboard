package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	config, client := init_github_connection(ctx)

	// Set up the route and handler
	http.HandleFunc("/config/hello_go", hello_go)
	http.HandleFunc("/config/get_teams", get_teams(ctx, client, config.Owner))
	http.HandleFunc("/config/get_members", get_members(ctx, client, config.Owner))
	http.HandleFunc("/get_pr_list", get_pr_list(ctx, client, config.Owner, config.Repo))

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Could not start server: ", err.Error())
	}
}
