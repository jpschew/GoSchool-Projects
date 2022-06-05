package utils

import (
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"sync"
	"time"
)

func ConvertMonthToInt(month string) (int, error) {

	switch month {
	case "January":
		return 1, nil
	case "Februray":
		return 2, nil
	case "March":
		return 3, nil
	case "April":
		return 4, nil
	case "May":
		return 5, nil
	case "June":
		return 6, nil
	case "July":
		return 7, nil
	case "August":
		return 8, nil
	case "September":
		return 9, nil
	case "October":
		return 10, nil
	case "November":
		return 11, nil
	case "December":
		return 12, nil
	default:
		return -1, errors.New("please spell out the month correctly")
	}
}

func SendingEmail(wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("Sending appointment details to your email...")
	time.Sleep(3 * time.Second)
	fmt.Println("Appointment details sent to your mail!")

}

func Waiting() {
	fmt.Print("Waiting.")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Println(".")
	time.Sleep(1 * time.Second)
}

func PrintMessage(messageList []string) {
	for _, message := range messageList {
		fmt.Println(message)
	}
}

// leap divisible by 400 and 4 but not 100
func IsLeap(year int) bool {
	if year%400 == 0 {
		return true
	} else {
		if year%100 == 0 {
			return false
		} else if year%4 == 0 {
			return true
		}
	}
	return false

}

func ConvertToUpper(input string) string {

	//// Titlt() deprecated in 1.18
	//return strings.Title(strings.ToLower(input))

	// use cases.TItle() instead
	return cases.Title(language.Und, cases.NoLower).String(strings.ToLower(input))
}

// function to clear screen
func ClearScreen() {
	fmt.Print("\033c")
}

// added for go in action 1
func ConvertIntToMonth(month int) (string, error) {

	switch month {
	case 1:
		return "January", nil
	case 2:
		return "Februray", nil
	case 3:
		return "March", nil
	case 4:
		return "April", nil
	case 5:
		return "May", nil
	case 6:
		return "June", nil
	case 7:
		return "July", nil
	case 8:
		return "August", nil
	case 9:
		return "September", nil
	case 10:
		return "October", nil
	case 11:
		return "November", nil
	case 12:
		return "December", nil
	default:
		return "", errors.New("please spell out the month correctly")
	}
}
