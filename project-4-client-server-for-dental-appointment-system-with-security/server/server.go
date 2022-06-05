// Package server implements the handler functions for the respective routes.
// Each handler function will execute the respective features according to the route.
package server

// important for assignment
import (
	"errors"
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/config"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/datatype"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/userpackage"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/utils"
	"html/template"
	"net/http"
	"strconv"

	// go install for 1.18
	// go get for version before 1.18

	uuid "github.com/satori/go.uuid"

	// add in encryption
	"golang.org/x/crypto/bcrypt"
)

// user struct will store the information of the user
type user struct {
	Username string
	Password []byte
	First    string
	Last     string
	Admin    bool
}

// userData struct will store the information keyed in by user to list available dentist time as well as data to be shown at the browser
type userData struct {
	Month       string
	Day         int
	Year        int
	DentistList []string
	DentistName string
	TimeSlot    []string
}

// appt struct will store the information keyed in by user to make an appointment as well as data to be shown at the browser
type appt struct {
	TimeSlot    []string
	Month       string
	Day         int
	Year        int
	TimeSession string
	DentistList []string
	Message     string
}

// apptEdit struct will store the information keyed in by user to edit an appointment as well as data to be shown at the browser
type apptEdit struct {
	// apptId int
	TimeSlot    []string
	ApptId      string
	Message     string
	Month       string
	Day         int
	Year        int
	TimeSession string
	DentistList []string
	DentistName string
	ApptMessage string
}

// variables declared here will be used for the various function in this package
var (
	tpl *template.Template

	// key is username
	mapUsers = map[string]user{}

	// key is cookie value, key is username
	mapSessions = map[string]string{}

	listTime userData

	apptInfo appt

	editApptInfo apptEdit

	userMonth, userDay, apptMonthStr, timeSession, userId   string
	userApptDay, userApptMonth, userTimeSession, userApptId int
	dentistResult, timeSlot                                 []string

	myUser user

	details *datatype.ApptDetails

	outputMessage string

	filename = [4]string{
		"user1",
		"user2",
		"user3",
		"admin1",
	}
)

// constants declared here will be used to specify the path for config file as well as the session timeout time
const (
	timeout = 1800
	path    = "config"
)

// init initialized 3 users and 1 admin for this application using the env file.
func init() {

	// for testing purpose
	userInfo := user{}
	for i, file := range filename {
		// fmt.Println(file)
		config, err := config.LoadConfig(path, file)
		if err != nil {
			// log.Fatalln(err)
			utils.ErrorLogging("Configuration error occured:", err)
		}
		// fmt.Println(config.UserName, config.Password, config.FirstName, config.LastName)
		bPassword, _ := bcrypt.GenerateFromPassword([]byte(config.Password), bcrypt.MinCost)
		if i == len(filename)-1 {
			userInfo = user{config.UserName, bPassword, config.FirstName, config.LastName, true}
		} else {
			userInfo = user{config.UserName, bPassword, config.FirstName, config.LastName, false}
		}

		mapUsers[userInfo.Username] = userInfo
	}
	// fmt.Println(mapUsers)

	// load all templates from template folder
	tpl = template.Must(template.ParseGlob("templates/*"))

}

// // StartServer will start the server and it include all the handler functions for the respective routes.
// func StartServer() {

// 	// gorilla mux is used here to limit the request method for different routes
// 	r := mux.NewRouter()
// 	r.HandleFunc("/", index).Methods(http.MethodGet)
// 	r.HandleFunc("/signup", signup)
// 	r.HandleFunc("/login", login)
// 	r.HandleFunc("/browse", browse)
// 	r.HandleFunc("/search", searchAppt)
// 	r.HandleFunc("/make", makeAppt)
// 	r.HandleFunc("/list", list)
// 	r.HandleFunc("/edit", editAppt)
// 	r.HandleFunc("/logout", logout).Methods(http.MethodGet)
// 	r.Handle("/favicon.ico", http.NotFoundHandler())
// 	http.ListenAndServe(":8080", r)

// 	// https is used for secure communications between client and server
// 	// http.ListenAndServeTLS(":8081",
// 	// 	"key/cert.pem",
// 	// 	"key/key.pem",
// 	// 	r)
// }

