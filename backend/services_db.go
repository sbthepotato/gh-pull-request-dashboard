package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

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

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalln("Error writing user JSON to file: ", err)
	}
}

func read_users() map[string]*Custom_User {
	file, err := os.Open("db/users.json")
	if err != nil {
		log.Println("error reading user file: ", err)
		return nil
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("Error reading file: ", err)
	}

	UserMap := make(map[string]*Custom_User)

	err = json.Unmarshal(jsonData, &UserMap)
	if err != nil {
		log.Fatalln("Error unmarshalling users from Json ", err)
	}

	return UserMap
}
