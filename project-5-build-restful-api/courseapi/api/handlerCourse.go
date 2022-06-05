package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jpschew/GoSchool-Projects/project-5-build-restful-api/courseapi/config"
	"github.com/jpschew/GoSchool-Projects/project-5-build-restful-api/courseapi/database"
	"io"
	"log"
	"net/http"
)

var uniCourses map[string]database.Course // map course id to course info
//var eeeMod, csMod, accMod map[string]module // map module id to module info
//var modMap map[string]map[string]module     // map course id to module map

var userAPIKey map[string]string

var (
	sqlDriver = "mysql"
	dbPath    = "database"
)
var dsn string
var dbName string

func init() {
	dbConfig, err := config.LoadDBConfig(dbPath, "db")
	if err != nil {
		log.Fatalln(err.Error())
	}
	dsn = dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.DBIP + ":" + dbConfig.DBPort + ")"
	dbName = dbConfig.DBName
	fmt.Println(dsn, dbName)

}

// Home is the handler function for displaying a welcome message to the user.
func Home(res http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(res, "Welcome to the BEST UNIVERSITY COURSES API!")

	// display the message to user
	fmt.Fprintln(res, "Welcome to the BEST UNIVERSITY COURSES API!")

}

// Allcourses is the handler function for retrieving all the course data.
func Allcourses(res http.ResponseWriter, req *http.Request) {

	// check for valid access token
	if !validKey(req) {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("401 -Invalid key"))
		return
	}

	//fmt.Println(uniCourses)
	for _, courseDetails := range uniCourses {
		//fmt.Println("start")
		fmt.Println(courseDetails.CourseName + " is " + courseDetails.Description + " and its course code is " + courseDetails.CourseCode)
		//fmt.Println("end")
	}
	// returns all the courses in JSON converted from GO data
	json.NewEncoder(res).Encode(uniCourses)
}

// Course is the handler function for retrieving a specific course or module, adding/updating course and/or module data as well as deleting a course or module.
func Course(res http.ResponseWriter, req *http.Request) {

	// check for valid access token
	if !validKey(req) {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("401 -Invalid key"))
		return
	}

	//fmt.Println("inside course")

	// Vars returns the route variables for the current request, if any from the gorilla mux
	// return a map with key string and value string
	// params is a map with key(courseid) as string and value(specified in url) as string
	params := mux.Vars(req)

	// Get method - retrieve and request data from a specified recourse (url)
	// does not require the body
	if req.Method == "GET" {

		if _, ok := uniCourses[params["courseid"]]; ok { // check if the key exists
			if params["moduleid"] != "" { // check if the module id is in the url

				if _, ok := uniCourses[params["courseid"]].CourseModules[params["moduleid"]]; ok { // check if key exists

					// returns all the courses in JSON converted from GO data
					json.NewEncoder(res).Encode(uniCourses[params["courseid"]].CourseModules[params["moduleid"]])
					fmt.Println("Retreiving data...", uniCourses[params["courseid"]].CourseModules[params["moduleid"]])

				} else {
					notFoundStatus(res, false)
				}
			} else { // if module id not in url, means only course id
				json.NewEncoder(res).Encode(uniCourses[params["courseid"]])
				fmt.Println("Retreiving data...", uniCourses[params["courseid"]])
			}

		} else {
			notFoundStatus(res, true)
		}
	}

	// Delete method - delete data indicated by a specified resource (url)
	// does not require body
	if req.Method == "DELETE" {

		// connect to database to delete course/module from course/module table
		db := database.CreateDBConn(sqlDriver, dsn, dbName)
		defer db.Close()

		// check if the key is in the uniCourse
		if _, ok := uniCourses[params["courseid"]]; ok {

			if params["moduleid"] != "" { // check if the module id is in the url

				if _, ok := uniCourses[params["courseid"]].CourseModules[params["moduleid"]]; ok { // check if key exists

					// get the module map
					modEdit := uniCourses[params["courseid"]].CourseModules

					// this will delete the key from a map
					// first arg is the map
					// second arg is the key
					delete(modEdit, params["moduleid"])

					database.DeleteModule(db, params["moduleid"], false)
					acceptedStatus(res, params, false, false)

				} else {
					notFoundStatus(res, false)
				}
			} else { // if module id not in url, means only course id

				// this will delete the key from a map
				// first arg is the map
				// second arg is the key
				delete(uniCourses, params["courseid"])

				database.DeleteCourse(db, params["courseid"])
				acceptedStatus(res, params, true, false)
			}

		} else {
			notFoundStatus(res, true)
		}
	}

	// get the header from request
	// check if body is in json format
	// use Content-Type to check for the resource type
	// for POST and PUT, information is sent via the request body
	if req.Header.Get("Content-Type") == "application/json" {
		//fmt.Println("json type")

		// POST is for creating new course
		// if duplicate course or module is created
		// will ask user to use PUT to update the content instead
		if req.Method == "POST" {

			// connect to database to delete course/module from course/module table
			db := database.CreateDBConn(sqlDriver, dsn, dbName)
			defer db.Close()

			//fmt.Println("in post")
			if _, ok := uniCourses[params["courseid"]]; ok { // if course id exists, display duplicate course error
				//fmt.Println("key exists")

				if params["moduleid"] != "" { // if the module id is in the url
					if _, ok := uniCourses[params["courseid"]].CourseModules[params["moduleid"]]; ok { // if module id exists
						conflictStatus(res, false)
					} else { // if module id does not exists

						// close the request body at the end of the function after reading
						defer req.Body.Close()

						// read in data from request body
						if reqBody, err := io.ReadAll(req.Body); err != nil { // error when reading request body
							//fmt.Println("body error")
							unprocessableEntityStatus(res, false)
						} else { // no error when reading request body

							//fmt.Println("no body error")

							newModule := database.Module{"", "", ""}

							json.Unmarshal(reqBody, &newModule)

							// create new module
							createModule(res, reqBody, params, newModule, db)
						}

					}
				} else { // if module not in url, but course id exists
					conflictStatus(res, true)
				}

			} else { // if course id doesnt exists, add to course map

				//fmt.Println("key does not exists")

				// close the request body at the end of the function after reading
				defer req.Body.Close()

				// read in data from request body
				if reqBody, err := io.ReadAll(req.Body); err != nil { // if there is error when reading request body

					//fmt.Println("body error")

					unprocessableEntityStatus(res, true)

				} else { // no error when reading request body
					//fmt.Println("body no error")

					newCourse := database.Course{"", "", "", make(map[string]database.Module)}

					json.Unmarshal(reqBody, &newCourse)

					// create new course
					createCourse(res, reqBody, params, newCourse, db)

				}
			}
		}

		if req.Method == "PUT" {

			// connect to database to delete course/module from course/module table
			db := database.CreateDBConn(sqlDriver, dsn, dbName)
			defer db.Close()

			//fmt.Println("in put")

			// close the request body at the end of the function after reading
			defer req.Body.Close()

			// read in data from request body
			if reqBody, err := io.ReadAll(req.Body); err != nil { // if there is error when reading request body

				//fmt.Println("body error")

				unprocessableEntityStatus(res, true)

			} else { // no error when reading request body
				//fmt.Println("body no error")

				if _, ok := uniCourses[params["courseid"]]; ok { // if course id exists, update content

					if params["moduleid"] != "" { // if the module id is in the url

						// if module id exists
						// unmarshall request body to newModule
						newModule := database.Module{"", "", ""}
						json.Unmarshal(reqBody, &newModule)

						if _, ok := uniCourses[params["courseid"]].CourseModules[params["moduleid"]]; ok { // if module id exists, update module content

							// update module
							updateModule(res, reqBody, params, newModule, db)

						} else { // if module id does not exists, create new

							// create new module
							createModule(res, reqBody, params, newModule, db)

						}
					} else { // if module not in url, but course id exists

						// if module id does not exists
						// unmarshall request body to newCourse
						newCourse := database.Course{"", "", "", make(map[string]database.Module)}
						json.Unmarshal(reqBody, &newCourse)

						// update course content
						updateCourse(res, reqBody, params, newCourse, db)

					}

				} else { // if course id does not exists, create new content

					newCourse := database.Course{"", "", "", make(map[string]database.Module)}
					json.Unmarshal(reqBody, &newCourse)

					// create new course
					createCourse(res, reqBody, params, newCourse, db)

				}
			}
		}
	}
	//fmt.Println("outside")
}

