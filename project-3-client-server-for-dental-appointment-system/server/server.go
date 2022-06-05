package server

// important for assignment
import (
	"errors"
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/datatype"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/userpackage"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/utils"
	"html/template"
	"net/http"
	"strconv"

	// go install for 1.18
	// go get for version before 1.18
	uuid "github.com/satori/go.uuid"

	// add in encryption
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password []byte
	First    string
	Last     string
	Admin    bool
}

type userData struct {
	Month       string
	Day         int
	Year        int
	DentistList []string
	DentistName string
	TimeSlot    []string
}

type appt struct {
	TimeSlot    []string
	Month       string
	Day         int
	Year        int
	TimeSession string
	DentistList []string
	Message     string
}

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

var (
	tpl *template.Template

	// key is username
	mapUsers = map[string]user{}

	// key is cookie value, key is username
	mapSessions = map[string]string{}

	// key is username, value is list of appointment id
	// mapApptId = map[string][]int{}

	Appt userData

	ApptInfo appt

	EditApptInfo apptEdit

	// var timeSlot []string
	userMonth, userDay, apptMonthStr, timeSession, userId   string
	userApptDay, userApptMonth, userTimeSession, userApptId int
	dentistResult, timeSlot                                 []string
	// dentistName                      string

	// Username = []string{"jpschew", "ken", "ivan"}

	myUser user

	details *datatype.ApptDetails

	outputMessage string
)

func init() {

	// for testing purpose
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("jas123per"), bcrypt.MinCost)
	user1 := user{"jpschew", bPassword, "jasper", "chew", false}
	mapUsers[user1.Username] = user1
	bPassword, _ = bcrypt.GenerateFromPassword([]byte("ken123tan"), bcrypt.MinCost)
	user2 := user{"ken", bPassword, "ken", "tan", false}
	mapUsers[user2.Username] = user2
	bPassword, _ = bcrypt.GenerateFromPassword([]byte("ivan123lim"), bcrypt.MinCost)
	user3 := user{"ivan", bPassword, "ivan", "lim", false}
	mapUsers[user3.Username] = user3
	// will create default admin acct in the server
	// do not use this in go in action 2. security breach
	// will be taught in go in action 2 how to code these 2 lines
	bPassword, _ = bcrypt.GenerateFromPassword([]byte("ad123min"), bcrypt.MinCost)
	admin1 := user{"admin", bPassword, "the", "admin", true}
	mapUsers[admin1.Username] = admin1

	// load all templates from template folder
	// var err error
	// tpl, err = template.ParseGlob("templates/*")
	tpl = template.Must(template.ParseGlob("templates/*"))
	// if err != nil {
	// 	fmt.Println(err)
	// }

}

func StartServer() {
	http.HandleFunc("/", index)
	// http.HandleFunc("/restricted", restricted)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/browse", browse)
	http.HandleFunc("/search", searchAppt)
	http.HandleFunc("/make", makeAppt)
	http.HandleFunc("/list", list)
	http.HandleFunc("/edit", editAppt)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	cleanData()
	// fmt.Println(userMonth, userDate, apptMonthStr, userApptDate, userApptMonth, dentistResult)
	myUser = getUser(res, req)
	tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

// func restricted(res http.ResponseWriter, req *http.Request) {
// 	myUser := getUser(res, req)
// 	if !alreadyLoggedIn(req) {
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}
// 	tpl.ExecuteTemplate(res, "restricted.gohtml", myUser)
// }

func signup(res http.ResponseWriter, req *http.Request) {

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
				http.Error(res, "Username already taken", http.StatusForbidden)
				// need to return as there is further processing
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)

			// set the cookie value to the username
			// so we know which username using which cookie value
			mapSessions[myCookie.Value] = username

			// 1st argument pass in password in []byte format]
			// 2nd arguemnt determine how difficult the decrytion will be
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				// need to return as there is further processing
				return
			}

			myUser = user{username, bPassword, firstName, lastName, false}
			mapUsers[username] = myUser
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

