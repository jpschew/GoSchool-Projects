// Package utils implements some utility functions to be used for the application.
package utils

import (
	"errors"
	"fmt"
	"github.com/jpschew/GoSchool-Projects/project-4-client-server-for-dental-appointment-system-with-security/datatype"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	traceFile = "log/trace.txt"
	warnFile  = "log/warn.txt"
	errFile   = "log/err.txt"
	traceLog  *log.Logger // log any information
	warnLog   *log.Logger // log activities that need to be concerned
	errLog    *log.Logger // log activities that are critical
)

// ConvertMonthToInt converts the month input string to an integer value.
// It takes in string input and return integer output.
// If the month is out of range, it will return an error.
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

// SendingEmail will print a message on console that the appointment details is sent to the user's email.
func SendingEmail(wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println("Sending appointment details to your email...")
	time.Sleep(3 * time.Second)
	fmt.Println("Appointment details sent to your mail!")

}

// Waiting will delay the next code execution for 3 seconds.
func Waiting() {
	fmt.Print("Waiting.")
	time.Sleep(1 * time.Second)
	fmt.Print(".")
	time.Sleep(1 * time.Second)
	fmt.Println(".")
	time.Sleep(1 * time.Second)
}

// PrintMessage will print the message to the console.
// It will take in a slice of messages and output them line by line to the console.
func PrintMessage(messageList []string) {
	for _, message := range messageList {
		fmt.Println(message)
	}
}

// IsLeap check if a particular year is a leap year.
// It will take in a year as integer and then return a true if it is a leap year, false if it is not a leap year.
func IsLeap(year int) bool {
	// leap divisible by 400 and 4 but not 100
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

// ConvertToUpper converts a string input to only the first letter is a caps while all the remaining letter will be small caps.
// It will take in a string input and return a string output.
func ConvertToUpper(input string) string {

	//// Titlt() deprecated in 1.18
	//return strings.Title(strings.ToLower(input))

	// use cases.TItle() instead
	return cases.Title(language.Und, cases.NoLower).String(strings.ToLower(input))
}

// ClearScreen is a function to clear screen.
func ClearScreen() {
	fmt.Print("\033c")
}

// ConvertIntToMonth converts the integer month input to a string.
// It takes in integer input and return a string output.
// If the month is out of range, it will return an error.
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

// CheckForValidDate checks if a given month and day is a valid date.
// This function also take into consideration of leap year.
// It will take in month and day as integer inputs and return true if it is a valid date, else false.
func CheckForValidDate(apptMonth int, apptDate int) bool {
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
		} else if apptDate >= 29 && apptMonth == 2 && !IsLeap(datatype.ApptYear) { // consider Feb without leap year
			// panic(errors.New("date is out of range"))
			return false
		}
	}
	return true
}

// TraceLogging traces all activities and write to "trace,txt" file.
func TraceLogging(message string) {

	// first args - file name
	// second args - flag for file
	// thrid args - unix file permissions (3 digit octal value)
	// typically 666 for file and 777 for directory
	// O_CREATE - create the file if file does exists
	// O_WRONLY - write only file
	// O_APPEND - append to file when writing
	file, err := os.OpenFile(traceFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed to open trace log file:", err)
	}
	traceLog = log.New(io.MultiWriter(file, os.Stderr), "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// add message to text file
	traceLog.Println(message)
}

// WarningLogging stores activities that include error or panic during the application and write to "warn,txt" file.
// However, these errors/panics are not severe to stop the application.
func WarnLogging(message string, warning error) {
	file, err := os.OpenFile(warnFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed to open warn log file:", err)
	}
	log.Println(message)
	warnLog = log.New(io.MultiWriter(file, os.Stderr), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	// add message to text file
	warnLog.Println(message, warning)

}

// ErrorLogging traces critical problem and write to a "error,txt" file.
// These errors/panics will stop the application immediately.
func ErrorLogging(message string, err error) {
	file, err := os.OpenFile(errFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	errLog = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	// add message to text file
	errLog.Println(message, err)
	// log.Fatalln(message, err)
	os.Exit(1)
}