// Index is a handler function that will ask the user to sign up for a new account or log in to their existing account.
// It will read in the request and response and display the index page using a template.
func Index(res http.ResponseWriter, req *http.Request) {
	cleanData()
	// fmt.Println(userMonth, userDate, apptMonthStr, userApptDate, userApptMonth, dentistResult)

	// check if user exists
	myUser = getUser(res, req)
	// utils.TraceLogging(myUser.Username + "has logged in.")
	tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

// Signup is a handler function that will ask the user to sign up for a new account.
// It will ask the user for his/her information as well as creating a new session using the cookie value, which is tagged to the username of the user.
// It will read in the request and response and display the signup page using a template.
func Signup(res http.ResponseWriter, req *http.Request) {

	cleanData()

	// check if user already logged in
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// fmt.Println(req.Method)

	var myUser user

	// process form submission
	// in the template signup, there is a post method form
	// if user submit the information will get into this if loop
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstName := req.FormValue("firstname")
		lastName := req.FormValue("lastname")
		// admin := false
		fmt.Println("logged in")
		if username != "" {
			// check if username exist/ taken
			if _, ok := mapUsers[username]; ok {
				utils.TraceLogging("Username already taken when signing up.")
				utils.WarnLogging("Username already taken when signing up.", nil)
				http.Error(res, "Username already taken", http.StatusForbidden)
				// need to return as there is further processing
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:   "myCookie",
				Value:  id.String(),
				MaxAge: timeout, // set a time of 30min for the session
			}
			http.SetCookie(res, myCookie)

			// set the cookie value to the username
			// so we know which username using which cookie value
			mapSessions[myCookie.Value] = username

			// 1st argument pass in password in []byte format]
			// 2nd arguemnt determine how difficult the decrytion will be
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				utils.TraceLogging("Password is not hashed when signing up.")
				utils.WarnLogging("Password is not hashed when signing up.", err)
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				// need to return as there is further processing
				return
			}

			myUser = user{username, bPassword, firstName, lastName, false}
			mapUsers[username] = myUser
			utils.TraceLogging("New user " + myUser.Username + " has signed up.")
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

// Log in is a handler function that will ask the user to log in to their existing account.
// It will ask the user for his/her username and passowrd to log in their account as well as creating a new session using the cookie value, which is tagged to the username of the user.
// It will read in the request and response and display the login page using a template.
func Login(res http.ResponseWriter, req *http.Request) {

	cleanData()

	// check if user already logged in
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		myUser, ok := mapUsers[username]
		if !ok {
			utils.TraceLogging("Username and/or password do not match when logging in.")
			utils.WarnLogging("Username and/or password do not match when logging in.", nil)
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			utils.TraceLogging("Username and/or password do not match when logging in.")
			utils.WarnLogging("Username and/or password do not match when logging in.", err)
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:   "myCookie",
			Value:  id.String(),
			MaxAge: timeout, // set a time of 30min for the session
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username
		utils.TraceLogging(myUser.Username + " has logged in.")
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

// Logout is a handler function that will log the user out from their account.
// It will also delete the session using the cookie value, which is tagged to the username of the user.
// It will read in the request and response and redirect the user to the index page after logging out.
func Logout(res http.ResponseWriter, req *http.Request) {

	cleanData()

	// check if user already logged in
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	// built in function to delete key value pair in a map
	// first arg is the map
	// second arg is the key to delete in the map
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	// set MaxAge to < 0 to delete cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	utils.TraceLogging(myUser.Username + " has logged out.")
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

// getUser will set cookie using github.com/satori/go.uuid package to set unique identifer to cookie value as well as a 30min timeout for the session.
// It will also check if the same session exists.
// If it exists, it will return the user information.
func getUser(res http.ResponseWriter, req *http.Request) user {

	// get current session cookie
	// cookie name can be anything
	// here we used myCookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:   "myCookie",
			Value:  id.String(),
			MaxAge: timeout, // set a time of 30min for the session
		}

	}
	http.SetCookie(res, myCookie)

	// if the user exists already, get user
	var myUser user
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
	}

	return myUser
}

// alreadyLoggedIn checks if the user is already log in.
// It will return true if the user is logged in and false if the user has not logged in yet.
func alreadyLoggedIn(req *http.Request) bool {

	// retrieve cookie and check if there is error
	myCookie, err := req.Cookie("myCookie")
	// if there is error means cookie not available
	if err != nil {
		return false
	}
	// check if there is any username tied to the cookie value
	username := mapSessions[myCookie.Value]
	// check if the username is valid
	_, ok := mapUsers[username]
	return ok
}

// Browse is a handler function that will ask the user for the dentist name to browse all the appointments by the dentist.
// It will read in the request and response and display the browse page using a template.
func Browse(res http.ResponseWriter, req *http.Request) {

	// check if user already logged in
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in browse appointment:", err)
			http.Error(res, "Invalid dentist name has been entered", http.StatusBadRequest)
			return
		}
	}()

	// var dentistName string
	var Output []string
	var err error

	if req.Method == http.MethodPost {

		dentistName := req.FormValue("dentist")

		if dentistName != "" {
			Output, err = userpackage.BrowseDrAppointment(dentistName)
			if err != nil {
				utils.TraceLogging(myUser.Username + " encountered some error when browsing appointment.")
				utils.WarnLogging(myUser.Username+" encountered some error when browsing appointment.", err)
			} else {
				utils.TraceLogging(myUser.Username + " is browsing dentist appointment.")
			}
		}
	}

	// if Error != nil {
	// 	http.Error(res, "Wrong dentist name", http.StatusBadRequest)
	// 	return
	// }
	tpl.ExecuteTemplate(res, "browse.gohtml", Output)
}

