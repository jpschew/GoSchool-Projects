package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ApptYear = 2022
)

var (
	TimeArr = [15]string{
		"09:00AM",
		"09:30AM",
		"10:00AM",
		"10:30AM",
		"11:00AM",
		"11:30AM",
		"01:00PM",
		"01:30PM",
		"02:00PM",
		"02:30PM",
		"03:00PM",
		"03:30PM",
		"04:00PM",
		"04:30PM",
		"05:00PM",
	}

	patientMessage = [4]string{
		"1. Make an appointment",
		"2. List available times of selected doctor",
		"3. Edit appointment",
	}

	invalidMessage = "You have entered an invalid choice."
)

func patientPage() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in patient page:", err)
		}
	}()

	var userInput string

	PrintMessage(patientMessage[:])

	fmt.Println("Please enter you choice")
	fmt.Scanln(&userInput)

	switch userInput {
	case "1":
		ClearScreen()
		makeAppointment()
	case "2":
		ClearScreen()
		listAvailableDentistTime()
	case "3":
		ClearScreen()
		editAppointment()
	default:
		fmt.Println(invalidMessage)
	}

}

func makeAppointment() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in make appointment:", err)
		}
	}()

	// var apptMonth string
	var apptMonth, apptDate, apptTime int

	apptMonth, apptDate, apptTime = askForApptDate(false)

	searchAvailDentist(apptMonth, apptDate, apptTime, false)

}

// if showAvailTime is true will show all available time for particular day
func askForApptDate(showAvailTime bool) (int, int, int) {

	ClearScreen()

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in ask appointment date:", err)
	// 	}
	// }()

	var apptMonth string
	var apptDate int
	var apptTime int

	// get the current year, month and date
	year, month, date := time.Now().Date()

	fmt.Printf("Today is %v %v %v.\n", date, month, year)
	fmt.Println("You appointment need to book one day in advance.")

	fmt.Println("Select the month that you want (from", month, "onwards):")

	if _, err := fmt.Scanln(&apptMonth); err != nil {
		panic(err)
	}
	apptMonth = ConvertToUpper(apptMonth)

	monthInt, err := ConvertMonthToInt(apptMonth)

	if err != nil {
		panic(err)
	}

	fmt.Println("Select the date that you want: ")
	fmt.Scanln(&apptDate)
	if apptDate < 1 || apptDate > 31 || monthInt < 1 || monthInt > 12 { // consider those with 31 days
		panic(errors.New("date is out of range"))
	} else {
		if apptDate >= 31 && (monthInt == 4 || monthInt == 6 || monthInt == 9 || monthInt == 11) { // consider those with 30 days
			panic(errors.New("date is out of range"))
		} else if apptDate >= 29 && monthInt == 2 && !IsLeap(ApptYear) { // consider Feb without leap year
			panic(errors.New("date is out of range"))
		}
	}

	if !showAvailTime {
		fmt.Println("Select the time slot that you want")
		fmt.Println("------------------------------------")

		for i, timeslot := range TimeArr {
			fmt.Print(i+1, ".", timeslot, "\t")
			if (i+1)%5 == 0 {
				fmt.Println()
			}
		}

		fmt.Scanln(&apptTime)
	}

	ClearScreen()

	return monthInt, apptDate, apptTime
}

// if isEdit is true means for editing appointment
// else is for making appointment
func searchAvailDentist(month int, date int, timeSession int, isEdit bool) {

	// ClearScreen()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search available dentist:", err)
		}
	}()

	var emptyDentistList bool

	fmt.Println("These are the available dentists on your selected timeslot:")
	fmt.Println("------------------------------------------------------------")

	emptyDentistList = printAvailDentist(month, date, timeSession, true)

	// if no available dentist, prompt user to choose another timeslot
	if emptyDentistList {
		fmt.Println("There are no dentists available on your chosen slot. Please choose another date or time slot.")

	} else if !isEdit { // for making appointment only as edit appointment will not add to the appointment list
		addToApptList(month, date, timeSession)
	}

}

// if isSearch is true is for searching of available dentist on a particular timeslot
// else just print all dentists in clinic
func printAvailDentist(month int, date int, timeSession int, isSearch bool) bool {

	// ClearScreen()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in print available dentist:", err)
		}
	}()

	var emptyDentistList = true

	for _, dentist := range (*NewArr).list { // NewArr is a pointer and the list field consists the updated dentist
		// if searching for dentist for a particular timeslot
		if isSearch {
			found, name, _ := DentistHash.Search(dentist, timeSession, month, date)
			if found {
				fmt.Println(name)
				emptyDentistList = false
			}
		} else { // if not searching then provide all the dentists in the clinic
			// in this case, we are printing out the dentist only, value of emptyDentistList will not affect
			fmt.Println(dentist)
		}

	}
	return emptyDentistList

}

