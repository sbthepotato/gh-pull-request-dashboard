package main

import (
	"context"
	"log"
	"slices"
	"sort"
	"strconv"
	"sync"

	"github.com/google/go-github/v66/github"
)

/*
get all repositories for the currently set org
*/
func gh_get_repos(ctx context.Context, c *github.Client, owner string) ([]*CustomRepo, error) {

	opt := &github.RepositoryListByOrgOptions{
		Sort:        "full_name",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var all_repos []*github.Repository

	for {
		repos, resp, err := c.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		all_repos = append(all_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}

	var custom_repos []*CustomRepo

	for _, repo := range all_repos {
		custom_repo := new(CustomRepo)

		custom_repo.Repository = repo
		custom_repos = append(custom_repos, custom_repo)

	}

	return custom_repos, nil
}

/*
get members and link them up to one of the active teams
*/
func gh_get_members(ctx context.Context, c *github.Client, owner string) ([]*CustomUser, error) {

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var all_users []*github.User

	for {
		users, resp, err := c.Organizations.ListMembers(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		all_users = append(all_users, users...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	teams, err := read_teams(true)
	if err != nil {
		return nil, err
	}

	user_teams := make(map[string]*CustomTeam)

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
			user_teams[*team_member.Login] = team
		}
	}

	var wg sync.WaitGroup
	user_channel := make(chan *CustomUser)

	// go through all org members to get extended user info, also add team info
	for _, member := range all_users {
		wg.Add(1)
		go process_user(user_channel, &wg, ctx, c, *member.Login, user_teams)
	}

	go func() {
		wg.Wait()
		close(user_channel)
	}()

	user_map := make(map[string]*CustomUser)
	custom_users := make([]*CustomUser, 0)

	for processed_user := range user_channel {
		user_map[*processed_user.Login] = processed_user
		custom_users = append(custom_users, processed_user)
	}

	write_users(user_map)

	return custom_users, nil
}

/*
process a member into the member channel
*/
func process_user(user_channel chan<- *CustomUser, wg *sync.WaitGroup, ctx context.Context, c *github.Client, login string, teams map[string]*CustomTeam) {
	defer wg.Done()

	user, _, err := c.Users.Get(ctx, login)
	if err != nil {
		log.Println("error fetching user", login, err)
	}

	customUser := new(CustomUser)
	customUser.User = user
	customUser.Team = teams[*user.Login]

	user_channel <- customUser
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
func gh_get_pr_list(ctx context.Context, c *github.Client, owner string, repo string) (*PullRequestInfo, error) {

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	gh_prs, _, err := c.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	users, err := read_users()
	if err != nil {
		return nil, err
	}

	teams, err := read_teams(true)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	idx := 0
	PrChannel := make(chan *CustomPullRequest)

	for _, pr := range gh_prs {
		if *pr.Draft {
			continue
		}

		wg.Add(1)
		go process_pr(PrChannel, &wg, ctx, c, owner, repo, pr, users, teams, idx)
		idx++ // manual index as we are skipping draft
	}

	result := new(PullRequestInfo)

	// sort the list of teams for the aggregation banner
	if teams != nil {
		slugs := make([]string, 0, len(teams))

		for slug := range teams {
			slugs = append(slugs, slug)
		}

		sort.SliceStable(slugs, func(i, j int) bool {
			return *teams[slugs[i]].ReviewOrder < *teams[slugs[j]].ReviewOrder
		})

		sorted_teams := make([]*CustomTeam, 0)
		for _, slug := range slugs {
			sorted_teams = append(sorted_teams, teams[slug])
		}

		result.ReviewTeams = sorted_teams
	}

	if users != nil {
		userList := make([]*CustomUser, 0)

		for _, user := range users {
			userList = append(userList, user)
		}

		result.Users = userList
	}

	go func() {
		wg.Wait()
		close(PrChannel)
	}()

	prs := make([]*CustomPullRequest, idx)

	for processed_pr := range PrChannel {
		prs[*processed_pr.Index] = processed_pr
	}

	result.PullRequests = prs

	return result, nil

}

/*
Process a pull request into the pull request channel
*/
func process_pr(PrChannel chan<- *CustomPullRequest, wg *sync.WaitGroup, ctx context.Context, c *github.Client, owner string, repo string, pr *github.PullRequest, users map[string]*CustomUser, teams map[string]*CustomTeam, idx int) {

	defer wg.Done()

	customPr := new(CustomPullRequest)
	var ErrorMessage string
	var ErrorText string

	detailedPr, _, err := c.PullRequests.Get(ctx, owner, repo, *pr.Number)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching detailed pr info for pr " + strconv.Itoa(*pr.Number)
		log.Println(ErrorMessage, err.Error())
	}

	if *detailedPr.Draft {
		*detailedPr.State = "draft"
	}

	customPr.PullRequest = detailedPr

	review := new(Review)
	reviewOverview := make([]Review, 0)
	userReviewList := make([]string, 0)
	statusRequested := "REVIEW_REQUESTED"
	changesRequested := "Changes Requested"
	statusApproved := "APPROVED"
	teamOther := "other"

	if teams == nil {
		teamOther = "review"
	}

	// first populate requested teams and users. any previous state doesn't matter if you're requested
	if detailedPr.RequestedTeams != nil {
		for _, requestedTeam := range detailedPr.RequestedTeams {
			review.State = &statusRequested

			// if the team map isn't available, just use what we have
			if val, ok := teams[*requestedTeam.Slug]; ok {
				review.Team = val
			} else {
				custom_team := new(CustomTeam)
				custom_team.Team = requestedTeam
				review.Team = custom_team
			}

			review.State = &statusRequested

			reviewOverview = append(reviewOverview, *review)
		}
	}

	if detailedPr.RequestedReviewers != nil {
		for _, requestedUser := range detailedPr.RequestedReviewers {

			// if the user map isn't available, just use what we have
			if val, ok := users[*requestedUser.Login]; ok {
				review.User = val
			} else {
				customUser := new(CustomUser)
				customUser.User = requestedUser
				review.User = customUser
			}

			if review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = &statusRequested

			reviewOverview = append(reviewOverview, *review)
			userReviewList = append(userReviewList, *requestedUser.Login)
		}
	}

	reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, *detailedPr.Number, nil)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching pull request reviews for pr " + strconv.Itoa(*pr.Number)
		teamOther = "error"
		customPr.Awaiting = &teamOther

		log.Println(ErrorMessage, err.Error())
	}

	// loop in reverse because we're only interested in the most recent event
	for i := len(reviews) - 1; i >= 0; i-- {
		//for _, review := range reviews {
		ghReview := reviews[i]
		if (!slices.Contains(userReviewList, *ghReview.User.Login)) &&
			(*detailedPr.User.Login != *ghReview.User.Login) &&
			(*ghReview.State != "COMMENTED") {

			review := new(Review)
			userReviewList = append(userReviewList, *ghReview.User.Login)
			if val, ok := users[*ghReview.User.Login]; ok {
				review.User = val
			} else {
				customUser := new(CustomUser)
				customUser.User = ghReview.User
				review.User = customUser
			}

			if review.User != nil && review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = ghReview.State
			reviewOverview = append(reviewOverview, *review)

		}
	}

	customPr.CreatedBy = users[*detailedPr.User.Login]
	customPr.ReviewOverview = make([]*Review, 0)
	customPr.Index = &idx
	currentPriority := 100
	unassigned := false
	approvedCount := 0

	for _, customReview := range reviewOverview {
		if *customReview.State != "DISMISSED" {

			review := new(Review)
			review.User = customReview.User
			review.Team = customReview.Team
			review.State = customReview.State

			if *review.State == statusRequested {
				if review.Team != nil {

					if review.Team.ReviewOrder != nil &&
						*review.Team.ReviewOrder < currentPriority &&
						(customPr.Awaiting != nil &&
							*customPr.Awaiting != "error" ||
							customPr.Awaiting == nil) {

						currentPriority = *review.Team.ReviewOrder
						customPr.Awaiting = review.Team.Name

						if review.User == nil {
							unassigned = true
							customPr.Unassigned = &unassigned
						} else {
							unassigned = false
							customPr.Unassigned = &unassigned
						}

					} else if review.Team.ReviewOrder != nil &&
						*review.Team.ReviewOrder == currentPriority &&
						review.User != nil {

						unassigned = false
						customPr.Unassigned = &unassigned

					} else if customPr.Awaiting == nil {
						customPr.Awaiting = &teamOther
					}

				} else if review.Team == nil {
					customPr.Awaiting = &teamOther
				}

			} else if *review.State == "CHANGES_REQUESTED" {
				customPr.Awaiting = &changesRequested
				currentPriority = -1
				unassigned = false
				customPr.Unassigned = &unassigned

			} else if *review.State == statusApproved {
				approvedCount++
			}
			customPr.ReviewOverview = append(customPr.ReviewOverview, review)

		}

		if customPr.Awaiting == nil && approvedCount >= 1 {
			customPr.Awaiting = &statusApproved
			unassigned = false
			customPr.Unassigned = &unassigned
		}
	}

	if ErrorMessage != "" {
		ErrorMessage = ErrorMessage + ". See the console on the server or the errorText field in the network tab for more information"
		customPr.ErrorMessage = &ErrorMessage
		customPr.ErrorText = &ErrorText
	}

	PrChannel <- customPr

}
