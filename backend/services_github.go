package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
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

		// Create a response object
		response := Response{
			Message: string(jsonData),
			Status:  http.StatusOK,
		}

		// Encode the response as JSON and send it
		json.NewEncoder(w).Encode(response)
	}
}

func get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		opts := &github.PullRequestListOptions{
			State:       "open",
			ListOptions: github.ListOptions{PerPage: 30},
		}

		var allPRs []*github.PullRequest

		for {
			prs, resp, err := c.PullRequests.List(ctx, owner, repo, opts)
			if err != nil {
				log.Fatalf("Error fetching pull requests: %e", err)
			}

			allPRs = append(allPRs, prs...)

			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

		// Print pull request information along with requested reviewers and their approval status
		for _, pr := range allPRs {
			// Fetch detailed information for each pull request, including reviewers
			detailedPR, _, err := c.PullRequests.Get(ctx, owner, repo, pr.GetNumber())
			if err != nil {
				log.Printf("Error fetching pull request #%d details: %v", pr.GetNumber(), err)
				continue
			}

			fmt.Printf("#%d %s\n", pr.GetNumber(), pr.GetTitle())

			// Fetch the reviews for the pull request
			reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, pr.GetNumber(), nil)
			if err != nil {
				log.Printf("Error fetching reviews for pull request #%d: %v", pr.GetNumber(), err)
				continue
			}

			// Map to store the latest review status for each reviewer
			reviewerStatus := make(map[string]string)

			// Process each review and update the status for each reviewer
			for _, review := range reviews {
				reviewer := review.GetUser().GetLogin()
				state := review.GetState()
				reviewerStatus[reviewer] = state
			}

			// Print requested reviewers and their approval status
			if len(detailedPR.RequestedReviewers) > 0 {
				fmt.Println("Requested Reviewers:")
				for _, reviewer := range detailedPR.RequestedReviewers {
					reviewerLogin := reviewer.GetLogin()
					status, ok := reviewerStatus[reviewerLogin]
					if ok {
						switch status {
						case "APPROVED":
							fmt.Printf("- %s: Approved\n", reviewerLogin)
						case "CHANGES_REQUESTED":
							fmt.Printf("- %s: Changes Requested\n", reviewerLogin)
						case "COMMENTED":
							fmt.Printf("- %s: Commented\n", reviewerLogin)
						default:
							fmt.Printf("- %s: Reviewed (State: %s)\n", reviewerLogin, status)
						}
					} else {
						fmt.Printf("- %s: Pending Review\n", reviewerLogin)
					}
				}
			} else {
				fmt.Println("No requested reviewers.")
			}

			fmt.Println() // Blank line between PRs for readability
		}
	}
}
