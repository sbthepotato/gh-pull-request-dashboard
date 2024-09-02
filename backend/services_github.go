package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/v64/github"
)

func get_contributors(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")

		opt := &github.ListContributorsOptions{
			ListOptions: github.ListOptions{PerPage: 5},
		}

		contributors, _, err := c.Repositories.ListContributors(ctx, owner, repo, opt)
		if err != nil {
			log.Fatalf("Error fetching contributors: %e", err)
		}

		jsonData, err := json.Marshal(contributors)
		if err != nil {
			log.Fatalf("Error marshalling contributors to JSON: %e", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")

		opts := &github.PullRequestListOptions{
			State:       "open",
			ListOptions: github.ListOptions{PerPage: 30},
		}

		prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching Pull Requests: %e", err)
		}

		jsonData, err := json.Marshal(prs)
		if err != nil {
			log.Fatalf("Error marshalling Pull Requests to JSON: %e", err)
		}

		w.Write(jsonData)
	}
}
