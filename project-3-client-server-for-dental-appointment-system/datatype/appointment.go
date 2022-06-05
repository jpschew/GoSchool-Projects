package datatype

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const apptArrSize = 10

var apptNum = 0
var ApptHash *apptHashTable   // hashtable for storing appointment using appointment id
var DrApptHash *apptHashTable // hashtable for storing appointment using dr Id

// need to change to env file
// brought over from server.g
var userName = []string{"jpschew", "ken", "ivan"}

// brought over from patient.go
var ApptYear = time.Now().Year()
var TimeArr = [15]string{
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

type apptHashTable struct {
	arr  [apptArrSize]*appointmentList
	size int
}

type appointmentList struct {
	head *appointmentNode
}

type appointmentNode struct {
	apptId  int
	details *ApptDetails
	next    *appointmentNode
}

type ApptDetails struct {
	// pId int
	name        string
	drId        int
	apptMonth   int
	apptDate    int
	apptTime    int
	patientUser string // added for go in action 1
}

func (d *ApptDetails) Name() string {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.name
}

func (d *ApptDetails) Patient() string {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.patientUser
}

func (d *ApptDetails) Month() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.apptMonth
}

func (d *ApptDetails) Day() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.apptDate
}

func (d *ApptDetails) Time() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.apptTime
}

func (d *ApptDetails) CheckChanges(timeSession int, month int, date int, name string) (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in check changes method in appointment:", err)
		}
	}()

	if d.name == name && d.apptTime == timeSession && d.apptMonth == month && d.apptDate == date {
		return false, errors.New("nothing has been changed to the appointment")
	}
	return true, nil
}

func (l *appointmentList) updateAppt(name string, drId int, timeSession int, month int, date int, apptId int, userName string) error {

	if l.head == nil {
		return fmt.Errorf("appointment id %d not found", apptId)
	} else {
		currentNode := l.head
		for currentNode != nil {
			if currentNode.apptId == apptId {
				currentNode.details.name = name
				currentNode.details.drId = drId
				currentNode.details.apptMonth = month
				currentNode.details.apptDate = date
				currentNode.details.apptTime = timeSession
				currentNode.details.patientUser = userName
				return nil
			}
			currentNode = currentNode.next
		}
	}
	return fmt.Errorf("appointment id %d not found", apptId)
}

// if isAppt is true means to update appointment using appointment id else means update appointment using drId
// this function is for updating appointment if same dentist is chosen
func (a *apptHashTable) Update(name string, drId int, timeSession int, month int, date int, apptId int, isAppt bool, wg *sync.WaitGroup, userName string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in update method in appointment:", err)
		}
	}()

	defer wg.Done()

	var index int

	if isAppt {
		// apptNum++
		index = hashFunction(apptId)
	} else {
		index = hashFunction(drId)
	}

	a.arr[index].updateAppt(name, drId, timeSession, month, date, apptId, userName)
	// a.size++
	// return apptNum
}

func (l *appointmentList) removeAppt(name string, drId int, timeSession int, month int, date int, apptId int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in remove appointment method in appointment:", err)
		}
	}()

	if l.head == nil {
		return fmt.Errorf("no appointment for Dr. %s to remove", name)
	} else {
		currentNode := l.head
		if currentNode.apptId == apptId {
			l.head = currentNode.next
			return nil
		}
		// prevNode := l.head
		for currentNode != nil {
			prevNode := currentNode
			currentNode = currentNode.next
			if currentNode.apptId == apptId {
				prevNode.next = currentNode.next
				return nil
			}

		}
	}
	return fmt.Errorf("appointment id %d not found", apptId)
}

// if isAppt is true means to update appointment using appointment id else means update appointment using dentist name
// this function is for updating appointment if a different dentist is chosen
func (a *apptHashTable) UpdateDiffDr(name string, newDrId int, oldDrId int, timeSession int, month int, date int, apptId int, isAppt bool, wg *sync.WaitGroup, userName string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in update different dentist method in appointment:", err)
		}
	}()

	defer wg.Done()

	var index int

	index = hashFunction(oldDrId)

	a.arr[index].removeAppt(name, oldDrId, timeSession, month, date, apptId)

	index = hashFunction(newDrId)
	a.arr[index].addAppt(name, newDrId, timeSession, month, date, apptId, userName)

}

func (l *appointmentList) addAppt(name string, drId int, timeSession int, month int, date int, apptId int, userName string) {

	details := &ApptDetails{name, drId, month, date, timeSession, userName}
	newAppt := &appointmentNode{apptId, details, nil}
	if l.head == nil {
		l.head = newAppt
	} else {
		newAppt.next = l.head
		l.head = newAppt
	}

}

// modified for go in action 1
// to add in username
func (a *apptHashTable) Add(name string, drId int, timeSession int, month int, date int, apptIdChn chan int, isAppt bool, wg *sync.WaitGroup, userName string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add method in appointment:", err)
		}
	}()

	defer wg.Done()

	var index int

	if isAppt {
		apptNum++
		index = hashFunction(apptNum)
		apptIdChn <- apptNum // send back the appointment id to the goroutines
	} else {
		index = hashFunction(drId)
	}

	a.arr[index].addAppt(name, drId, timeSession, month, date, apptNum, userName)
	a.size++

}

