package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v64/github"
)

var last_fetched_pr time.Time
var prs []Custom_Pull_Request

var teams []*github.Team

func hello_go(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}

func get_teams(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if len(teams) < 1 {
			teams = gh_get_teams(ctx, c, owner)
			log.Println("get new teams")
		} else {
			log.Println("use cached teams")
		}

		jsonData, err := json.Marshal(teams)
		if err != nil {
			log.Fatalln("Error marshalling teams to JSON: ", err)
		}

		setHeaders(&w)
		w.Write(jsonData)
	}
}

func get_members(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonData, err := json.Marshal(gh_get_members(ctx, c, owner))
		if err != nil {
			log.Fatalf("Error marshalling members to JSON: %e", err)
		}

		setHeaders(&w)
		w.Write(jsonData)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		currentTime := time.Now()

		if currentTime.Sub(last_fetched_pr).Minutes() > 30 {
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

		setHeaders(&w)
		w.Write(jsonData)
	}
}
