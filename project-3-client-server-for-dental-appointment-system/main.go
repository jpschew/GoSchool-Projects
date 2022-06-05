package main

import (
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/datatype"
	"github.com/jpschew/GoSchool-Projects/project-3-client-server-for-dental-appointment-system/server"
	"runtime"
)

// var (
// 	toContinue = true
// )

func init() {

	datatype.InitDentist()
	datatype.InitApptHashTable()
	datatype.InitAppt()
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

	server.StartServer()
	// StartServer()
	// for toContinue {

	// 	ClearScreen()
	// 	checkForAdmin()
	// 	toContinue = toContinueApp()

	// }

}

// commented out for go in action 1
// func checkForAdmin() {

// 	var role string

// 	fmt.Println("Please state if you are a patient or admin:")
// 	fmt.Scanln(&role)

// 	role = ConvertToUpper(role)

// 	switch role {
// 	case "Patient":
// 		patientPage()
// 	case "Admin":
// 		adminPage()
// 	default:
// 		fmt.Println(invalidMessage)
// 	}

// }

// func toContinueApp() bool {

// 	var userChoice string

// 	fmt.Println("Do you still want to continue making new appointment or editing existing appointment? (Y/N)")

// 	if _, err := fmt.Scanln(&userChoice); err != nil {
// 		panic(err)
// 	} else {
// 		if userChoice == "Y" {
// 			return true
// 		} else if userChoice == "N" {
// 			return false
// 		} else {
// 			panic(errors.New("enter wrong choice"))
// 		}
// 	}
// }