// modify for go in action 1
func (l *appointmentList) searchAppt(apptId int) (bool, *ApptDetails, string, error) {

	// var output string

	if l.head == nil {
		return false, nil, "", fmt.Errorf("appointment id %d not found", apptId)
	} else {
		currentNode := l.head
		for currentNode != nil {
			if currentNode.apptId == apptId {
				output := fmt.Sprintf("Appointment Id %d is on %d-%02d-%02d at %s with Dr. %s for %s.\n",
					apptId, ApptYear, currentNode.details.apptMonth, currentNode.details.apptDate, TimeArr[currentNode.details.apptTime-1], currentNode.details.name, currentNode.details.patientUser)
				fmt.Printf("Appointment Id %d is on %d-%02d-%02d at %s with Dr. %s for %s.\n",
					apptId, ApptYear, currentNode.details.apptMonth, currentNode.details.apptDate, TimeArr[currentNode.details.apptTime-1], currentNode.details.name, currentNode.details.patientUser)
				return true, currentNode.details, output, nil
			}
			currentNode = currentNode.next
		}
	}
	return false, nil, "", fmt.Errorf("appointment id %d not found", apptId)

}

// modify for go in action 1
// added userName and return additonal string
func (a *apptHashTable) Search(apptId int) (bool, *ApptDetails, string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search method in appointment:", err)
		}
	}()

	index := hashFunction(apptId)

	if _, details, output, err := a.arr[index].searchAppt(apptId); err != nil {
		// fmt.Println("cant find name")
		return false, nil, output, err
	} else {
		return true, details, output, nil
	}
}

// added one more return type []string
// to return appointment list
func (l *appointmentList) browseDrAppt() ([]string, error) {

	var output []string

	if l.head == nil {
		return []string{}, errors.New("no appointment has been made for this dentist")
	} else {
		currentNode := l.head
		for currentNode != nil {
			s := fmt.Sprintf("Appointment Id %d is on %d-%02d-%02d %s for %s.\n",
				currentNode.apptId, ApptYear, currentNode.details.apptMonth, currentNode.details.apptDate, TimeArr[currentNode.details.apptTime-1], currentNode.details.patientUser)
			output = append(output, s)
			fmt.Printf("Appointment Id %d is on %d-%02d-%02d %s for %s.\n",
				currentNode.apptId, ApptYear, currentNode.details.apptMonth, currentNode.details.apptDate, TimeArr[currentNode.details.apptTime-1], currentNode.details.patientUser)
			currentNode = currentNode.next
		}
	}
	return output, nil

}

// modify for go in action 1
func (a *apptHashTable) Browse(drId int) ([]string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in browse method in appointment:", err)
		}
	}()

	index := hashFunction(drId)
	return a.arr[index].browseDrAppt()
}

func hashFunction(id int) int {
	return id % apptArrSize
}

// this function is for initializing 2 appointment hashtable
// 1 for using appointment id
// anohter for using dr id
func InitApptHashTable() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in intialize appointment hashtable:", err)
		}
	}()

	ApptHash = &apptHashTable{}
	for i := 0; i < apptArrSize; i++ {
		ApptHash.arr[i] = &appointmentList{}
	}
	DrApptHash = &apptHashTable{}
	for i := 0; i < apptArrSize; i++ {
		DrApptHash.arr[i] = &appointmentList{}
	}

}

// this function is for initializing dentist appointment
func InitAppt() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize appointment:", err)
		}
	}()

	// var monthArr = []int{5, 6, 7}
	var dayArr = []int{11, 12, 13}
	var timeArr = []int{1, 2, 3}

	var wg = sync.WaitGroup{}

	// used in go advanced
	// for _, dentist := range DentistArr.list {
	// 	// for _, month := range monthArr {
	// 	for _, day := range dayArr {
	// 		for _, session := range timeArr {
	// 			var chn = make(chan int, 1)
	// 			wg.Add(3)
	// 			chn <- 1
	// 			DentistHash.UpdateTimeSlot(dentist, session, 5, day, chn, &wg)
	// 			// close(chn)
	// 			ApptHash.Add(dentist, DentistHash.getDrId(dentist), session, 5, day, chn, true, &wg)
	// 			DrApptHash.Add(dentist, DentistHash.getDrId(dentist), session, 5, day, chn, false, &wg)
	// 			close(chn)

	// 		}
	// 	}
	// }

	// used in go in action 1
	for i, userName := range userName {
		for _, day := range dayArr {
			for _, session := range timeArr {
				var chn = make(chan int, 1)
				wg.Add(3)
				chn <- 1
				DentistHash.UpdateTimeSlot(DentistArr.list[i], session, 5, day, chn, &wg)
				// close(chn)
				ApptHash.Add(DentistArr.list[i], DentistHash.GetDrId(DentistArr.list[i]), session, 5, day, chn, true, &wg, userName)
				DrApptHash.Add(DentistArr.list[i], DentistHash.GetDrId(DentistArr.list[i]), session, 5, day, chn, false, &wg, userName)
				close(chn)
			}
		}
	}

}
