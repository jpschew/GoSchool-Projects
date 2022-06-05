package userpackage

import (
	"errors"
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/datatype"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/utils"
	"sync"
)

// const (
// 	ApptYear = 2022
// )

var (
	// TimeArr = [15]string{
	// 	"09:00AM",
	// 	"09:30AM",
	// 	"10:00AM",
	// 	"10:30AM",
	// 	"11:00AM",
	// 	"11:30AM",
	// 	"01:00PM",
	// 	"01:30PM",
	// 	"02:00PM",
	// 	"02:30PM",
	// 	"03:00PM",
	// 	"03:30PM",
	// 	"04:00PM",
	// 	"04:30PM",
	// 	"05:00PM",
	// }

	patientMessage = [4]string{
		"1. Make an appointment",
		"2. List available times of selected doctor",
		"3. Edit appointment",
	}

	// invalidMessage = "You have entered an invalid choice."

	// ApptYear = time.Now().Year()
)

// func patientPage() {

// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Panic occurred in patient page:", err)
// 		}
// 	}()

// 	var userInput string

// 	PrintMessage(patientMessage[:])

// 	fmt.Println("Please enter you choice")
// 	fmt.Scanln(&userInput)

// 	switch userInput {
// 	case "1":
// 		ClearScreen()
// 		makeAppointment()
// 	case "2":
// 		ClearScreen()
// 		listAvailableDentistTime()
// 	case "3":
// 		ClearScreen()
// 		editAppointment()
// 	default:
// 		fmt.Println(invalidMessage)
// 	}

// }

func MakeAppointment(apptMonth int, apptDate int, apptTime int) []string {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in make appointment:", err)
		}
	}()

	// var apptMonth string
	// var apptMonth, apptDate, apptTime int

	// apptMonth, apptDate, apptTime = askForApptDate(false)

	result := SearchAvailDentist(apptMonth, apptDate, apptTime, false)

	return result

}

// commented out for go in action 1
// // if showAvailTime is true will show all available time for particular day
// func askForApptDate(showAvailTime bool) (int, int, int) {

// 	ClearScreen()

// 	// defer func() {
// 	// 	if err := recover(); err != nil {
// 	// 		fmt.Println("Panic occurred in ask appointment date:", err)
// 	// 	}
// 	// }()

// 	var apptMonth string
// 	var apptDate int
// 	var apptTime int

// 	// get the current year, month and date
// 	year, month, date := time.Now().Date()

// 	fmt.Printf("Today is %v %v %v.\n", date, month, year)
// 	fmt.Println("You appointment need to book one day in advance.")

// 	fmt.Println("Select the month that you want (from", month, "onwards):")

// 	if _, err := fmt.Scanln(&apptMonth); err != nil {
// 		panic(err)
// 	}
// 	apptMonth = ConvertToUpper(apptMonth)

// 	monthInt, err := ConvertMonthToInt(apptMonth)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Select the date that you want: ")
// 	fmt.Scanln(&apptDate)
// 	if apptDate < 1 || apptDate > 31 { // consider those with 31 days
// 		panic(errors.New("date is out of range"))
// 	} else {
// 		if apptDate == 31 && (monthInt == 4 || monthInt == 6 || monthInt == 9 || monthInt == 11) { // consider those with 30 days
// 			panic(errors.New("date is out of range"))
// 		} else if apptDate >= 29 && monthInt == 2 && !IsLeap(ApptYear) { // consider Feb without leap year
// 			panic(errors.New("date is out of range"))
// 		}
// 	}

// 	if !showAvailTime {
// 		fmt.Println("Select the time slot that you want")
// 		fmt.Println("------------------------------------")

// 		for i, timeslot := range TimeArr {
// 			fmt.Print(i+1, ".", timeslot, "\t")
// 			if (i+1)%5 == 0 {
// 				fmt.Println()
// 			}
// 		}

// 		fmt.Scanln(&apptTime)
// 	}

// 	ClearScreen()

// 	return monthInt, apptDate, apptTime
// }

// if isEdit is true means for editing appointment
// else is for making appointment
func SearchAvailDentist(month int, date int, timeSession int, isEdit bool) []string {

	// ClearScreen()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search available dentist:", err)
		}
	}()

	var emptyDentistList bool
	var dentistResult []string
	var result []string

	fmt.Println("These are the available dentists on your selected timeslot:")
	fmt.Println("------------------------------------------------------------")

	s := fmt.Sprintln("These are the available dentists on your selected timeslot:")
	dentistResult = append(result, s)

	emptyDentistList, result = printAvailDentist(month, date, timeSession, true)

	// if no available dentist, prompt user to choose another timeslot
	if emptyDentistList {
		fmt.Println("There are no dentists available on your chosen slot. Please choose another date or time slot.")

		s = fmt.Sprintln("There are no dentists available on your chosen slot. Please choose another date or time slot.")
		dentistResult = append(dentistResult, s)

	} else {

		// commented away for go ina ction 1
		// if !isEdit{ // for making appointment only as edit appointment will not add to the appointment list
		// addToApptList(month, date, timeSession, userName)
		// }

		// added for go in action 1
		// edit or add appointment also need this list
		dentistResult = append(dentistResult, result...)

	}
	return dentistResult

}

