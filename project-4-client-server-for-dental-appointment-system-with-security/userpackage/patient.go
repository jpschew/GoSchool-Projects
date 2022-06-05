package userpackage

import (
	"errors"
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/datatype"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/utils"
	"sync"
)

var (
	patientMessage = [4]string{
		"1. Make an appointment",
		"2. List available times of selected doctor",
		"3. Edit appointment",
	}
)

func MakeAppointment(apptMonth int, apptDate int, apptTime int) []string {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in make appointment:", err)
		}
	}()

	result := SearchAvailDentist(apptMonth, apptDate, apptTime, false)

	return result

}

// SearchAvailDentist will search for available dentist given a date and time.
// It will take in integer inputs for month, day, time and return a slice of string as an output.
// If isEdit is true means for editing appointment, else is for making appointment.
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

// AddToApptList will add the new appointment to the system appointment list and return message of the appointment details as a string.
// It will take in integer inputs for month, day, time, string input for dentist name and username and return a string as an output.
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
	go datatype.ApptHash.Add(dentist, datatype.DentistHash.DrId(dentist), timeSession, month, date, apptIdChn, true, &wg, userName)
	apptId := <-apptIdChn
	go datatype.DrApptHash.Add(dentist, datatype.DentistHash.DrId(dentist), timeSession, month, date, apptIdChn, false, &wg, userName)
	close(apptIdChn)

	result := printAppt(apptId, dentist, timeSession, date, month, false)

	wg.Add(1)
	go utils.SendingEmail(&wg)

	// wait for all goroutines to finish
	wg.Wait()
	return result
}

// if isUpadate is true will update the appointment
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

// ListAvailableDentist will check the user key in the valid date and return list of available dentist as a slice.
// It will take in integer inputs for month, day and return a slice as an output.
// If the date is not a valid date, an error will be return as well.
func ListAvailableDentist(apptMonth int, apptDate int) ([]string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in list available dentist time:", err)
	// 	}
	// }()

	// ClearScreen()

	var dentistList []string

	// commented out for go in action 1
	// apptMonth, apptDate, _ := askForApptDate(true)

	if utils.CheckForValidDate(apptMonth, apptDate) {
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

		// // to list the available dentist time
		// DentistHash.listAvailTime(dentistName, apptMonth, apptDate)

		return dentistList, nil
	} else {
		return []string{}, errors.New("date is out of range")
	}

}

// EditAppointment will ask the user for the Appointment Id that she/he wants to edit and return a pointer to the appointment details as well as the appointmen details as a string.
// It will take in Appointment Id as integer input and return a pointer and a string.
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

// UpdateAppt will update the appointment list given the appointment details and return a message of the appointment details as a string.
// It will take in appointment details pointer, the appointment id, month, date and time slot as integer input as well as dentist name and username as string input and return a string as output.
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
		go datatype.DentistHash.AddTimeSlot(details.DentistName(), details.Time(), details.Month(), details.Day(), &wg)

		// need handle race condition issue for updating
		// pass in buffered chan with length 1 so only one can update at a time
		// as this will delete a timenode from timeBST
		// might clash when others making appointment
		go datatype.DentistHash.UpdateTimeSlot(name, apptTime, apptMonth, apptDate, updateTime, &wg)
		updateTime <- 1
		close(updateTime)

		go datatype.ApptHash.Update(name, datatype.DentistHash.DrId(name), apptTime, apptMonth, apptDate, apptId, true, &wg, userName)

		if details.DentistName() == name { // for updating if dentist chosen is the same
			go datatype.DrApptHash.Update(name, datatype.DentistHash.DrId(name), apptTime, apptMonth, apptDate, apptId, false, &wg, userName)
		} else { // for updating if dentist chosen is the different
			go datatype.DrApptHash.UpdateDiffDr(name, datatype.DentistHash.DrId(name), datatype.DentistHash.DrId(details.DentistName()), apptTime, apptMonth, apptDate, apptId, false, &wg, userName)
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