// init initialize the course and key map from the database.
func init() {
	// make a map for uniCourses and userAPIKey
	uniCourses = make(map[string]database.Course)
	userAPIKey = make(map[string]string)

	// initially manually initialized to course map for testing

	//// make maps for different course
	//eeeMod = make(map[string]module)
	//csMod = make(map[string]module)
	//accMod = make(map[string]module)
	//
	//// assign module info (value) to module code(key) in respective course module map
	//eeeMod["EEE2004"] = module{"Digital Electronics", "EEE2004", "Teaching about digital electronics."}
	//eeeMod["EEE2002"] = module{"Analog Electronics", "EEE2002", "Teaching about analog electronics."}
	//
	//csMod["CS1002"] = module{"Data Structures and Algorithms", "CS1002", "Teaching about data structure and algorithms."}
	//csMod["CS1008"] = module{"Golang", "CS1008", "Teaching about Golang programming."}
	//
	//accMod["ACC3010"] = module{"Accounting I", "ACC3010", "Teaching about accounting."}
	//accMod["ACC3014"] = module{"Statistics and Analysis", "ACC3014", "Teaching about statistics."}
	//
	//// assign course info (value) to course code(key) in overall course map
	//uniCourses["EEE"] = courseInfo{"Electrical and Electronic Engineering", "Teaching about EEE.", "EEE", eeeMod}
	//uniCourses["CS"] = courseInfo{"Computer Science", "Teaching about CS.", "CS", csMod}
	//uniCourses["ACC"] = courseInfo{"Accountancy", "Teaching about ACC.", "ACC", accMod}

	// retrieve data from database and assign to the course and key map
	db := database.CreateDBConn(sqlDriver, dsn, dbName)
	defer db.Close()
	uniCourses = database.GetAllRecords(db)
	userAPIKey = database.GetAllKeys(db)

	fmt.Println("uniCourse map:", uniCourses)
	//fmt.Println("User api key:", userAPIKey)

}
