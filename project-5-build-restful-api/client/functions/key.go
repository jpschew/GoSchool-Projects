package functions

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const keyURL = "http://127.0.0.1:5000/api/v1"
const returnAPI = "Returning to API key selection menu"

var (
	message = [4]string{
		"1. Generate a new API key.",
		"2. Generate a new API key if you forget the API key that is tagged to you.",
		"3. Delete an API key.",
		"4. Continue to use the rest API if you have your API key",
	}
)

// GenerateDeleteKey will ask user if they want to generate a new API key, update an API key tagged to a user, delete an API key or continue using existing API key.
func GenerateDeleteKey() (*string, *string) {

	var userChoice string
	var userName string
	var apiKey string
	var success bool

	toLoop := true

	for toLoop {
		ClearScreen()
		fmt.Println("You can generate an API key here to use for the course rest API.")
		fmt.Println("Please generate an API key if you do not have one.")
		fmt.Println("The API key will tag to the username you registered.")

		for _, msg := range message {
			fmt.Println(msg)
		}

		fmt.Scanln(&userChoice)
		switch userChoice {
		case "1":
			ClearScreen()
			if userName, apiKey, success = addUpdateKey(false); success {
				toLoop = false
			} else {
				returnToAPI()
			}
		case "2":
			ClearScreen()
			if userName, apiKey, success = addUpdateKey(true); success {
				toLoop = false
			} else {
				returnToAPI()
			}
		case "3":
			ClearScreen()
			if success = deleteKey(); success {
				var choice string
				fmt.Println("Key in 'Y' to continue using the rest API else it will end the application:")
				fmt.Scanln(&choice)
				if strings.ToUpper(choice) != "Y" {
					os.Exit(0)
				}
			} else {
				returnToAPI()
			}
		case "4":
			ClearScreen()
			userName, apiKey = askForUserAPI(false)
			toLoop = false
		default:
			fmt.Println("Wrong Choice.")
		}
	}
	//fmt.Println(userName, apiKey)

	return &userName, &apiKey
}

// deleteKey deletes an API key tagged to a user from the database via REST API.
func deleteKey() bool {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in add key:", err)
		}
	}()

	var user string
	var userAPI map[string]string

	user, _ = askForUserAPI(true)

	userAPI = map[string]string{"user": user, "apiKey": ""}

	// Marshall converts go data type to json type
	// in []byte format
	jsonData, _ := json.Marshal(userAPI)

	// route url for deleting key
	url := keyURL + "/deletekey"

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonData))

	if err != nil {
		//fmt.Println("error at course")
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in add course with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		if response.StatusCode == 404 {
			log.Panicln(errors.New("404 - User not found!.\n"))
		}
	}
	return true

}

// addUpdateKey adds/updates an API key to a user in the database via REST API.
func addUpdateKey(isUpdate bool) (string, string, bool) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in add key:", err)
		}
	}()

	var userAPI map[string]string
	var user string
	var apiKey string
	//isStart = true
	generateKey()

	user, apiKey = askForUserAPI(false)

	userAPI = map[string]string{"user": user, "apiKey": apiKey}
	//fmt.Println(userAPI)

	// Marshall converts go data type to json type
	// in []byte format
	jsonData, _ := json.Marshal(userAPI)

	var response *http.Response
	var request *http.Request
	var err error

	// route url for adding/updating key
	url := keyURL + "/addupdatekey"

	if isUpdate {
		request, err = http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonData))
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err = client.Do(request)
	} else {
		response, err = http.Post(url, "application/json",
			bytes.NewReader(jsonData))
	}

	if err != nil {
		//fmt.Println("error at course")
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in add course with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		if isUpdate {
			if response.StatusCode == 404 {
				log.Panicln(errors.New("404 - User not found!\n"))
			}
		} else {
			// duplicate ID - 409 response code or unprocessable entity - 422 response code
			if response.StatusCode == 409 || response.StatusCode == 422 {
				log.Panicln(errors.New("409 - Duplicate user. You already have one api key/" +
					"422 - Please supply user and key information in JSON format.\n"))
			}
		}

	}
	return user, apiKey, true
}

// generateKey will generate a new API key for the user via REST API.
func generateKey() {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in get course:", err)
		}
	}()

	// route url for generating key
	url := keyURL

	response, err := http.Get(url + "/genkey")
	if err != nil {
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		// ioutil.ReaAll reads from body of request
		// and returns the data it reads s []byte
		data, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(response.StatusCode)
		fmt.Println(string(data))
	}
}

func askForUserAPI(isDelete bool) (string, string) {

	var user string
	var apiKey string

	fmt.Println("Please enter your username:")
	scanner := bufio.NewScanner(os.Stdin)
	// scan for course name input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	user = scanner.Text()
	user = strings.ReplaceAll(user, " ", "")

	if !isDelete { // if not delete method, ask for this info
		fmt.Println("Please key in the generated api key:")
		scanner = bufio.NewScanner(os.Stdin)
		// scan for course name input and save it to scanner
		scanner.Scan()
		// get all the data from scanner
		apiKey = scanner.Text()
	}

	return user, apiKey
}

func returnToAPI() {
	fmt.Println("Please select your choice again.")
	fmt.Print(returnAPI + ".")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
}
