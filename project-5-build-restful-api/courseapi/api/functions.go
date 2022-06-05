package api

import (
	"database/sql"
	"github.com/jpschew/GoSchool-Projects/project-5-build-restful-api/courseapi/database"
	"net/http"
)

// createModule will create a new module and add to the overall course map as well as adding to the database
func createModule(res http.ResponseWriter, body []byte, urlParams map[string]string, newModule database.Module, db *sql.DB) {

	if newModule.ModuleName == "" || newModule.ModuleCode == "" {
		unprocessableEntityStatus(res, false)
	} else {
		// retrieve the course to edit
		courseEdit := uniCourses[urlParams["courseid"]]
		//fmt.Println(courseEdit, courseEdit.CourseModules)
		// set the module map of key (moduleid) to the value passed in
		courseEdit.CourseModules[urlParams["moduleid"]] = newModule
		//fmt.Println(courseEdit)
		database.InsertModule(db, newModule, urlParams["courseid"])
		createdStatus(res, urlParams, false)
	}
}

// createCourse will create a new course and add to the overall course map as well as adding to the database
func createCourse(res http.ResponseWriter, body []byte, urlParams map[string]string, newCourse database.Course, db *sql.DB) {

	if newCourse.CourseName == "" || newCourse.CourseCode == "" { // if one/both field empty

		//fmt.Println("one or two field empty")

		unprocessableEntityStatus(res, true)

	} else { // if both fields present
		//fmt.Println("all fields present")
		uniCourses[urlParams["courseid"]] = newCourse
		database.InsertCourse(db, newCourse)
		createdStatus(res, urlParams, true)
	}
}

// updateModule will update a module in the overall course map as well as updating the database
func updateModule(res http.ResponseWriter, body []byte, urlParams map[string]string, newModule database.Module, db *sql.DB) {

	if newModule.ModuleName == "" || newModule.ModuleCode == "" {
		unprocessableEntityStatus(res, false)
	} else {
		// retrieve the course to edit
		courseEdit := uniCourses[urlParams["courseid"]]
		// set the module map of key (moduleid) to the value passed in
		courseEdit.CourseModules[urlParams["moduleid"]] = newModule
		//acceptedStatusUpdate(res, params, false, newModule.ModuleCode, newModule.ModuleName)
		database.EditModule(db, newModule, urlParams["courseid"])
		acceptedStatus(res, urlParams, false, true)
	}
}

// updateCourse will update a course in the overall course map as well as updating the database
func updateCourse(res http.ResponseWriter, body []byte, urlParams map[string]string, newCourse database.Course, db *sql.DB) {

	if newCourse.CourseName == "" || newCourse.CourseCode == "" { // if one/both field empty

		//fmt.Println("one or two field empty")

		unprocessableEntityStatus(res, true)

	} else { // if both fields present
		if len(uniCourses[urlParams["courseid"]].CourseModules) != 0 { // there is existing module in course
			currMod := uniCourses[urlParams["courseid"]].CourseModules
			newCourse.CourseModules = currMod
		}

		uniCourses[urlParams["courseid"]] = newCourse

		database.EditCourse(db, newCourse)

		if len(newCourse.CourseModules) == 0 { // if no module in course, check for empty map

			// delete all module associated with the course
			database.DeleteModule(db, newCourse.CourseCode, true)
		}

		acceptedStatus(res, urlParams, true, true)

	}
}

// validKey checks for a valid key to secure the REST API
// so that only authenticated user can use the REST API
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	//fmt.Println(v)
	//fmt.Println(userAPIKey)
	// check if user exists
	if user, ok := v["user"]; ok && user[0] != "" {
		//fmt.Println(user[0])
		// check if key exists
		if key, ok := v["key"]; ok && key[0] != "" {
			//fmt.Println(user[0], key[0])
			// check if key tagger to user is correct
			if userAPIKey[user[0]] == key[0] {
				return true
			} else {
				return false
			}
		}
		return false
	}
	return false
}
