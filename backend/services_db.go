package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Write users to file in db/users.json.
func write_users(users map[string]*CustomUser) error {
	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Println("error marshalling users to JSON: ", err)
		return err
	}

	file, err := os.Create("db/users.json")
	if err != nil {
		log.Println("Error creating users file: ", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("Error writing user JSON to file: ", err)
		return err
	}

	return nil
}

// Read users from db/users.json.
// Returns a map where the key is the user.login and the value is the custom user struct
func read_users() (map[string]*CustomUser, error) {
	file, err := os.Open("db/users.json")
	if err != nil {
		log.Println("error reading user file: ", err.Error())
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading user file: ", err.Error())
		return nil, err
	}

	UserMap := make(map[string]*CustomUser)

	err = json.Unmarshal(jsonData, &UserMap)
	if err != nil {
		log.Println("Error unmarshalling users from Json ", err.Error())
		return nil, err
	}

	return UserMap, nil
}

// Write teams to file in db/teams.json
func write_teams(teams map[string]*CustomTeam, active_only bool) error {
	jsonData, err := json.Marshal(teams)
	if err != nil {
		log.Println("Error marshalling teams to JSON: ", err.Error())
		return err
	}

	var file *os.File

	if active_only {
		file, err = os.Create("db/teams_active.json")
	} else {
		file, err = os.Create("db/teams.json")
	}

	if err != nil {
		log.Println("Error creating teams file: ", err.Error())
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("error writing team JSON to file: ", err.Error())
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
	if err != nil {
		log.Println("error reading team file", err.Error())
		return nil, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading team file: ", err.Error())
		return nil, err
	}

	teamMap := make(map[string]*CustomTeam)

	err = json.Unmarshal(jsonData, &teamMap)
	if err != nil {
		log.Println("Error unmarshalling teams from Json ", err.Error())
		return nil, err
	}

	return teamMap, nil

}
