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

	var wg sync.WaitGroup
	user_channel := make(chan *CustomUser)

	// go through all org members to get extended user info, also add team info
	for _, member := range members {
		wg.Add(1)
		go process_member(user_channel, &wg, ctx, c, *member.Login, userTeams)
	}

	go func() {
		wg.Wait()
		close(user_channel)
	}()

	userMap := make(map[string]*CustomUser)
	for processed_user := range user_channel {
		userMap[*processed_user.Login] = processed_user
		users = append(users, processed_user)
	}

	write_users(userMap)

	return users, nil
}

/*
process a member into the member channel
*/
func process_member(user_channel chan<- *CustomUser, wg *sync.WaitGroup, ctx context.Context, c *github.Client, login string, teams map[string]*CustomTeam) {
	defer wg.Done()

	user, _, err := c.Users.Get(ctx, login)
	if err != nil {
		log.Println("error fetching user", login, err)
	}

	custom_user := new(CustomUser)
	custom_user.User = user
	custom_user.Team = teams[*user.Login]

	user_channel <- custom_user
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

	go func() {
		wg.Wait()
		close(PrChannel)
	}()

	prs := make([]*CustomPullRequest, idx)

	for processed_pr := range PrChannel {
		prs[*processed_pr.Index] = processed_pr
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

	detailed_pr, _, err := c.PullRequests.Get(ctx, owner, repo, *pr.Number)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching detailed pr info for pr " + strconv.Itoa(*pr.Number)
		log.Println(ErrorMessage, err.Error())
	}

	if *detailed_pr.Draft {
		*detailed_pr.State = "draft"
	}

	customPr.PullRequest = detailed_pr

	review := new(Review)
	review_overview := make([]Review, 0)
	userReviewList := make([]string, 0)
	status_requested := "REVIEW_REQUESTED"
	changes_requested := "Changes Requested"
	status_approved := "APPROVED"
	team_other := "other"

	if teams == nil {
		team_other = "review"
	}

	// first populate requested teams and users. any previous state doesn't matter if you're requested
	if detailed_pr.RequestedTeams != nil {
		for _, req_team := range detailed_pr.RequestedTeams {
			review.State = &status_requested

			// if the team map isn't available, just use what we have
			if val, ok := teams[*req_team.Slug]; ok {
				review.Team = val
			} else {
				custom_team := new(CustomTeam)
				custom_team.Team = req_team
				review.Team = custom_team
			}

			review.State = &status_requested

			review_overview = append(review_overview, *review)
		}
	}

	if detailed_pr.RequestedReviewers != nil {
		for _, req_review := range detailed_pr.RequestedReviewers {

			// if the user map isn't available, just use what we have
			if val, ok := users[*req_review.Login]; ok {
				review.User = val
			} else {
				custom_user := new(CustomUser)
				custom_user.User = req_review
				review.User = custom_user
			}

			if review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = &status_requested

			review_overview = append(review_overview, *review)
			userReviewList = append(userReviewList, *req_review.Login)
		}
	}

	reviews, _, err := c.PullRequests.ListReviews(ctx, owner, repo, *detailed_pr.Number, nil)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching pull request reviews for pr " + strconv.Itoa(*pr.Number)
		team_other = "error"
		customPr.Awaiting = &team_other

		log.Println(ErrorMessage, err.Error())
	}

	// loop in reverse because we're only interested in the most recent event
	for i := len(reviews) - 1; i >= 0; i-- {
		//for _, review := range reviews {
		gh_review := reviews[i]
		if (!slices.Contains(userReviewList, *gh_review.User.Login)) &&
			(*detailed_pr.User.Login != *gh_review.User.Login) &&
			(*gh_review.State != "COMMENTED") {

			review := new(Review)
			userReviewList = append(userReviewList, *gh_review.User.Login)
			if val, ok := users[*gh_review.User.Login]; ok {
				review.User = val
			} else {
				custom_user := new(CustomUser)
				custom_user.User = gh_review.User
				review.User = custom_user
			}

			if review.User != nil && review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = gh_review.State
			review_overview = append(review_overview, *review)

		}
	}

	customPr.CreatedBy = users[*detailed_pr.User.Login]
	customPr.ReviewOverview = make([]*Review, 0)
	customPr.Index = &idx
	currentPriority := 100
	approvedCount := 0

	for _, customReview := range review_overview {
		if *customReview.State != "DISMISSED" {

			review := new(Review)
			review.User = customReview.User
			review.Team = customReview.Team
			review.State = customReview.State

			if *review.State == status_requested {
				if review.Team != nil {
					if review.Team.ReviewOrder != nil &&
						*review.Team.ReviewOrder < currentPriority &&
						*customPr.Awaiting != "error" {

						currentPriority = *review.Team.ReviewOrder
						customPr.Awaiting = review.Team.Name

					} else if customPr.Awaiting == nil {
						customPr.Awaiting = &team_other
					}
				} else if review.Team == nil {
					customPr.Awaiting = &team_other
				}

			} else if *review.State == "CHANGES_REQUESTED" {
				customPr.Awaiting = &changes_requested
				currentPriority = -1
			} else if *review.State == status_approved {
				approvedCount++
			}
			customPr.ReviewOverview = append(customPr.ReviewOverview, review)
		}

		if customPr.Awaiting == nil && approvedCount >= 1 {
			customPr.Awaiting = &status_approved
		}
	}

	customPr.ErrorMessage = &ErrorMessage
	customPr.ErrorText = &ErrorText

	PrChannel <- customPr

}