// SearchAppt is a handler function that will ask the user for the Appointment Id to search for a particular appointment
// It will read in the request and response and display the search page using a template.
func SearchAppt(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search appointment:", err)
			http.Error(res, "Invalid appointment id has been entered", http.StatusBadRequest)
			return
		}
	}()

	// var dentistName string
	var Output string
	// var err error
	// fmt.Println("in search")
	if req.Method == http.MethodPost {

		idStr := req.FormValue("apptId")
		id, err := strconv.Atoi(idStr)

		// fmt.Printf("%v, %T", id, id)
		if err != nil {
			utils.TraceLogging(myUser.Username + " has keyed in an invalid appointment id.")
			utils.WarnLogging(myUser.Username+" has keyed in an invalid appointment id.", err)
			http.Error(res, "Invalid appointment id has been entered", http.StatusBadRequest)
			return
		}

		Output, err = userpackage.SearchAppointment(id)
		if err != err {
			utils.TraceLogging(myUser.Username + " encountered some error when searching appointment.")
			utils.WarnLogging(myUser.Username+" encountered some error when searching appointment.", err)
		} else {
			utils.TraceLogging(myUser.Username + " is searching for an appointment.")
		}

	}

	// if Error != nil {
	// 	http.Error(res, "Wrong dentist name", http.StatusBadRequest)
	// 	return
	// }
	tpl.ExecuteTemplate(res, "search.gohtml", Output)
}