func addToApptList(month int, date int, timeSession int) {

	// ClearScreen()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add to appointment list:", err)
		}
	}()

	userChoice := askForDentistName()

	wg := sync.WaitGroup{}

	apptIdChn := make(chan int)
	updateTime := make(chan int, 1)

	wg.Add(3)

	// need handle race condition issue for updating
	// pass in buffered chan with length 1 so only one can update at a time
	// as this will delete a timenode from timeBST
	// might clash when others editing appointment
	go DentistHash.UpdateTimeSlot(userChoice, timeSession, month, date, updateTime, &wg)
	updateTime <- 1
	close(updateTime)

	// can use goroutines to add new appointment to appointment list and dr appointment list concurrently
	go ApptHash.Add(userChoice, DentistHash.getDrId(userChoice), timeSession, month, date, apptIdChn, true, &wg)
	apptId := <-apptIdChn
	go DrApptHash.Add(userChoice, DentistHash.getDrId(userChoice), timeSession, month, date, apptIdChn, false, &wg)
	close(apptIdChn)

	printAppt(apptId, userChoice, timeSession, month, date, false)

	wg.Add(1)
	go SendingEmail(&wg)

	// wait for all goroutines to finish
	wg.Wait()
}

func askForDentistName() string {

	// ClearScreen()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in ask for dentist name:", err)
		}
	}()

	var userChoice string
	fmt.Println("Please choose the doctor that you want to make an appointment with: ")
	fmt.Scanln(&userChoice)
	userChoice = ConvertToUpper(userChoice)
	return userChoice
}

// // if isUpadate is true will update the appointment
// else it will mean that new appontment had been made
func printAppt(apptId int, name string, apptTime int, apptDate int, apptMonth int, isUpdate bool) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in print appointment:", err)
		}
	}()

	if isUpdate {
		fmt.Printf("Your appointment id is %d and your appointment has been updated to be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, ApptYear, apptMonth, apptDate, TimeArr[apptTime-1], name)
	} else {
		fmt.Printf("Your appointment id is %d and your appointment will be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, ApptYear, apptMonth, apptDate, TimeArr[apptTime-1], name)
	}

}

func listAvailableDentistTime() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in list available dentist time:", err)
		}
	}()

	ClearScreen()

	apptMonth, apptDate, _ := askForApptDate(true)

	fmt.Println("These are the available dentists in our clinic:")
	fmt.Println("-----------------------------------------------")

	_ = printAvailDentist(apptMonth, apptDate, 0, false)

	dentistName := askForDentistName()

	ClearScreen()

	DentistHash.listAvailTime(dentistName, apptMonth, apptDate)
}

func editAppointment() {

	var apptId int

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in edit appointment:", err)
		}
	}()

	fmt.Println("Please enter you Appointment Id that you want to edit:")
	fmt.Scanln(&apptId)

	_, apptDetails, err := ApptHash.Search(apptId)
	if err != nil {
		panic(err)
		// fmt.Println(err)
	}

	askForPatientChoice(apptDetails, apptId)

}

func askForPatientChoice(details *apptDetails, apptId int) {

	var apptMonth int
	var apptDate int
	var apptTime int
	var userChoice string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in ask for patient choice:", err)
		}
	}()

	fmt.Println("Do you still want to change your appointment? (Y/N)")

	if _, err := fmt.Scanln(&userChoice); err != nil {
		panic(err)
	} else {
		userChoice = ConvertToUpper(userChoice)

		if userChoice == "Y" {

			apptMonth, apptDate, apptTime = askForApptDate(false)

			searchAvailDentist(apptTime, apptMonth, apptDate, true)
			userChoice = askForDentistName()

			if _, err := details.CheckChanges(apptTime, apptMonth, apptDate, userChoice); err != nil {
				// fmt.Println("panic")
				panic(err)
			} else {
				updateAppt(details, apptTime, apptMonth, apptDate, userChoice, apptId)
			}
		} else if userChoice == "N" {
			return
		} else {
			panic(errors.New("enter wrong choice"))
		}
	}
}

func updateAppt(details *apptDetails, apptTime int, apptMonth int, apptDate int, name string, apptId int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in update appointment in patient page:", err)
		}
	}()

	wg := sync.WaitGroup{}

	if _, err := details.CheckChanges(apptTime, apptMonth, apptDate, name); err != nil {
		panic(err)
	} else {

		// addTime := make(chan int, 1)
		updateTime := make(chan int, 1)

		wg.Add(4)

		// start := time.Now()

		// can use goroutines here to do all jobs concurrently
		// the following 4 goroutines are for background computation, no need wait
		go DentistHash.AddTimeSlot(details.name, details.apptTime, details.apptMonth, details.apptDate, &wg)

		// need handle race condition issue for updating
		// pass in buffered chan with length 1 so only one can update at a time
		// as this will delete a timenode from timeBST
		// might clash when others making appointment
		go DentistHash.UpdateTimeSlot(name, apptTime, apptMonth, apptDate, updateTime, &wg)
		updateTime <- 1
		close(updateTime)

		go ApptHash.Update(name, DentistHash.getDrId(name), apptTime, apptMonth, apptDate, apptId, true, &wg)

		if details.name == name { // for updating if dentist chosen is the same
			go DrApptHash.Update(name, DentistHash.getDrId(name), apptTime, apptMonth, apptDate, apptId, false, &wg)
		} else { // for updating if dentist chosen is the different
			go DrApptHash.UpdateDiffDr(name, DentistHash.getDrId(name), DentistHash.getDrId(details.name), apptTime, apptMonth, apptDate, apptId, false, &wg)
		}

		// end := time.Now()

		// diff := time.Since(start)
		// fmt.Println(diff.Seconds())
	}

	printAppt(apptId, name, apptTime, apptMonth, apptDate, true)

	wg.Add(1)
	go SendingEmail(&wg)

	// wait for all goroutines to finish
	wg.Wait()
}
