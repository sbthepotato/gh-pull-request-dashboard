package main

import "github.com/google/go-github/v64/github"

type Config struct {
	Token string `json:"token"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type Review_Overview struct {
	User  string `json:"user"`
	State string `json:"state"`
}

type Custom_Pull_Request struct {
	github.PullRequest
	Review_Overview []Review_Overview `json:"review_overview"`
}

type Custom_User struct {
	github.User
	Team_Name string `json:"team_name"`
	Team_Slug string `json:"team_slug"`
}

type Custom_Team struct {
	github.Team
	Review_Enabled bool `json:"review_enabled"`
	Review_Order   int  `json:"review_order"`
}

type Set_Team struct {
	Slug           string `json:"slug"`
	Review_Enabled bool   `json:"review_enabled"`
	Review_Order   int    `json:"review_order"`
}
