package main

import (
	"fmt"
)

var (
	adminMessage = [2]string{
		"1. Browse appointments for a dentist",
		"2. Search for appointment",
	}

	// invalidMessage = "You have entered an invalid choice."
)

func adminPage() {

	var adminInput string
	PrintMessage(adminMessage[:])

	fmt.Println("Please enter you choice")
	fmt.Scanln(&adminInput)

	switch adminInput {
	case "1":
		ClearScreen()
		browseDrAppointment()
	case "2":
		ClearScreen()
		searchAppointment()
	default:
		fmt.Println(invalidMessage)
	}
}

func browseDrAppointment() {

	var dentist string

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in browse appointment:", err)
		}
	}()

	fmt.Println("Please enter the dentist that you want to browse:")
	fmt.Scanln(&dentist)
	dentist = ConvertToUpper(dentist)

	if err := DrApptHash.Browse(DentistHash.getDrId(dentist)); err != nil {
		panic(err)
	}

	fmt.Printf("\nThe appointment list for Dr. %s is shown below:\n", dentist)

	// if err := DrApptHash.Browse(DentistHash.getDrId(dentist)); err != nil {
	// 	panic(err)
	// }

	Waiting()

}

func searchAppointment() {

	var id int

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search appointment:", err)
		}
	}()

	fmt.Println("Please enter the appointment Id that you want to search:")
	fmt.Scanln(&id)

	if _, _, err := ApptHash.Search(id); err != nil {
		panic(err)
	}

	Waiting()

}
