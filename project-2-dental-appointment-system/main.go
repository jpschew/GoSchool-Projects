package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	toContinue = true
)

func init() {

	InitDentist()
	InitApptHashTable()
	InitAppt()
}

func main() {

	// 8 phyiscal processor
	runtime.GOMAXPROCS(runtime.NumCPU())
	// fmt.Println(runtime.NumCPU())

	// for panic recovery
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in  main:", err)
		}
	}()

	for toContinue {

		ClearScreen()
		checkForAdmin()
		toContinue = toContinueApp()

	}

}

func checkForAdmin() {

	var role string

	fmt.Println("Please state if you are a patient or admin:")
	fmt.Scanln(&role)

	role = ConvertToUpper(role)

	switch role {
	case "Patient":
		patientPage()
	case "Admin":
		adminPage()
	default:
		fmt.Println(invalidMessage)
	}

}

func toContinueApp() bool {

	var userChoice string

	fmt.Println("Do you still want to continue making new appointment or editing existing appointment? (Y/N)")

	if _, err := fmt.Scanln(&userChoice); err != nil {
		panic(err)
	} else {
		if strings.ToUpper(userChoice) == "Y" {
			return true
		} else if strings.ToUpper(userChoice) == "N" {
			return false
		} else {
			panic(errors.New("enter wrong choice"))
		}
	}
}
