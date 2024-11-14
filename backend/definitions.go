package main

import "github.com/google/go-github/v66/github"

/*
core application settings
*/
type Config struct {
	Token string `json:"token,omitempty"`
	Owner string `json:"owner,omitempty"`
	Repo  string `json:"repo,omitempty"`
}

/*
Overview of a Pull Request Review
*/
type Review struct {
	User  *CustomUser `json:"user,omitempty"`
	Team  *CustomTeam `json:"team,omitempty"`
	State *string     `json:"state,omitempty"`
}

/*
Pull Request with extra fields for custom objects
*/
type CustomPullRequest struct {
	*github.PullRequest
	CreatedBy      *CustomUser `json:"created_by,omitempty"`
	ReviewOverview []*Review   `json:"review_overview,omitempty"`
	Awaiting       *string     `json:"awaiting,omitempty"`
	ErrorMessage   *string     `json:"error_message,omitempty"`
	ErrorText      *string     `json:"error_text,omitempty"`
	Index          *int        `json:"-"`
}

/*
User with a Custom Team attached
*/
type CustomUser struct {
	*github.User
	Team *CustomTeam `json:"team,omitempty"`
}

/*
Team with review info
*/
type CustomTeam struct {
	*github.Team
	ReviewEnabled *bool `json:"review_enabled,omitempty"`
	ReviewOrder   *int  `json:"review_order,omitempty"`
}

/*
the POST from frontend to set team info
*/
type SetTeam struct {
	Slug          string `json:"slug,omitempty"`
	ReviewEnabled bool   `json:"review_enabled,omitempty"`
	ReviewOrder   int    `json:"review_order,omitempty"`
}

/*
pull request list with accompanying information for list
*/
type PullRequestInfo struct {
	PullRequests []*CustomPullRequest `json:"pull_requests,omitempty"`
	ReviewTeams  []*CustomTeam        `json:"review_teams,omitempty"`
}
