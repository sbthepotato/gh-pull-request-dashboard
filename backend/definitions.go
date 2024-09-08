package main

import "github.com/google/go-github/v64/github"

type Config struct {
	Token string `json:"token,omitempty"`
	Owner string `json:"owner,omitempty"`
	Repo  string `json:"repo,omitempty"`
}

type ReviewAggregation struct {
	github.Team
	PendingReviewCount *int `json:"pending_review_count,omitempty"`
}

type Review struct {
	User  *CustomUser `json:"user,omitempty"`
	Team  *CustomTeam `json:"team,omitempty"`
	State *string     `json:"state,omitempty"`
}

type CustomPullRequest struct {
	github.PullRequest
	CreatedBy      *CustomUser `json:"created_by,omitempty"`
	ReviewOverview []*Review   `json:"review_overview,omitempty"`
	Awaiting       *string     `json:"awaiting,omitempty"`
}

type CustomUser struct {
	*github.User
	Team *CustomTeam `json:"team,omitempty"`
}

type CustomTeam struct {
	github.Team
	ReviewEnabled *bool `json:"review_enabled,omitempty"`
	ReviewOrder   *int  `json:"review_order,omitempty"`
}

type SetTeam struct {
	Slug          string `json:"slug,omitempty"`
	ReviewEnabled bool   `json:"review_enabled,omitempty"`
	ReviewOrder   int    `json:"review_order,omitempty"`
}
