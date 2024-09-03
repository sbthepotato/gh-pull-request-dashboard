package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/go-github/v64/github"
)

func get_contributors(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func get_reviews(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pr_number, err := strconv.Atoi(r.URL.Query().Get("pr"))
		if err != nil {
			log.Println("pr num ", r.URL.Query().Get("pr"))
			log.Fatalf("PR is not number: %e", err)
		}

		reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, pr_number, nil)
		if err != nil {
			log.Printf("Error fetching reviews for pull request #%d: %v", pr_number, err)
		}

		jsonData, err := json.Marshal(reviews)
		if err != nil {
			log.Fatalf("Error marshalling reviews to JSON: %e", err)
		}

		enableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