// if isSearch is true is for searching of available dentist on a particular timeslot
// else just print all dentists in clinic
func printAvailDentist(month int, date int, timeSession int, isSearch bool) (bool, []string) {

	// ClearScreen()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in print available dentist:", err)
		}
	}()

	var emptyDentistList = true
	var dentistList []string

	for _, dentist := range datatype.NewArr.List() { // NewArr is a pointer and the list field consists the updated dentist
		// if searching for dentist for a particular timeslot
		if isSearch {
			found, name, _ := datatype.DentistHash.Search(dentist, timeSession, month, date)
			if found {
				fmt.Println(name)
				emptyDentistList = false

				dentistList = append(dentistList, dentist)
			}
		} else { // if not searching then provide all the dentists in the clinic
			// in this case, we are printing out the dentist only, value of emptyDentistList will not affect
			fmt.Println(dentist)

			// added for go in action 1
			dentistList = append(dentistList, dentist)
		}

	}
	return emptyDentistList, dentistList

}

func AddToApptList(month int, date int, timeSession int, dentist string, userName string) string {

	// ClearScreen()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add to appointment list:", err)
		}
	}()

	// userChoice := askForDentistName()

	wg := sync.WaitGroup{}

	apptIdChn := make(chan int)
	updateTime := make(chan int, 1)

	wg.Add(3)

	// need handle race condition issue for updating
	// pass in buffered chan with length 1 so only one can update at a time
	// as this will delete a timenode from timeBST
	// might clash when others editing appointment
	go datatype.DentistHash.UpdateTimeSlot(dentist, timeSession, month, date, updateTime, &wg)
	updateTime <- 1
	close(updateTime)

	// can use goroutines to add new appointment to appointment list and dr appointment list concurrently
	go datatype.ApptHash.Add(dentist, datatype.DentistHash.GetDrId(dentist), timeSession, month, date, apptIdChn, true, &wg, userName)
	apptId := <-apptIdChn
	go datatype.DrApptHash.Add(dentist, datatype.DentistHash.GetDrId(dentist), timeSession, month, date, apptIdChn, false, &wg, userName)
	close(apptIdChn)

	result := printAppt(apptId, dentist, timeSession, date, month, false)

	wg.Add(1)
	go utils.SendingEmail(&wg)

	// wait for all goroutines to finish
	wg.Wait()
	return result
}

// commented out for go in action 1
// func askForDentistName() string {

// 	// ClearScreen()

// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Panic occurred in ask for dentist name:", err)
// 		}
// 	}()

// 	var userChoice string
// 	fmt.Println("Please choose the doctor that you want to make an appointment with: ")
// 	fmt.Scanln(&userChoice)
// 	userChoice = ConvertToUpper(userChoice)
// 	return userChoice
// }