// MakeAppt is a handler function that will ask the user for information to make their appointment.
// It will read in the request and response and display the make page using a template.
func MakeAppt(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// var timeSession string
	// var month string
	var err error
	// var userTimeSession int

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in make appointment:", err)
			http.Error(res, "Error occurred when making appointment", http.StatusBadRequest)
			return
		}
	}()

	apptInfo = appt{datatype.TimeArr[:], "", 0, datatype.ApptYear, "", []string{}, ""}

	if req.Method == http.MethodPost {
		// timeSession = req.FormValue("timeSession")
		if userMonth == "" && userDay == "" {
			userDay = req.FormValue("day")
			userMonth = req.FormValue("month")

			userApptMonth, userApptDay, err = checkValidDate(userMonth, userDay)
			if err != nil {
				userMonth = ""
				userDay = ""
				utils.TraceLogging(myUser.Username + " has keyed in invalid date while making appointment.")
				utils.WarnLogging(myUser.Username+" has keyed in invalid date while making appointment.", err)
				http.Error(res, "Invalid month/day has been entered. Month need to be integer between 1 and 12. Day need to be integer between 1 and 31.", http.StatusBadRequest)
				return
			}
			apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

			timeSession = req.FormValue("time")
			// userTimeSession, err = strconv.Atoi(timeSession)
			userTimeSession, err = checkValidTime(timeSession)
			timeSession = datatype.TimeArr[userTimeSession-1]
			if err != nil {
				userMonth = ""
				userDay = ""
				timeSession = ""
				utils.TraceLogging(myUser.Username + " has keyed in invalid time session while making appointment.")
				utils.WarnLogging(myUser.Username+" has keyed in invalid time session while making appointment.", err)
				http.Error(res, "Invalid times ession has been entered. Time session need to be integer between 1 and 15", http.StatusBadRequest)
				return
			}
			// fmt.Println(day, month)

			dentistResult = userpackage.MakeAppointment(userApptMonth, userApptDay, userTimeSession)
			// fmt.Println(dentistResult)
			apptInfo = appt{datatype.TimeArr[:], apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, ""}
		} else {
			dentistName := req.FormValue("dentist")
			dentistName = utils.ConvertToUpper(dentistName)
			// apptInfo = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
			// fmt.Println(dentistName)
			if dentistName != "" {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println("Panic occurred as wrong dentist name is keyed:", err)
						http.Error(res, "Error occurred when editing appointment. Wrong Dentist name entered!", http.StatusBadRequest)
						return
					}
				}()
				if found, _, err := datatype.DentistHash.Search(dentistName, userTimeSession, userApptMonth, userApptDay); !found {
					utils.TraceLogging(myUser.Username + " has keyed in invalid dentist while making appointment.")
					utils.WarnLogging(myUser.Username+" has keyed in invalid dentist while making appointment.", err)
					panic(err)
				} else {
					result := userpackage.AddToApptList(userApptMonth, userApptDay, userTimeSession, dentistName, myUser.Username)
					// if err != nil {
					// 	utils.TraceLogging(myUser.Username + " has keyed in invalid dentist for the date and time for the appointment.")
					// 	utils.WarnLogging(myUser.Username+" has keyed in invalid dentist for the date and time for the appointment.", err)
					// 	panic(err)
					// }

					apptInfo = appt{datatype.TimeArr[:], apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, result}
					// fmt.Println(apptInfo)
					utils.TraceLogging(myUser.Username + " has made an appointment.")
				}
			}
		}

	}

	tpl.ExecuteTemplate(res, "make.gohtml", apptInfo)
}

// checkValidDate checks if the month and day is a in the correct range for the date
// It will take in the month and day as string and return them as integer.
func checkValidDate(m string, d string) (int, int, error) {

	var mInt, dInt int
	var err error
	// convert from string to int
	// if input is not integer return error
	mInt, err = strconv.Atoi(m)
	if err != nil {
		// http.Error(res, "Invalid month has been entered. Month need to be integer between 1 and 12", http.StatusBadRequest)
		return 0, 0, err
	}

	dInt, err = strconv.Atoi(d)
	if err != nil {
		// http.Error(res, "Invalid date has been entered. Date need to be integer between 1 and 31", http.StatusBadRequest)
		return 0, 0, err
	}

	// checkForValidDate is a function declared in patient.go
	// if both are integer, check if range is correct
	// fmt.Println(mInt, dInt, checkForValidDate(mInt, dInt))
	if !utils.CheckForValidDate(mInt, dInt) {
		return 0, 0, errors.New("date selected is out of range")
	}

	return mInt, dInt, nil

}

// checkValidTime checks if time input is valid.
// It will take in time as string and return back as integer.
func checkValidTime(t string) (int, error) {

	var tInt int
	var err error
	// convert from string to int
	// if input is not integer return error
	tInt, err = strconv.Atoi(t)
	if err != nil {
		// http.Error(res, "Invalid month has been entered. Month need to be integer between 1 and 12", http.StatusBadRequest)
		return 0, err
	}

	if !checkForValidTime(tInt) {
		return 0, errors.New("time slot selected is out of range")
	}

	return tInt, nil

}

// checkValidForTime checks if time slot input is in the correct range.
// It will take in time slot as int and return true if it is in a valid range, else false.
func checkForValidTime(timeSession int) bool {

	if timeSession < 1 || timeSession > 15 {
		// panic(errors.New("date is out of range"))
		return false
	}
	return true
}