func login(res http.ResponseWriter, req *http.Request) {

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
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

func logout(res http.ResponseWriter, req *http.Request) {

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

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func getUser(res http.ResponseWriter, req *http.Request) user {

	// get current session cookie
	// cookie name can be anything
	// here we used myCookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
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

func browse(res http.ResponseWriter, req *http.Request) {

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
	// var Error error

	if req.Method == http.MethodPost {

		dentistName := req.FormValue("dentist")

		if dentistName != "" {
			Output, _ = userpackage.BrowseDrAppointment(dentistName)
		}
	}

	// if Error != nil {
	// 	http.Error(res, "Wrong dentist name", http.StatusBadRequest)
	// 	return
	// }
	tpl.ExecuteTemplate(res, "browse.gohtml", Output)
}

func searchAppt(res http.ResponseWriter, req *http.Request) {

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
	// var Error error
	// fmt.Println("in search")
	if req.Method == http.MethodPost {

		idStr := req.FormValue("apptId")
		id, err := strconv.Atoi(idStr)

		// fmt.Printf("%v, %T", id, id)
		if err != nil {
			http.Error(res, "Invalid appointment id has been entered", http.StatusBadRequest)
			return
		}

		Output = userpackage.SearchAppointment(id)

	}

	// if Error != nil {
	// 	http.Error(res, "Wrong dentist name", http.StatusBadRequest)
	// 	return
	// }
	tpl.ExecuteTemplate(res, "search.gohtml", Output)
}

func makeAppt(res http.ResponseWriter, req *http.Request) {

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

	ApptInfo = appt{datatype.TimeArr[:], "", 0, datatype.ApptYear, "", []string{}, ""}

	if req.Method == http.MethodPost {
		// timeSession = req.FormValue("timeSession")
		if userMonth == "" && userDay == "" {
			userDay = req.FormValue("day")
			userMonth = req.FormValue("month")

			userApptMonth, userApptDay, err = checkValidInput(res, userMonth, userDay)
			if err != nil {
				http.Error(res, "Invalid month/day has been entered. Month need to be integer between 1 and 12. Day need to be integer between 1 and 31.", http.StatusBadRequest)
				return
			}
			apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

			timeSession = req.FormValue("time")
			userTimeSession, err = strconv.Atoi(timeSession)
			timeSession = datatype.TimeArr[userTimeSession-1]
			if err != nil {
				http.Error(res, "Invalid times ession has been entered. Time session need to be integer between 1 and 15", http.StatusBadRequest)
				return
			}
			// fmt.Println(day, month)

			dentistResult = userpackage.MakeAppointment(userApptMonth, userApptDay, userTimeSession)
			// dentistResult = searchAvailDentist(userApptMonth, userApptDay, userTimeSession, false, myUser.Username)
			// fmt.Println(dentistResult)
			// result := addToApptList(userApptMonth, userApptDay, userTimeSession, myUser.Username)
			ApptInfo = appt{datatype.TimeArr[:], apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, ""}
		} else {
			dentistName := req.FormValue("dentist")
			dentistName = utils.ConvertToUpper(dentistName)
			// Appt = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
			// fmt.Println(dentistName)
			if dentistName != "" {
				// fmt.Println("enter")
				// fmt.Println(userApptMonth, userApptDate)
				// timeSlot, err := DentistHash.listAvailTime(dentistName, userApptMonth, userApptDay)
				result := userpackage.AddToApptList(userApptMonth, userApptDay, userTimeSession, dentistName, myUser.Username)
				if err != nil {
					panic(err)
				}

				ApptInfo = appt{datatype.TimeArr[:], apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, result}
				// fmt.Println(Appt)
			}
		}

	}

	tpl.ExecuteTemplate(res, "make.gohtml", ApptInfo)
}

func checkValidInput(res http.ResponseWriter, m string, d string) (int, int, error) {

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
	// if err != nil {
	// 	http.Error(res, "Invalid date has been entered. Date need to be integer between 1 and 31", http.StatusBadRequest)
	// 	return
	// }

	return mInt, dInt, err

}

func list(res http.ResponseWriter, req *http.Request) {

	// cleanData()
	// fmt.Println(userMonth, userDate, apptMonthStr, userApptDate, userApptMonth, dentistResult)

	// var Appt userDate
	// // var timeSlot []string
	// var month, date, apptMonthStr string
	// var apptDate, apptMonth, apptYear int
	// var dentist []string

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
			userApptMonth, err = strconv.Atoi(userMonth)
			if err != nil {
				http.Error(res, "Invalid month has been entered. Month need to be integer between 1 and 12", http.StatusBadRequest)
				return
			}

			userApptDay, err = strconv.Atoi(userDay)
			if err != nil {
				http.Error(res, "Invalid date has been entered. Date need to be integer between 1 and 31", http.StatusBadRequest)
				return
			}

			apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

			// apptYear := time.Now().Year()
			// Appt = userDate{apptMonthStr, apptDate, apptYear, []string{}}

			dentistResult, err = userpackage.ListAvailableDentistTime(userApptMonth, userApptDay)
			if err != nil {
				panic(err)
			}
			// fmt.Println(dentist)
			Appt = userData{apptMonthStr, userApptDay, datatype.ApptYear, dentistResult, "", []string{}}
		} else {

			// var timeSlot []string
			// var err error
			// month := req.FormValue("month")
			// date := req.FormValue("date")

			// // convert from string to int
			// // if input is not integer return error
			// apptMonth, err := strconv.Atoi(month)
			// if err != nil {
			// 	http.Error(res, "Invalid month has been entered. Month need to be integer between 1 and 12", http.StatusBadRequest)
			// 	return
			// }

			// apptDate, err := strconv.Atoi(date)
			// if err != nil {
			// 	http.Error(res, "Invalid date has been entered. Date need to be integer between 1 and 31", http.StatusBadRequest)
			// 	return
			// }

			// apptMonthStr, _ := ConvertIntToMonth(apptMonth)

			// apptYear := time.Now().Year()
			// // Appt = userDate{apptMonthStr, apptDate, apptYear, []string{}}

			// dentist := listAvailableDentistTime(apptMonth, apptDate)
			// // fmt.Println(dentist)
			// Appt = userDate{apptMonthStr, apptDate, apptYear, dentist}

			dentistName := req.FormValue("dentist")
			dentistName = utils.ConvertToUpper(dentistName)
			// Appt = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
			// fmt.Println(dentistName)
			if dentistName != "" {
				// fmt.Println("enter")
				// fmt.Println(userApptMonth, userApptDate)
				timeSlot, err = datatype.DentistHash.ListAvailTime(dentistName, userApptMonth, userApptDay)
				if err != nil {
					panic(err)
				}

				Appt = userData{apptMonthStr, userApptDay, datatype.ApptYear, dentistResult, dentistName, timeSlot}
				// fmt.Println(Appt)
			}
		}
		// dentistName := req.FormValue("dentist")
		// dentistName = ConvertToUpper(dentistName)
		// Appt = userDate{apptMonthStr, apptDate, apptYear, dentist, dentistName, []string{}}
		// fmt.Println(dentistName)
		// if dentistName != "" {
		// 	fmt.Println("enter")
		// 	timeSlot, err := DentistHash.listAvailTime(dentistName, apptMonth, apptDate)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	Appt = userDate{apptMonthStr, apptDate, apptYear, dentist, dentistName, timeSlot}
		// }
		// if month != "" && date != "" {
		// 	http.Redirect(res, req, "/ask", http.StatusSeeOther)
		// 	return
		// }

	}

	tpl.ExecuteTemplate(res, "list.gohtml", Appt)

}

func editAppt(res http.ResponseWriter, req *http.Request) {

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

	EditApptInfo = apptEdit{datatype.TimeArr[:], userId, "", apptMonthStr, userApptDay, datatype.ApptYear, timeSession, []string{}, "", ""}

	if req.Method == http.MethodPost {
		// timeSession = req.FormValue("timeSession")
		if userId == "" {
			userId = req.FormValue("apptId")

			userApptId, err = strconv.Atoi(userId)
			if err != nil {
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
					panic(errors.New("wrong appointment id entered"))
				}
			}
			// fmt.Println(details, outputMessage)
			// dentistResult = searchAvailDentist(userApptMonth, userApptDay, userTimeSession, false, myUser.Username)
			// fmt.Println(dentistResult)
			// result := addToApptList(userApptMonth, userApptDay, userTimeSession, myUser.Username)
			EditApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, []string{}, "", ""}
		} else {

			if userMonth == "" && userDay == "" {
				userDay = req.FormValue("day")
				userMonth = req.FormValue("month")

				userApptMonth, userApptDay, err = checkValidInput(res, userMonth, userDay)
				if err != nil {
					http.Error(res, "Invalid month/day has been entered. Month need to be integer between 1 and 12. Day need to be integer between 1 and 31.", http.StatusBadRequest)
					return
				}
				apptMonthStr, _ = utils.ConvertIntToMonth(userApptMonth)

				timeSession = req.FormValue("time")
				userTimeSession, err = strconv.Atoi(timeSession)
				timeSession = datatype.TimeArr[userTimeSession-1]
				if err != nil {
					http.Error(res, "Invalid times ession has been entered. Time session need to be integer between 1 and 15", http.StatusBadRequest)
					return
				}

				dentistResult = userpackage.SearchAvailDentist(userApptMonth, userApptDay, userTimeSession, true)
				EditApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, "", ""}
			} else {
				dentistName := req.FormValue("dentist")
				dentistName = utils.ConvertToUpper(dentistName)

				// EditApptInfo = apptEdit{TimeArr[:], outputMessage, apptMonthStr, userApptDay, ApptYear, timeSession, dentistResult, dentistName, ""}
				// Appt = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
				// fmt.Println(dentistName)
				if dentistName != "" {
					// fmt.Println("enter")
					// fmt.Println(userApptMonth, userApptDate)
					if _, err := details.CheckChanges(userTimeSession, userApptMonth, userApptDay, dentistName); err != nil {
						// fmt.Println("panic")
						panic(err)
					} else {
						var result string
						if myUser.Username == "admin" { // if admin edit appointment, use back patient username
							result = userpackage.UpdateAppt(details, userTimeSession, userApptMonth, userApptDay, dentistName, userApptId, details.Patient())
						} else { // else use current user username
							result = userpackage.UpdateAppt(details, userTimeSession, userApptMonth, userApptDay, dentistName, userApptId, myUser.Username)
						}
						EditApptInfo = apptEdit{datatype.TimeArr[:], userId, outputMessage, apptMonthStr, userApptDay, datatype.ApptYear, timeSession, dentistResult, dentistName, result}
					}

					// Appt = userData{apptMonthStr, userApptDay, ApptYear, dentistResult, dentistName, timeSlot}
					// fmt.Println(Appt)
				}
			}

			// dentistName := req.FormValue("dentist")
			// dentistName = ConvertToUpper(dentistName)
			// // Appt = userData{apptMonthStr, userApptDate, ApptYear, dentist, dentistName, []string{}}
			// // fmt.Println(dentistName)
			// if dentistName != "" {
			// 	// fmt.Println("enter")
			// 	// fmt.Println(userApptMonth, userApptDate)
			// 	// timeSlot, err := DentistHash.listAvailTime(dentistName, userApptMonth, userApptDay)
			// 	result := addToApptList(userApptMonth, userApptDay, userTimeSession, dentistName, myUser.Username)
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	// EditApptInfo = appt{TimeArr[:], apptMonthStr, userApptDay, ApptYear, timeSession, dentistResult, result}
			// 	// fmt.Println(Appt)
			// }
		}

	}

	tpl.ExecuteTemplate(res, "edit.gohtml", EditApptInfo)
}

func cleanData() {
	// everytime goes back to other route
	// clean data for listing
	userMonth, userDay, apptMonthStr, userId = "", "", "", ""
	userApptDay, userApptMonth = 0, 0
	dentistResult, timeSlot = []string{}, []string{}

	Appt = userData{}
	EditApptInfo = apptEdit{}
}

// func ask

// func edit(res http.ResponseWriter, req *http.Request) {

// 	if !alreadyLoggedIn(req) {
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}

// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Panic occurred in edit appointment:", err)
// 			http.Error(res, "Error occurred when editing appointment", http.StatusBadRequest)
// 			return
// 		}
// 	}()
// }
