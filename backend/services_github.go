package main

import (
	"context"
	"log"
	"sync"

	"github.com/google/go-github/v64/github"
)

func gh_get_contributors(ctx context.Context, c *github.Client, owner string, repo string) []*github.Contributor {

	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 5},
	}

	contributors, _, err := c.Repositories.ListContributors(ctx, owner, repo, opt)
	if err != nil {
		log.Fatalf("Error fetching contributors: %e", err)
	}

	return contributors

}

func gh_get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) []Custom_Pull_Request {

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 30},
	}

	gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		log.Fatalf("Error fetching Pull Requests: %e", err)
	}

	prs := make([]Custom_Pull_Request, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, pr := range gh_prs {

		wg.Add(1)

		go func() {
			defer wg.Done()

			if *pr.Draft {
				*pr.State = "draft"
			}

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

			mu.Lock()
			prs = append(prs, *custom_pr)
			mu.Unlock()

		}()

		wg.Wait()
	}

	return prs

}
