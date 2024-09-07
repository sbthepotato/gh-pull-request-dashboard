package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Write users to file in db/users.json.
func write_users(users map[string]*Custom_User) {
	jsonData, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("error marshalling users to JSON: ", err)
	}

	file, err := os.Create("db/users.json")
	if err != nil {
		log.Fatalln("Error creating users file: ", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalln("Error writing user JSON to file: ", err)
	}
}

// Read users from db/users.json.
// Returns a map where the key is the user.login and the value is the custom user struct
func read_users() map[string]*Custom_User {
	file, err := os.Open("db/users.json")
	if err != nil {
		log.Println("error reading user file: ", err)
		return nil
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading user file: ", err)
	}

	UserMap := make(map[string]*Custom_User)

	err = json.Unmarshal(jsonData, &UserMap)
	if err != nil {
		log.Fatalln("Error unmarshalling users from Json ", err)
	}

	return UserMap
}

// Write teams to file in db/teams.json
func write_teams(teams map[string]*Custom_Team) {
	jsonData, err := json.Marshal(teams)
	if err != nil {
		log.Fatalln("Error marshalling teams to JSON: ", err.Error())
	}

	file, err := os.Create("db/teams.json")
	if err != nil {
		log.Fatalln("Error creating teams file: ", err.Error())
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalln("error writing team JSON to file: ", err)
	}
}

// Read users from db/teams.json.
// Returns a map where the team.slug is the key and the value is the custom team struct
func read_teams() map[string]*Custom_Team {
	file, err := os.Open("db/teams.json")
	if err != nil {
		log.Println("error reading team file", err)
		return nil
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading team file: ", err)
	}

	teamMap := make(map[string]*Custom_Team)

	err = json.Unmarshal(jsonData, &teamMap)
	if err != nil {
		log.Fatalln("Error unmarshalling teams from Json ", err)
	}

	return teamMap

}
