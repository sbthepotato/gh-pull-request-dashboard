package main

import (
	"fmt"
	"context"
	"log"
	"os"

	"github.com/google/go-github/v64/github"
	"golang.org/x/oauth2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
	}
}


func main() {

	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
			log.Fatal("AUTH_TOKEN not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var allPRs []*github.PullRequest

	owner := "google"
	repo := "go-github"

	for {
		prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			log.Fatalf("Error fetching pull requests: %v", err)
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
		detailedPR, _, err := client.PullRequests.Get(ctx, owner, repo, pr.GetNumber())
		if err != nil {
			log.Printf("Error fetching pull request #%d details: %v", pr.GetNumber(), err)
			continue
		}

		fmt.Printf("#%d %s\n", pr.GetNumber(), pr.GetTitle())

		// Fetch the reviews for the pull request
		reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, pr.GetNumber(), nil)
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