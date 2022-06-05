package userpackage

import (
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/datatype"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/utils"
)

var (
	adminMessage = [2]string{
		"1. Browse appointments for a dentist",
		"2. Search for appointment",
	}

	// invalidMessage = "You have entered an invalid choice."
)

// comment out in go in action 1
// func adminPage() {

// 	var adminInput string
// 	PrintMessage(adminMessage[:])

// 	fmt.Println("Please enter you choice")
// 	fmt.Scanln(&adminInput)

// 	switch adminInput {
// 	case "1":
// 		ClearScreen()
// 		browseDrAppointment()
// 	case "2":
// 		ClearScreen()
// 		searchAppointment()
// 	default:
// 		fmt.Println(invalidMessage)
// 	}
// }

// modified for go in action 1 to add in parameter
// and return []string
func BrowseDrAppointment(dentist string) ([]string, error) {

	// added for go in action 1
	var output []string

	// comment out for go in action 1
	// var dentist string

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in browse appointment:", err)
	// 	}
	// }()

	// fmt.Println("Please enter the dentist that you want to browse:")
	// fmt.Scanln(&dentist)
	dentist = utils.ConvertToUpper(dentist)

	fmt.Printf("\nThe appointment list for Dr. %s is shown below:\n", dentist)

	result, err := datatype.DrApptHash.Browse(datatype.DentistHash.GetDrId(dentist))

	if err != nil {
		panic(err)
	}

	// fmt.Printf("\nThe appointment list for Dr. %s is shown below:\n", dentist)

	// if output, err := DrApptHash.Browse(DentistHash.getDrId(dentist)); err != nil {
	// 	panic(err)
	// }

	// Waiting()

	// for go in action 1
	message := fmt.Sprintf("\nThe appointment list for Dr. %s is shown below:\n", dentist)
	output = append(output, message)

	// result, err := DrApptHash.Browse(DentistHash.getDrId(dentist))

	if err != nil {
		panic(err)
	}

	output = append(output, result...)

	return output, err

}

// modified for go in action 1 to add in parameter
// and return string and error
func SearchAppointment(id int) string {

	// comment out for go in action 1
	// var id int

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in search appointment:", err)
	// 	}
	// }()

	// fmt.Println("Please enter the appointment Id that you want to search:")
	// fmt.Scanln(&id)

	// if _, _, err := ApptHash.Search(id); err != nil {
	// 	panic(err)
	// }

	// Waiting()

	_, _, output, err := datatype.ApptHash.Search(id)

	if err != nil {
		panic(err)
	}

	return output

}
