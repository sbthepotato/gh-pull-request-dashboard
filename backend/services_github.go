package main

import (
	"context"
	"log"
	"slices"
	"sync"

	"github.com/google/go-github/v64/github"
)

/*
get members and link them up to one of the active teams
*/
func gh_get_members(ctx context.Context, c *github.Client, owner string) ([]*CustomUser, error) {

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get members of org
	members, _, err := c.Organizations.ListMembers(ctx, owner, opt)
	if err != nil {
		return nil, err
	}

	users := make([]*CustomUser, 0)
	userTeams := make(map[string]*CustomTeam)

	teams, err := read_teams(true)
	if err != nil {
		return nil, err
	}

	// find team members of each team in org and add it to a map
	for _, team := range teams {

		if (team.ReviewEnabled == nil) || (!*team.ReviewEnabled) {
			continue
		}

		team_members, _, err := c.Teams.ListTeamMembersBySlug(ctx, owner, *team.Slug, nil)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		custom_user := new(CustomUser)
		custom_user.User = user
		custom_user.Team = userTeams[*user.Login]

		users = append(users, custom_user)
		userMap[*user.Login] = custom_user
	}

	write_users(userMap)

	return users, nil
}

/*
get list of all teams for a given organisation
*/
func gh_get_teams(ctx context.Context, c *github.Client, owner string) ([]*CustomTeam, error) {

	opt := &github.ListOptions{
		PerPage: 100,
	}

	teams, _, err := c.Teams.ListTeams(ctx, owner, opt)
	if err != nil {
		return nil, err
	}

	detailed_teams := make([]*CustomTeam, 0)
	teamMap := make(map[string]*CustomTeam)
	default_review := false
	default_order := 0

	for _, team := range teams {
		Custom_Team := new(CustomTeam)
		Custom_Team.Team = team
		Custom_Team.ReviewEnabled = &default_review
		Custom_Team.ReviewOrder = &default_order

		teamMap[*team.Slug] = Custom_Team
		detailed_teams = append(detailed_teams, Custom_Team)
	}

	write_teams(teamMap, false)

	return detailed_teams, nil
}

/*
get list of github pull requests and process them with review information
*/
func gh_get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) ([]*CustomPullRequest, error) {

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	pr_channel := make(chan *CustomPullRequest)

	users, err := read_users()
	if err != nil {
		return nil, err
	}

	teams, err := read_teams(true)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for _, pr := range gh_prs {
		if *pr.Draft {
			continue
		}

		wg.Add(1)
		go process_pr(pr_channel, &wg, ctx, c, owner, repo, *pr.Number, users, teams)
	}

	go func() {
		wg.Wait()
		close(pr_channel)
	}()

	prs := make([]*CustomPullRequest, 0)

	for processed_pr := range pr_channel {
		prs = append(prs, processed_pr)
	}

	return prs, nil

}

/*
Process a pull request into the pull request channel
*/
func process_pr(pr_channel chan<- *CustomPullRequest, wg *sync.WaitGroup, ctx context.Context, c *github.Client, owner string, repo string, pr_num int, users map[string]*CustomUser, teams map[string]*CustomTeam) {

	defer wg.Done()

	detailed_pr, _, err := c.PullRequests.Get(ctx, owner, repo, pr_num)
	if err != nil {
		log.Println("error fetching detailed pr info for pr", pr_num)
	}

	if *detailed_pr.Draft {
		*detailed_pr.State = "draft"
	}

	review := new(Review)
	review_overview := make([]Review, 0)
	user_review_list := make([]string, 0)
	status_requested := "REVIEW_REQUESTED"
	changes_requested := "Changes Requested"
	status_approved := "APPROVED"
	team_other := "OTHER"

	// first populate requested teams and users. any previous state doesn't matter if you're requested
	if detailed_pr.RequestedTeams != nil {
		for _, req_team := range detailed_pr.RequestedTeams {
			review.State = &status_requested
			review.Team = teams[*req_team.Slug]
			review.State = &status_requested

			review_overview = append(review_overview, *review)
		}
	}

	if detailed_pr.RequestedReviewers != nil {
		for _, req_review := range detailed_pr.RequestedReviewers {
			review.User = users[*req_review.Login]
			review.Team = teams[*review.User.Team.Slug]
			review.State = &status_requested

			review_overview = append(review_overview, *review)
			user_review_list = append(user_review_list, *req_review.Login)
		}
	}

	reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, *detailed_pr.Number, nil)
	if err != nil {
		log.Println("error fetching pull request reviews for pr", pr_num)
	}

	// loop in reverse because we're only interested in the most recent event
	for i := len(reviews) - 1; i >= 0; i-- {
		//for _, review := range reviews {
		gh_review := reviews[i]
		if (!slices.Contains(user_review_list, *gh_review.User.Login)) && (*detailed_pr.User.Login != *gh_review.User.Login) {
			review := new(Review)
			user_review_list = append(user_review_list, *gh_review.User.Login)
			review.User = users[*gh_review.User.Login]
			if review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = gh_review.State
			review_overview = append(review_overview, *review)
		}
	}

	custom_pr := new(CustomPullRequest)
	custom_pr.CreatedBy = users[*detailed_pr.User.Login]
	custom_pr.PullRequest = detailed_pr
	custom_pr.ReviewOverview = make([]*Review, 0)
	current_priority := 100
	approved_count := 0

	for _, custom_review := range review_overview {
		if (*custom_review.State != "DISMISSED") && (*custom_review.State != "COMMENTED") {
			review := new(Review)
			review.User = custom_review.User
			review.Team = custom_review.Team
			review.State = custom_review.State

			if *review.State == status_requested {
				if review.Team != nil && *review.Team.ReviewOrder < current_priority {
					current_priority = *review.Team.ReviewOrder
					custom_pr.Awaiting = review.Team.Name
				} else if review.Team == nil {
					custom_pr.Awaiting = &team_other
				}
			} else if *review.State == "CHANGES_REQUESTED" {
				custom_pr.Awaiting = &changes_requested
				current_priority = -1
			} else if *review.State == status_approved {
				approved_count++
			}
			custom_pr.ReviewOverview = append(custom_pr.ReviewOverview, review)
		}

		if custom_pr.Awaiting == nil && approved_count >= 1 {
			custom_pr.Awaiting = &status_approved
		}
	}

	pr_channel <- custom_pr

}
