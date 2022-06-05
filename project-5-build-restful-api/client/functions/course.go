package functions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Course stores the course information
type Course struct { // map this type to the record in the table
	CourseCode  string
	Description string
	CourseName  string
}

// Module stores the module information
type Module struct { // map this type to the record in the table
	ModuleCode  string
	ModuleName  string
	Description string
}

const baseURL = "http://127.0.0.1:5000/api/v1/courses"

// getCourse will retrieve the course data from the database via REST API.
func getCourse(courseCode string, moduleCode string, user string, apiKey string) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in get course:", err)
		}
	}()

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "?user=" + user + "&key=" + apiKey

	if courseCode != "" || moduleCode != "" { // if one/two fields present

		if courseCode == "" && moduleCode != "" { // if only module code
			log.Panicln(errors.New("course code is needed"))
			//fmt.Println("Course code is needed.")
			//return
		}

		if courseCode != "" && moduleCode == "" { // if only course code
			url = baseURL + "/" + courseCode + "?user=" + user + "&key=" + apiKey
		}

		if courseCode != "" && moduleCode != "" { // if both fields present
			url = baseURL + "/" + courseCode + "/" + moduleCode + "?user=" + user + "&key=" + apiKey
		}
	}
	//fmt.Println(url)
	response, err := http.Get(url)
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

// addCourse will add the course data to the database via REST API.
func addCourse(newCourse Course, user string, apiKey string) bool {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in add course:", err)
		}
	}()

	if err := checkValidInput(newCourse.CourseCode, newCourse.CourseName, true); err != nil {
		log.Panicln(err)
	}

	var goCourse map[string]string

	goCourse = map[string]string{"courseCode": newCourse.CourseCode, "courseName": newCourse.CourseName,
		"description": newCourse.Description}

	// Marshall converts go data type to json type
	// in []byte format
	jsonCourse, _ := json.Marshal(goCourse)

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + newCourse.CourseCode + "?user=" + user + "&key=" + apiKey

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonCourse))

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

		// duplicate ID - 409 response code or unprocessable entity - 422 response code
		if response.StatusCode == 409 || response.StatusCode == 422 {
			log.Panicln(errors.New("409 - Duplicate Course Code / 422 - " +
				"Please supply course information with course code and course name in JSON format\n"))
		}
	}

	return true
}

// addModule will add the module data to the database via REST API.
func addModule(newModule Module, courseCode string, user string, apiKey string) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in add module:", err)
			//return false
		}
	}()

	if err := checkValidInput(newModule.ModuleCode, newModule.ModuleName, false); err != nil {
		log.Panicln(err)
	}

	var goModule map[string]string

	goModule = map[string]string{"moduleCode": newModule.ModuleCode, "moduleName": newModule.ModuleName, "description": newModule.Description}

	jsonModule, _ := json.Marshal(goModule)

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + courseCode + "/" + newModule.ModuleCode + "?user=" + user + "&key=" + apiKey

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonModule))
	if err != nil {
		//fmt.Println("error at course")
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in add module with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		// duplicate ID - 409 response code or unprocessable entity - 422 response code
		if response.StatusCode == 409 || response.StatusCode == 422 {
			log.Panicln(errors.New("409 - Duplicate Course Code / 422 - " +
				"Please supply course information with module code and module name in JSON format\n"))
		}

	}
	//return true
}

// updateCourse will update the course data to the database via REST API.
func updateCourse(newCourse Course, user string, apiKey string) bool {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in update course:", err)
			//return false
		}
	}()

	if err := checkValidInput(newCourse.CourseCode, newCourse.CourseName, true); err != nil {
		log.Panicln(err)
	}

	var goCourse map[string]string

	goCourse = map[string]string{"courseCode": newCourse.CourseCode, "courseName": newCourse.CourseName,
		"description": newCourse.Description}

	// Marshall converts go data type to json type
	// in []byte format
	jsonCourse, _ := json.Marshal(goCourse)

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + newCourse.CourseCode + "?user=" + user + "&key=" + apiKey

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonCourse))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		//fmt.Println("error at course")
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in update course with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		// unprocessable entity - 422 response code
		if response.StatusCode == 422 {
			log.Panicln(errors.New("422 - Please supply course information with course code and course name in JSON format\n"))
		}
	}

	return true
}

// updateModule will update the module data to the database via REST API.
func updateModule(newModule Module, courseCode string, user string, apiKey string) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered in update module:", err)
			//return false
		}
	}()

	if err := checkValidInput(newModule.ModuleCode, newModule.ModuleName, false); err != nil {
		log.Panicln(err)
	}

	var goModule map[string]string

	goModule = map[string]string{"moduleCode": newModule.ModuleCode, "moduleName": newModule.ModuleName, "description": newModule.Description}

	jsonModule, _ := json.Marshal(goModule)

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + courseCode + "/" + newModule.ModuleCode + "?user=" + user + "&key=" + apiKey

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonModule))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		//fmt.Println("error at course")
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in update module with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))

		// unprocessable entity - 422 response code
		if response.StatusCode == 422 {
			log.Panicln(errors.New("422 - Please supply course information with module code and module name in JSON format\n"))
		}
	}

}

// deleteCourse will delete the course data from the database via REST API.
func deleteCourse(courseCode string, user string, apiKey string) {

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + courseCode + "?user=" + user + "&key=" + apiKey

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in delete course with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
	}
}

// deleteModule will delete the module data from the database via REST API.
func deleteModule(courseCode string, moduleCode string, user string, apiKey string) {

	// url with query strings to check for
	// authenticate user and valid api key
	url := baseURL + "/" + courseCode + "/" + moduleCode + "?user=" + user + "&key=" + apiKey

	request, err := http.NewRequest(http.MethodDelete, url, nil)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		//fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Panicln(errors.New("The HTTP request failed in delete module with error %s\n"), err)
	} else {
		// close the body at the end
		defer response.Body.Close()

		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
	}
}

// checkValidInput will check if both the course/module code and name are present.
func checkValidInput(code string, name string, isCourse bool) error {

	if isCourse {
		if name == "" || code == "" { // if one/both field empty
			return errors.New("need to have both course name and course code")
		}
	} else {
		if !((name == "" && code == "") || (name != "" && code != "")) { // either one is present only
			return errors.New("need to have both module code and name or nothing at all")
		}
	}

	return nil
}

// ClearScreen is a function to clear screen.
func ClearScreen() {
	fmt.Print("\033c")
}
