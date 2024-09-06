package main

import "github.com/google/go-github/v64/github"

type Config struct {
	Token string
	Owner string
	Repo  string
}

type Review_Overview struct {
	User  string
	State string
}

type Custom_Pull_Request struct {
	github.PullRequest
	Review_Overview []Review_Overview
}

type Custom_User struct {
	github.User
	Team_Name string
	Team_Slug string
}
