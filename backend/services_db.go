package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Write users to file in db/users.json.
func write_users(users map[string]*CustomUser) error {
	jsonData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	file, err := os.Create("db/users.json")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Read users from db/users.json.
// Returns a map where the key is the user.login and the value is the custom user struct
func read_users() (map[string]*CustomUser, error) {
	file, err := os.Open("db/users.json")

	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	UserMap := make(map[string]*CustomUser)

	err = json.Unmarshal(jsonData, &UserMap)
	if err != nil {
		return nil, err
	}

	return UserMap, nil
}

// Write teams to file in db/teams.json
func write_teams(teams map[string]*CustomTeam, active_only bool) error {
	jsonData, err := json.Marshal(teams)
	if err != nil {
		return err
	}

	var file *os.File

	if active_only {
		file, err = os.Create("db/teams_active.json")
	} else {
		file, err = os.Create("db/teams.json")
	}

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Read users from db/teams.json.
// Returns a map where the team.slug is the key and the value is the custom team struct
func read_teams(active_only bool) (map[string]*CustomTeam, error) {
	var file *os.File
	var err error

	if active_only {
		file, err = os.Open("db/teams_active.json")
	} else {
		file, err = os.Open("db/teams.json")
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	teamMap := make(map[string]*CustomTeam)

	err = json.Unmarshal(jsonData, &teamMap)
	if err != nil {
		return nil, err
	}

	return teamMap, nil

}

/*
write repo list to db/repos.json
*/
func write_repos(repos []*CustomRepo) error {
	jsonData, err := json.Marshal(repos)
	if err != nil {
		return err
	}

	file, err := os.Create("db/repos.json")

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

/*
read repo list from db/repos.json
parameter decides if list should only return active repos
*/
func read_repos(active_only bool) ([]*CustomRepo, error) {

	file, err := os.Open("db/repos.json")

	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	repos := make([]*CustomRepo, 0)

	err = json.Unmarshal(jsonData, &repos)
	if err != nil {
		return nil, err
	}

	if active_only {
		active_repos := make([]*CustomRepo, 0)

		for _, repo := range repos {
			if *repo.Enabled {
				active_repos = append(active_repos, repo)
			}
		}

		return active_repos, nil
	} else {
		return repos, nil
	}

}