// List is a handler function that will list the available time slot for a particular dentist.
// It will read in the request and response and display the list page using a template.
func List(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in list available dentist:", err)
			http.Error(res, "Date selected is out of range or Wrong Dentist selected", http.StatusBadRequest)
			return
		}
	}()

	var err error

	// fmt.Println(userMonth, "1", userDay)
	if req.Method == http.MethodPost {
		// fmt.Println(userMonth, "2", userDay)
		if userMonth == "" && userDay == "" {

			// var err error

			userMonth = req.FormValue("month")
			userDay = req.FormValue("day")

			// convert from string to int
			// if input is not integer return error
			userApptMonth, userApptDay, err = checkValidDate(userMonth, userDay)
			if err != nil {
				userMonth = ""
				userDay = ""
				utils.TraceLogging(myUser.Username + " has keyed in invalid date while listing dentist available time.")
				utils.WarnLogging(myUser.Username+" has keyed in invalid date while listing dentist available time.", err)
				http.Error(res, "Invalid month has been entered. Month need to be integer between 1 and 12", http.StatusBadRequest)
				return
			}

			apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

			dentistResult, err = userpackage.ListAvailableDentist(userApptMonth, userApptDay)
			if err != nil {
				panic(err)
			}
			listTime = userData{apptMonthStr, userApptDay, datatype.ApptYear, dentistResult, "", []string{}}
		} else {

			dentistName := req.FormValue("dentist")
			dentistName = utils.ConvertToUpper(dentistName)
			// listTime = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
			// fmt.Println(dentistName)
			if dentistName != "" {
				// fmt.Println("enter")
				// fmt.Println(userApptMonth, userApptDate)
				timeSlot, err = datatype.DentistHash.ListAvailTime(dentistName, userApptMonth, userApptDay)
				if err != nil {
					utils.TraceLogging(myUser.Username + " has keyed in invalid dentist while listing dentist available time.")
					utils.WarnLogging(myUser.Username+" has keyed in invalid dentist while listing dentist available time.", err)
					panic(err)
				}

				listTime = userData{apptMonthStr, userApptDay, datatype.ApptYear, dentistResult, dentistName, timeSlot}
				// fmt.Println(listTime)
				utils.TraceLogging(myUser.Username + " is listing the available dentist time slot.")
			}
		}
	}

	tpl.ExecuteTemplate(res, "list.gohtml", listTime)

}

