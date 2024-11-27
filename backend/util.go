package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
)

func setHeaders(w *http.ResponseWriter, content_type string) {

	if content_type == "text" {
		(*w).Header().Set("Content-Type", "text/plain")
	} else if content_type == "json" {
		(*w).Header().Set("Content-Type", "application/json")
	}
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func enable_cors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Serve the actual request
		handler.ServeHTTP(w, r)
	})
}

func load_config() *Config {
	content, err := os.ReadFile("./db/config.json")
	if err != nil {
		log.Fatal("Error when opening config: ", err)
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal of config: ", err)
	}

	return &payload
}

func init_github_connection(ctx context.Context) (*Config, *github.Client) {
	config := load_config()

	authToken := config.Token

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	config.Token = ""

	return config, client

}
