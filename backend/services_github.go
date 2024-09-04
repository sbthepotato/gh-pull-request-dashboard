package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

		gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching Pull Requests: %e", err)
		}

		prs := make([]Custom_Pull_Request, 0)

		for _, pr := range gh_prs {

			review_overview := make(map[string]string, 0)

			// first populate requested teams and users. any previous state doesn't matter if you're requested
			if pr.RequestedTeams != nil {
				for _, req_team := range pr.RequestedTeams {
					review_overview[*req_team.Name] = "REQUESTED"
				}
			}

			if pr.RequestedReviewers != nil {
				for _, req_review := range pr.RequestedReviewers {
					review_overview[*req_review.Login] = "REQUESTED"
				}
			}

			reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, *pr.Number, nil)
			if err != nil {
				log.Fatalf("error fetching pull request reviews")
			}

			// loop in reverse because we're only interested in the most recent event
			for i := len(reviews) - 1; i >= 0; i-- {
				//for _, review := range reviews {
				review := reviews[i]
				_, exists := review_overview[*review.User.Login]
				if (!exists) && (*pr.User.Login != *review.User.Login) {
					review_overview[*review.User.Login] = *review.State
				}
			}

			custom_pr := new(Custom_Pull_Request)
			custom_pr.PullRequest = *pr
			custom_pr.Review_Overview = make([]Review_Overview, 0)

			for user, state := range review_overview {
				if (state != "DISMISSED") && (state != "COMMENTED") {
					review_overview := new(Review_Overview)
					review_overview.User = user
					review_overview.State = state

					custom_pr.Review_Overview = append(custom_pr.Review_Overview, *review_overview)
				}
			}

			prs = append(prs, *custom_pr)

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