// EditAppt is a handler function that will ask the user for their Appointment Id to edit their appointment.
// It will read in the request and response and display the edit page using a template.
func EditAppt(res http.ResponseWriter, req *http.Request) {

	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// var timeSession string
	// var month string
	var err error
	// var userTimeSession int

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in make appointment:", err)
			http.Error(res, "Error occurred when editing appointment. Wrong Appointment Id entered!", http.StatusBadRequest)
			return
		}
	}()

	// EditApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, []string{}, "", ""}

	if req.Method == http.MethodPost {
		// fmt.Println(userId, "1")
		if userId == "" {

			defer func() {
				if err := recover(); err != nil {
					userId = ""
					fmt.Println("Panic occurred as wrong appointment id is keyed:", err)
					http.Error(res, "Error occurred when editing appointment. Wrong Appointment Id entered!", http.StatusBadRequest)
					return
				}
			}()

			userId = req.FormValue("apptId")

			userApptId, err = strconv.Atoi(userId)
			if err != nil {
				userId = ""
				utils.TraceLogging(myUser.Username + " has keyed in invalid appointment id while editing appointment.")
				utils.WarnLogging(myUser.Username+" has keyed in invalid appointment id while editing appointment.", err)
				http.Error(res, "Invalid appointment id has been entered. Appointment Id must be integer.", http.StatusBadRequest)
				return
			}

			// var details *apptDetails
			details, outputMessage = userpackage.EditAppointment(userApptId)

			// if user key in the appointment not tagged to its username and is not admin
			// panic an error
			// else if admin, continue to edit the appointment
			if details.Patient() != myUser.Username {
				if myUser.Username != "admin" {
					utils.TraceLogging(myUser.Username + " has keyed in an appointment id that is not tagged to him/her while editing appointment.")
					utils.WarnLogging(myUser.Username+" has keyed in invalid appointment id that is not tagged to him/her while editing appointment.", err)
					panic(errors.New("wrong appointment id entered"))
				}
			}
			// fmt.Println(details, outputMessage)
			editApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, []string{}, "", ""}
		} else {

			if userMonth == "" && userDay == "" {
				userDay = req.FormValue("day")
				userMonth = req.FormValue("month")

				userApptMonth, userApptDay, err = checkValidDate(userMonth, userDay)
				if err != nil {
					userMonth = ""
					userDay = ""
					utils.TraceLogging(myUser.Username + " has keyed in invalid date while editing appointment.")
					utils.WarnLogging(myUser.Username+" has keyed in invalid date while editing appointment.", err)
					http.Error(res, "Invalid month/day has been entered. Month need to be integer between 1 and 12. Day need to be integer between 1 and 31.", http.StatusBadRequest)
					return
				}
				apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

				timeSession = req.FormValue("time")
				userTimeSession, err = checkValidTime(timeSession)
				timeSession = datatype.TimeArr[userTimeSession-1]
				if err != nil {
					userMonth = ""
					userDay = ""
					timeSession = ""
					utils.TraceLogging(myUser.Username + " has keyed in invalid time session while editing appointment.")
					utils.WarnLogging(myUser.Username+" has keyed in invalid times session while editing appointment.", err)
					http.Error(res, "Invalid times ession has been entered. Time session need to be integer between 1 and 15", http.StatusBadRequest)
					return
				}

				dentistResult = userpackage.SearchAvailDentist(userApptMonth, userApptDay, userTimeSession, true)
				editApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, "", ""}
			} else {
				dentistName := req.FormValue("dentist")
				dentistName = utils.ConvertToUpper(dentistName)

				// EditApptInfo = apptEdit{TimeArr[:], outputMessage, apptMonthStr, userApptDay, ApptYear, timeSession, dentistResult, dentistName, ""}
				// fmt.Println(dentistName)
				if dentistName != "" {

					defer func() {
						if err := recover(); err != nil {
							fmt.Println("Panic occurred as wrong dentist name is keyed:", err)
							http.Error(res, "Error occurred when editing appointment. Wrong Dentist name entered!", http.StatusBadRequest)
							return
						}
					}()

					if found, _, err := datatype.DentistHash.Search(dentistName, userTimeSession, userApptMonth, userApptDay); !found {
						utils.TraceLogging(myUser.Username + " has keyed in invalid dentist while editing appointment.")
						utils.WarnLogging(myUser.Username+" has keyed in invalid dentist while editing appointment.", err)
						panic(err)
					} else {

						// fmt.Println(userApptMonth, userApptDate)
						if _, err := details.CheckChanges(userTimeSession, userApptMonth, userApptDay, dentistName); err != nil {
							// fmt.Println("panic")
							utils.TraceLogging(myUser.Username + " has not changed anything while editing appointment.")
							utils.WarnLogging(myUser.Username+" has not changed anything while editing appointment.", err)
							panic(err)
						} else {
							var result string
							if myUser.Username == "admin" { // if admin edit appointment, use back patient username
								result = userpackage.UpdateAppt(details, userTimeSession, userApptMonth, userApptDay, dentistName, userApptId, details.Patient())
							} else { // else use current user username
								result = userpackage.UpdateAppt(details, userTimeSession, userApptMonth, userApptDay, dentistName, userApptId, myUser.Username)
							}
							editApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, dentistName, result}
							utils.TraceLogging(myUser.Username + " has edited an appointment.")
						}
					}
				}
			}
		}

	}

	tpl.ExecuteTemplate(res, "edit.gohtml", editApptInfo)
}

func cleanData() {
	// everytime goes back to other route
	// clean data for listing
	userMonth, userDay, apptMonthStr, userId = "", "", "", ""
	userApptDay, userApptMonth = 0, 0
	dentistResult, timeSlot = []string{}, []string{}

	listTime = userData{}
	editApptInfo = apptEdit{}
}
