package main

import (
	"context"
	"log"
	"sync"

	"github.com/google/go-github/v64/github"
)

func gh_get_members(ctx context.Context, c *github.Client, owner string) []*Custom_User {

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 99},
	}

	// get members of org
	members, _, err := c.Organizations.ListMembers(ctx, owner, opt)
	if err != nil {
		log.Fatalf("Error fetching contributors: %e", err)
	}

	users := make([]*Custom_User, 0)
	userTeams := make(map[string]string)

	teams := gh_get_teams(ctx, c, owner)

	// find team members of each team in org and add it to a map
	for _, team := range teams {

		team_members, _, err := c.Teams.ListTeamMembersBySlug(ctx, owner, *team.Slug, nil)
		if err != nil {
			log.Fatalf("error fetching team members: %e", err)
		}

		for _, team_member := range team_members {
			userTeams[*team_member.Login] = *team.Name
		}
	}

	userMap := make(map[string]*Custom_User)

	// go through all org members to get extended user info, also add team info
	for _, member := range members {
		user, _, err := c.Users.Get(ctx, *member.Login)
		if err != nil {
			log.Fatalf("Error fetching user: %e", err)
		}

		custom_user := new(Custom_User)
		custom_user.User = *user
		custom_user.Team_Name = userTeams[*user.Login]

		users = append(users, custom_user)

		userMap[*user.Login] = custom_user
	}

	write_users(userMap)

	return users
}

func gh_get_teams(ctx context.Context, c *github.Client, owner string) []*Custom_Team {

	opt := &github.ListOptions{
		PerPage: 99,
	}

	teams, _, err := c.Teams.ListTeams(ctx, owner, opt)
	if err != nil {
		log.Fatalln("Error fetching teams: ", err.Error())
	}

	detailed_teams := make([]*Custom_Team, 0)
	teamMap := make(map[string]*Custom_Team)

	for _, team := range teams {
		detailed_team, _, err := c.Teams.GetTeamBySlug(ctx, owner, *team.Slug)
		if err != nil {
			log.Fatal("Error fetching detailed team info: ", err.Error())
		}

		Custom_Team := new(Custom_Team)
		Custom_Team.Team = *detailed_team

		teamMap[*detailed_team.Slug] = Custom_Team
		detailed_teams = append(detailed_teams, Custom_Team)
	}

	write_teams(teamMap)

	return detailed_teams
}

func gh_get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) []Custom_Pull_Request {

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 99},
	}

	gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		log.Fatalf("Error fetching Pull Requests: %e", err)
	}

	prs := make([]Custom_Pull_Request, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, pr := range gh_prs {

		if *pr.Draft {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			/* currently draft PRs are skipped, if they are to be included then this is needed
			if *pr.Draft {
				*pr.State = "draft"
			}
			*/

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