// // if isUpadate is true will update the appointment
// else it will mean that new appontment had been made
func printAppt(apptId int, name string, apptTime int, apptDate int, apptMonth int, isUpdate bool) string {

	var result string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in print appointment:", err)
		}
	}()

	if isUpdate {
		fmt.Printf("Your appointment id is %d and your appointment has been updated to be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, datatype.ApptYear, apptMonth, apptDate, datatype.TimeArr[apptTime-1], name)

		result = fmt.Sprintf("Your appointment id is %d and your appointment has been updated to be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, datatype.ApptYear, apptMonth, apptDate, datatype.TimeArr[apptTime-1], name)
	} else {
		fmt.Println(apptId, datatype.ApptYear, apptMonth, apptDate, apptTime-1, name)
		fmt.Printf("Your appointment id is %d and your appointment will be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, datatype.ApptYear, apptMonth, apptDate, datatype.TimeArr[apptTime-1], name)

		result = fmt.Sprintf("Your appointment id is %d and your appointment will be on %d-%02d-%02d %s with Dr. %s.\n",
			apptId, datatype.ApptYear, apptMonth, apptDate, datatype.TimeArr[apptTime-1], name)
	}

	return result

}

// modified to return []string for go in action 1
func ListAvailableDentistTime(apptMonth int, apptDate int) ([]string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in list available dentist time:", err)
	// 	}
	// }()

	// ClearScreen()

	var dentistList []string

	// commented out for go in action 1
	// apptMonth, apptDate, _ := askForApptDate(true)

	if checkForValidDate(apptMonth, apptDate) {
		fmt.Println("These are the available dentists in our clinic:")
		fmt.Println("-----------------------------------------------")

		s := fmt.Sprintln("These are the available dentists in our clinic:")
		dentistList = append(dentistList, s)

		// added dentistList for go in action 1
		_, result := printAvailDentist(apptMonth, apptDate, 0, false)
		dentistList = append(dentistList, result...)

		// commented out for go in action 1
		// dentistName := askForDentistName()

		// ClearScreen()

		// DentistHash.listAvailTime(dentistName, apptMonth, apptDate)

		return dentistList, nil
	} else {
		return []string{}, errors.New("date is out of range")
	}

}

// added in go in action 1
func checkForValidDate(apptMonth int, apptDate int) bool {
	// if apptMonth < 1 || apptMonth > 12 {
	// 	return false
	// }
	if apptMonth < 1 || apptMonth > 12 || apptDate < 1 || apptDate > 31 { // consider those with 31 days
		// panic(errors.New("date is out of range"))
		return false
	} else {
		if apptDate >= 31 && (apptMonth == 4 || apptMonth == 6 || apptMonth == 9 || apptMonth == 11) { // consider those with 30 days
			// panic(errors.New("date is out of range"))
			return false
		} else if apptDate >= 29 && apptMonth == 2 && !utils.IsLeap(datatype.ApptYear) { // consider Feb without leap year
			// panic(errors.New("date is out of range"))
			return false
		}
	}
	return true
}

func EditAppointment(apptId int) (*datatype.ApptDetails, string) {

	// var apptId int

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in edit appointment:", err)
		}
	}()

	// fmt.Println("Please enter you Appointment Id that you want to edit:")
	// fmt.Scanln(&apptId)

	// modifiedd for go in action 1
	_, apptDetails, output, err := datatype.ApptHash.Search(apptId)
	if err != nil {
		panic(err)
		// fmt.Println(err)
	}

	// askForPatientChoice(apptDetails, apptId, userName)
	return apptDetails, output

}

// func askForPatientChoice(details *apptDetails, apptId int, userName string) {

// 	var apptMonth int
// 	var apptDate int
// 	var apptTime int
// 	var userChoice string

// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Println("Panic occurred in ask for patient choice:", err)
// 		}
// 	}()

// 	fmt.Println("Do you still want to change your appointment? (Y/N)")

// 	if _, err := fmt.Scanln(&userChoice); err != nil {
// 		panic(err)
// 	} else {
// 		userChoice = ConvertToUpper(userChoice)

// 		if userChoice == "Y" {

// 			// commented out for go in action 1
// 			// apptMonth, apptDate, apptTime = askForApptDate(false)

// 			searchAvailDentist(apptTime, apptMonth, apptDate, true)

// 			// commented out for go in action 1
// 			// userChoice = askForDentistName()

// 			if _, err := details.CheckChanges(apptTime, apptMonth, apptDate, userChoice); err != nil {
// 				// fmt.Println("panic")
// 				panic(err)
// 			} else {
// 				updateAppt(details, apptTime, apptMonth, apptDate, userChoice, apptId, userName)
// 			}
// 		} else if userChoice == "N" {
// 			return
// 		} else {
// 			panic(errors.New("enter wrong choice"))
// 		}
// 	}
// }

// modified to return a string for go in action 1
func UpdateAppt(details *datatype.ApptDetails, apptTime int, apptMonth int, apptDate int, name string, apptId int, userName string) string {

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
		go datatype.DentistHash.AddTimeSlot(details.Name(), details.Time(), details.Month(), details.Day(), &wg)

		// need handle race condition issue for updating
		// pass in buffered chan with length 1 so only one can update at a time
		// as this will delete a timenode from timeBST
		// might clash when others making appointment
		go datatype.DentistHash.UpdateTimeSlot(name, apptTime, apptMonth, apptDate, updateTime, &wg)
		updateTime <- 1
		close(updateTime)

		go datatype.ApptHash.Update(name, datatype.DentistHash.GetDrId(name), apptTime, apptMonth, apptDate, apptId, true, &wg, userName)

		if details.Name() == name { // for updating if dentist chosen is the same
			go datatype.DrApptHash.Update(name, datatype.DentistHash.GetDrId(name), apptTime, apptMonth, apptDate, apptId, false, &wg, userName)
		} else { // for updating if dentist chosen is the different
			go datatype.DrApptHash.UpdateDiffDr(name, datatype.DentistHash.GetDrId(name), datatype.DentistHash.GetDrId(details.Name()), apptTime, apptMonth, apptDate, apptId, false, &wg, userName)
		}

		// end := time.Now()

		// diff := time.Since(start)
		// fmt.Println(diff.Seconds())
	}

	result := printAppt(apptId, name, apptTime, apptMonth, apptDate, true)

	wg.Add(1)
	go utils.SendingEmail(&wg)

	// wait for all goroutines to finish
	wg.Wait()
	return result
}
