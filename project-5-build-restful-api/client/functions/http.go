package functions

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// HttpMethod asks user on the method they want to access the REST API.
func HttpMethod(user *string, apiKey *string) {

	var httpChoice string

	fmt.Println("Please choose the http method (GET, POST, PUT, DELETE) that you want to execute:")
	//fmt.Scanln(&httpChoice)
	scanner := bufio.NewScanner(os.Stdin)
	// scan for input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	httpChoice = scanner.Text()

	//httpChoice = strings.ToUpper(httpChoice)

	switch strings.ToUpper(httpChoice) {
	case "GET":
		ClearScreen()
		getMethod(user, apiKey)
	case "POST":
		ClearScreen()
		postMethod(user, apiKey)
	case "PUT":
		ClearScreen()
		putMethod(user, apiKey)
	case "DELETE":
		ClearScreen()
		deleteMethod(user, apiKey)
	default:
		fmt.Println("Please key in the correct http method.")
	}
}

// getMethod asks user for the course and/or module code to retrieve the course and/or module data via the REST API.
func getMethod(user *string, apiKey *string) {

	var courseCode string
	var moduleCode string

	fmt.Println("You can retrieve the data from the rest API using this http method.")
	fmt.Println("Leave the Course Code and Module Code blank if you want to retrieve all data.")

	fmt.Println("\nDO NOT leave any space in between as Course Code and Module Code should be one word.")
	fmt.Println("If the words are separated by space, it will automatically merge the words together as one.")

	fmt.Println("\nPlease enter the Course Code (max. 5 char) you want to retrieve:")
	scanner := bufio.NewScanner(os.Stdin)
	// scan for course code input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	courseCode = scanner.Text()

	fmt.Println("\nPlease enter the Module Code (max. 10 char) you want to retrieve:")
	scanner = bufio.NewScanner(os.Stdin)
	// scan for module code input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	moduleCode = scanner.Text()

	//fmt.Println(strings.ToUpper(strings.ReplaceAll(courseCode, " ", "")),
	//	strings.ToUpper(strings.ReplaceAll(moduleCode, " ", "")))

	getCourse(strings.ToUpper(strings.ReplaceAll(courseCode, " ", "")),
		strings.ToUpper(strings.ReplaceAll(moduleCode, " ", "")), *user, *apiKey)

}

// postMethod asks user for the course and/or module code to create the course and/or module data to the database via the REST API.
func postMethod(user *string, apiKey *string) {

	var newCourse Course
	var newModule Module
	var choice string

	fmt.Println("You can create data to the rest API using this http method.")

	newCourse.CourseCode, newCourse.CourseName, newCourse.Description = askForInfo(true)

	// success means course successfully added
	if success := addCourse(newCourse, *user, *apiKey); success { // ask for new module to add to course
		fmt.Println("Key 'Y' to add module to the course:")
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) == "Y" {
			newModule.ModuleCode, newModule.ModuleName, newModule.Description = askForInfo(false)
			addModule(newModule, newCourse.CourseCode, *user, *apiKey)
		}
	}

}

// putMethod asks user for the course and/or module code to create/update the course and/or module data to the database via the REST API.
func putMethod(user *string, apiKey *string) {

	var newCourse Course
	var newModule Module
	var choice string

	fmt.Println("You can add/update data to the rest API using this http method.")

	newCourse.CourseCode, newCourse.CourseName, newCourse.Description = askForInfo(true)

	// success means course success
	if success := updateCourse(newCourse, *user, *apiKey); success { // ask for new module to add to course
		fmt.Println("Key 'Y' to add/update module to the course?")
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) == "Y" {
			newModule.ModuleCode, newModule.ModuleName, newModule.Description = askForInfo(false)
			updateModule(newModule, newCourse.CourseCode, *user, *apiKey)
		}
	}
}

// deleteMethod asks user for the course and/or module code to delete the course and/or module data from the database via the REST API.
func deleteMethod(user *string, apiKey *string) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in delete method:", err)
			//return false
		}
	}()

	var courseCode string
	var moduleCode string

	fmt.Println("You can delete data from the rest API using this http method.")
	fmt.Println("Leave the Course Code and Module Code blank.")

	fmt.Println("\nDO NOT leave any space in between as Course Code and Module Code should be one word.")
	fmt.Println("If the words are separated by space, it will automatically merge the words together as one.")

	fmt.Println("\nPlease enter the Course Code (max. 5 char) you want to delete:")
	scanner := bufio.NewScanner(os.Stdin)
	// scan for course code input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	courseCode = scanner.Text()

	if courseCode == "" {
		log.Panicln(errors.New("need to specify at least course code to delete"))
	}

	fmt.Println("\nPlease enter the Module Code (max. 10 char) you want to delete:")
	scanner = bufio.NewScanner(os.Stdin)
	// scan for module code input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	moduleCode = scanner.Text()

	//fmt.Println(strings.ToUpper(strings.ReplaceAll(courseCode, " ", "")),
	//	strings.ToUpper(strings.ReplaceAll(moduleCode, " ", "")))

	if moduleCode == "" {
		deleteCourse(strings.ReplaceAll(courseCode, " ", ""), *user, *apiKey)
	} else {
		deleteModule(strings.ToUpper(strings.ReplaceAll(courseCode, " ", "")),
			strings.ToUpper(strings.ReplaceAll(moduleCode, " ", "")), *user, *apiKey)
	}
}

// askForInfo asks user for information like course/module code, name and description.
func askForInfo(isCourse bool) (string, string, string) {

	var code string
	var name string
	var description string

	if isCourse {
		fmt.Println("\nPlease enter the Course Code (max. 5 char) you want to add/update:")
		scanner := bufio.NewScanner(os.Stdin)
		// scan for course code input and save it to scanner
		scanner.Scan()
		// get all the data from scanner
		code = scanner.Text()

		fmt.Println("\nPlease enter the Course Name (max. 50 char) you want to add/update:")
		scanner = bufio.NewScanner(os.Stdin)
		// scan for course name input and save it to scanner
		scanner.Scan()
		// get all the data from scanner
		name = scanner.Text()
	} else {
		fmt.Println("\nPlease enter the Module Code (max. 10 char) you want to add/update:")
		scanner := bufio.NewScanner(os.Stdin)
		// scan for module code input and save it to scanner
		scanner.Scan()
		// get all the data from scanner
		code = scanner.Text()

		fmt.Println("\nPlease enter the Module Name (max. 30 char) you want to add/update:")
		scanner = bufio.NewScanner(os.Stdin)
		// scan for module name input and save it to scanner
		scanner.Scan()
		// get all the data from scanner
		name = scanner.Text()

	}

	fmt.Println("\nPlease enter the Description (max. 50 char ) you want to add/update:")
	scanner := bufio.NewScanner(os.Stdin)
	// scan for input and save it to scanner
	scanner.Scan()
	// get all the data from scanner
	description = scanner.Text()

	return code, name, description

}
