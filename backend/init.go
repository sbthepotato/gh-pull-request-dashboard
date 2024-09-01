package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	Token string
	Owner string
	Repo  string
}

func load_config() *Config {
	content, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return &payload
}

func init_github_connection(ctx context.Context) (*Config, *github.Client) {
	config := load_config()

	authToken := config.Token

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	config.Token = ""

	return config, client

}
