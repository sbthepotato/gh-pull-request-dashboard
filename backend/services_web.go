package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/go-github/v67/github"
)

var mu sync.Mutex

var last_fetched_teams time.Time
var cached_teams []*CustomTeam

var last_fetched_members time.Time
var cached_members []*CustomUser

var last_fetched_prs time.Time
var cached_prs *PullRequestInfo

func hello_go(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}

func get_repos(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		repos, err := gh_get_repos(ctx, c, owner)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(repos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)

	}
}

func set_repos(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")

	mu.Lock()
	defer mu.Unlock()

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

	setRepos := make([]setRepo, 0)

	err = json.Unmarshal(body, &setRepos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := read_repos(false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// probably not the fastest way to do this but the list should never be huge...
	for _, setRepo := range setRepos {
		if setRepo.Enabled {
			for _, repo := range repos {
				if *repo.Name == setRepo.Name {
					repo.Enabled = &setRepo.Enabled
				}
			}
		}
	}

	err = write_repos(repos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Repo data saved successfully"))
}

func get_teams(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		var err error

		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if refresh == "y" {
			cached_teams, err = gh_get_teams(ctx, c, owner)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			last_fetched_teams = time.Now()

		} else if (currentTime.Sub(last_fetched_teams).Hours() < 1) || (len(cached_teams) == 0) {

			teams, err := read_teams(false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, team := range teams {
				cached_teams = append(cached_teams, team)
			}
		}

		jsonData, err := json.Marshal(cached_teams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}

func set_teams(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")

	mu.Lock()
	defer mu.Unlock()

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

	team_data := make([]SetTeam, 0)
	cached_teams = make([]*CustomTeam, 0)

	err = json.Unmarshal(body, &team_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team_map, err := read_teams(false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	active_team_map := make(map[string]*CustomTeam)

	for _, team := range team_data {
		*team_map[team.Slug].ReviewEnabled = team.ReviewEnabled
		*team_map[team.Slug].ReviewOrder = team.ReviewOrder

		if team.ReviewEnabled {
			active_team_map[team.Slug] = team_map[team.Slug]
		}

		updated_team := team_map[team.Slug]

		cached_teams = append(cached_teams, updated_team)
	}

	write_teams(active_team_map, true)
	write_teams(team_map, false)

	w.Write([]byte("Team data saved successfully"))
}

func get_members(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		var err error
		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if refresh == "y" {
			last_fetched_members = time.Now()
			cached_members, err = gh_get_members(ctx, c, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if currentTime.Sub(last_fetched_members).Hours() < 1 || (len(cached_members) == 0) {
			cached_members = make([]*CustomUser, 0)
			users, err := read_users()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, user := range users {
				cached_members = append(cached_members, user)
			}
		}

		jsonData, err := json.Marshal(cached_members)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		defer mu.Unlock()

		setHeaders(&w, "json")
		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if currentTime.Sub(last_fetched_prs).Minutes() > 5 || (refresh == "y" && currentTime.Sub(last_fetched_prs).Minutes() > 2) {
			cached_prs = new(PullRequestInfo)

			prs, err := gh_get_pr_list(ctx, c, owner, repo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cached_prs = prs
			last_fetched_prs = time.Now()
		}

		jsonData, err := json.Marshal(cached_prs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
