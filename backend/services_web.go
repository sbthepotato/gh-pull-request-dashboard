package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v64/github"
)

var last_fetched_pr time.Time
var prs []Custom_Pull_Request

var last_fetched_teams time.Time
var cached_teams []*Custom_Team

func hello_go(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "json")
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}

func get_teams(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if refresh == "y" {
			log.Println("get new teams")
			cached_teams = gh_get_teams(ctx, c, owner)
		} else if (currentTime.Sub(last_fetched_teams).Hours() < 1) || (len(cached_teams) == 0) {
			team_map := read_teams()
			cached_teams = make([]*Custom_Team, 0)
			for _, team := range team_map {
				cached_teams = append(cached_teams, team)
			}
			log.Println("get teams from file")
		} else {
			log.Println("get cached teams")
		}

		jsonData, err := json.Marshal(cached_teams)
		if err != nil {
			log.Fatalln("Error marshalling teams to JSON: ", err)
		}

		setHeaders(&w, "json")
		w.Write(jsonData)
	}
}

func set_teams(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")

	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	team_data := make([]Set_Team, 0)
	cached_teams = make([]*Custom_Team, 0)

	err = json.Unmarshal(body, &team_data)
	if err != nil {
		log.Println("error unmarshaling team data: ", err.Error())
	}

	team_map := read_teams()

	for _, team := range team_data {
		*team_map[team.Slug].ReviewEnabled = team.ReviewEnabled
		*team_map[team.Slug].ReviewOrder = team.ReviewOrder

		updated_team := team_map[team.Slug]

		cached_teams = append(cached_teams, updated_team)
	}

	write_teams(team_map)

	setHeaders(&w, "text")
	w.Write([]byte("Team data saved successfully"))
}

func get_members(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jsonData, err := json.Marshal(gh_get_members(ctx, c, owner))
		if err != nil {
			log.Fatalf("Error marshalling members to JSON: %e", err)
		}

		setHeaders(&w, "json")
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

		setHeaders(&w, "json")
		w.Write(jsonData)
	}
}
