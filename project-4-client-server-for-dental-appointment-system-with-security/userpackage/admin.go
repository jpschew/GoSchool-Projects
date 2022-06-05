// Package userpackage implements the different functions/features used by admin and patient.
// The functions/features that are associated with admin and patient are in admin.go and patient.go respectively.
package userpackage

import (
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/datatype"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/utils"
)

var (
	adminMessage = [2]string{
		"1. Browse appointments for a dentist",
		"2. Search for appointment",
	}
)

// BrowseDrAppointment takes in a dentist name and return his/her appointment list.
// It will take in string as an input and return a slice of string as an output.
// If the dentist is not available, an error will be return as well.
func BrowseDrAppointment(dentist string) ([]string, error) {

	var output []string

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in browse appointment:", err)
	// 	}

	// }()
	dentist = utils.ConvertToUpper(dentist)

	fmt.Printf("\nThe appointment list for Dr. %s is shown below:\n", dentist)

	result, err := datatype.DrApptHash.Browse(datatype.DentistHash.DrId(dentist))

	if err != nil {
		panic(err)
	}

	message := fmt.Sprintf("\nThe appointment list for Dr. %s is shown below:\n", dentist)
	output = append(output, message)

	if err != nil {
		panic(err)
	}

	output = append(output, result...)

	return output, err

}

// SearchAppointment takes in an Appointment Id and return a message with the appointment details.
// It will take in integer as an input and return a string as an output.
// If the Appointment Id is invalid, an error will be return as well.
func SearchAppointment(id int) (string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in search appointment:", err)
	// 	}
	// }()

	_, _, output, err := datatype.ApptHash.Search(id)

	if err != nil {
		panic(err)
	}

	return output, err

}
