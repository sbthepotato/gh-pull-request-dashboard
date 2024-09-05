package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v64/github"
)

type Review_Overview struct {
	User  string
	State string
}

type Custom_Pull_Request struct {
	github.PullRequest
	Review_Overview []Review_Overview
}

var last_fetched_pr time.Time
var prs []Custom_Pull_Request

func get_contributors(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonData, err := json.Marshal(gh_get_contributors(ctx, c, owner, repo))
		if err != nil {
			log.Fatalf("Error marshalling contributors to JSON: %e", err)
		}

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		currentTime := time.Now()

		if currentTime.Sub(last_fetched_pr).Minutes() > 10 {
			log.Print("get new")
			prs = gh_get_pr_list(ctx, c, owner, repo)
			last_fetched_pr = time.Now()
		} else {
			log.Print("use cached")
		}

		jsonData, err := json.Marshal(prs)
		if err != nil {
			log.Fatalf("Error marshalling Pull Requests to JSON: %e", err)
		}

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
