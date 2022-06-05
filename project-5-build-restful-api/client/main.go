package main

import (
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-5-build-restful-api/client/functions"
	"strings"
)

var (
	toContinue = true
)

func main() {

	isStart := true
	var user *string
	var apiKey *string

	for toContinue {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered in main:", err)
			}
		}()

		// start of program will ask user whether to generate an api key for accessing the REST API
		// subsequently will use the username and apikey keyed in by the user
		if isStart {
			user, apiKey = functions.GenerateDeleteKey()
			isStart = false
		}

		// ask user which http method they want to access the REST API
		functions.HttpMethod(user, apiKey)

		var userChoice string

		// after each call to the REST API
		// ask user if they still want to continue to access the REST API
		fmt.Println("Key 'Y' if you still want to continue using the rest API:")
		fmt.Scanln(&userChoice)

		//userChoice = strings.ToUpper(userChoice)

		if strings.ToUpper(userChoice) != "Y" {
			toContinue = false
		}
	}
}
