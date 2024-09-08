package main

import "github.com/google/go-github/v64/github"

type Config struct {
	Token string `json:"token,omitempty"`
	Owner string `json:"owner,omitempty"`
	Repo  string `json:"repo,omitempty"`
}

type Review_aggregation struct {
	github.Team
	PendingReviewCount *int `json:"pending_review_count,omitempty"`
}

type Review_Overview struct {
	User  *string `json:"user,omitempty"`
	State *string `json:"state,omitempty"`
}

type Custom_Pull_Request struct {
	github.PullRequest
	Review_Overview []*Review_Overview `json:"review_overview,omitempty"`
}

type Custom_User struct {
	*github.User
	Team *Custom_Team `json:"team,omitempty"`
}

type Custom_Team struct {
	github.Team
	ReviewEnabled *bool `json:"review_enabled,omitempty"`
	ReviewOrder   *int  `json:"review_order,omitempty"`
}

type Set_Team struct {
	Slug          string `json:"slug,omitempty"`
	ReviewEnabled bool   `json:"review_enabled,omitempty"`
	ReviewOrder   int    `json:"review_order,omitempty"`
}
