package main

import (
	"context"
	"log"
	"slices"
	"sync"

	"github.com/google/go-github/v64/github"
)

func gh_get_members(ctx context.Context, c *github.Client, owner string) []*CustomUser {

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get members of org
	members, _, err := c.Organizations.ListMembers(ctx, owner, opt)
	if err != nil {
		log.Fatalf("Error fetching contributors: %e", err)
	}

	users := make([]*CustomUser, 0)
	userTeams := make(map[string]*CustomTeam)

	teams := read_teams(true)

	// find team members of each team in org and add it to a map
	for _, team := range teams {

		if (team.ReviewEnabled == nil) || (!*team.ReviewEnabled) {
			continue
		}

		team_members, _, err := c.Teams.ListTeamMembersBySlug(ctx, owner, *team.Slug, nil)
		if err != nil {
			log.Fatalf("error fetching team members: %e", err)
		}

		for _, team_member := range team_members {
			userTeams[*team_member.Login] = team
		}
	}

	userMap := make(map[string]*CustomUser)

	// go through all org members to get extended user info, also add team info
	for _, member := range members {
		user, _, err := c.Users.Get(ctx, *member.Login)
		if err != nil {
			log.Fatalf("Error fetching user: %e", err)
		}

		custom_user := new(CustomUser)
		custom_user.User = user
		custom_user.Team = userTeams[*user.Login]

		users = append(users, custom_user)
		userMap[*user.Login] = custom_user
	}

	write_users(userMap)

	return users
}

func gh_get_teams(ctx context.Context, c *github.Client, owner string) []*CustomTeam {

	opt := &github.ListOptions{
		PerPage: 100,
	}

	teams, _, err := c.Teams.ListTeams(ctx, owner, opt)
	if err != nil {
		log.Fatalln("Error fetching teams: ", err.Error())
	}

	detailed_teams := make([]*CustomTeam, 0)
	teamMap := make(map[string]*CustomTeam)
	default_review := false
	default_order := 0

	for _, team := range teams {
		/* detailed team functionality temporarily removed as it gives us nothing useful (thanks github)
		 detailed_team, _, err := c.Teams.GetTeamBySlug(ctx, owner, *team.Slug)
		if err != nil {
			log.Fatal("Error fetching detailed team info: ", err.Error())
		} */

		Custom_Team := new(CustomTeam)
		Custom_Team.Team = team
		Custom_Team.ReviewEnabled = &default_review
		Custom_Team.ReviewOrder = &default_order

		teamMap[*team.Slug] = Custom_Team
		detailed_teams = append(detailed_teams, Custom_Team)
	}

	write_teams(teamMap, false)

	return detailed_teams
}

func gh_get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) []*CustomPullRequest {

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		log.Fatalf("Error fetching Pull Requests: %e", err)
	}

	prs := make([]*CustomPullRequest, 0)
	users := read_users()
	teams := read_teams(true)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, pr := range gh_prs {

		if *pr.Draft {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			if *pr.Draft {
				*pr.State = "draft"
			}

			review := new(Review)
			review_overview := make([]Review, 0)
			user_review_list := make([]string, 0)
			status_requested := "REVIEW_REQUESTED"
			changes_requested := "Changes Requested"

			// first populate requested teams and users. any previous state doesn't matter if you're requested
			if pr.RequestedTeams != nil {
				for _, req_team := range pr.RequestedTeams {
					review.State = &status_requested
					review.Team = teams[*req_team.Slug]
					review.State = &status_requested

					review_overview = append(review_overview, *review)
				}
			}

			if pr.RequestedReviewers != nil {
				for _, req_review := range pr.RequestedReviewers {
					review.User = users[*req_review.Login]
					review.Team = teams[*review.User.Team.Slug]
					review.State = &status_requested

					review_overview = append(review_overview, *review)
					user_review_list = append(user_review_list, *req_review.Login)
				}
			}

			reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, *pr.Number, nil)
			if err != nil {
				log.Fatalf("error fetching pull request reviews")
			}

			// loop in reverse because we're only interested in the most recent event
			for i := len(reviews) - 1; i >= 0; i-- {
				//for _, review := range reviews {
				gh_review := reviews[i]
				if (!slices.Contains(user_review_list, *gh_review.User.Login)) && (*pr.User.Login != *gh_review.User.Login) {
					review := new(Review)
					user_review_list = append(user_review_list, *gh_review.User.Login)
					review.User = users[*gh_review.User.Login]
					review.Team = teams[*review.User.Team.Slug]
					review.State = gh_review.State
					review_overview = append(review_overview, *review)
				}
			}

			custom_pr := new(CustomPullRequest)
			custom_pr.CreatedBy = users[*pr.User.Login]
			custom_pr.PullRequest = pr
			custom_pr.ReviewOverview = make([]*Review, 0)
			current_priority := 100

			for _, custom_review := range review_overview {
				if (*custom_review.State != "DISMISSED") && (*custom_review.State != "COMMENTED") {
					review := new(Review)
					review.User = custom_review.User
					review.Team = custom_review.Team
					review.State = custom_review.State

					if *review.State == status_requested {
						if *review.Team.ReviewOrder < current_priority {
							current_priority = *review.Team.ReviewOrder
							custom_pr.Awaiting = review.Team.Name
						}
					} else if *review.State == "CHANGES_REQUESTED" {
						custom_pr.Awaiting = &changes_requested
						current_priority = -1
					}

					custom_pr.ReviewOverview = append(custom_pr.ReviewOverview, review)
				}
			}

			mu.Lock()
			prs = append(prs, custom_pr)
			mu.Unlock()

		}()

		wg.Wait()
	}

	return prs

}
