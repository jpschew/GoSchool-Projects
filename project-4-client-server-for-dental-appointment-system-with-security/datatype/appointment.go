package datatype

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const ApptArrSize = 10

var (
	ApptHash   *ApptHashTable      // hashtable for storing appointment using appointment id
	DrApptHash *ApptHashTable      // hashtable for storing appointment using dr Id
	ApptYear   = time.Now().Year() // current year
	TimeArr    = [15]string{
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
	} // appointment timing
	apptNum = 0
)

// ApptHashTable is a structure that will contains all the appointment list.
// It has an array of size 10 and each index will be a pointer to a linkedlist.
type ApptHashTable struct {
	Arr  [ApptArrSize]*AppointmentList // contain a list of pointer pointing to linked-link of size 10
	Size int                           // number of available dentists
}

// AppointmentList is a linkedlink structure where its node will contain appointment information.
type AppointmentList struct {
	Head *AppointmentNode // points to the head of appointment linkedlist
}

// ApoointmentNode is a linkedlink node where appointment information is stored.
type AppointmentNode struct {
	ApptId  int              // name of dentist
	Details *ApptDetails     // details will contains appointment details, refer to ApptDetails
	Next    *AppointmentNode // points to the next node
}

// ApptDetails is a structure that contains the appointment information.
type ApptDetails struct {
	Name        string // name of dentist
	DrId        int    // dr Id
	ApptMonth   int    // month of appointment
	ApptDate    int    // day of appointment
	ApptTime    int    // time of appointment
	PatientUser string // patient username
}

// DentistName will return the name of dentist as string.
func (d *ApptDetails) DentistName() string {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.Name
}

// Patient will return the username of patient as string.
func (d *ApptDetails) Patient() string {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.PatientUser
}

// Month will return the appointment month as integer.
func (d *ApptDetails) Month() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.ApptMonth
}

// Day will return the appointment day as integer.
func (d *ApptDetails) Day() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.ApptDate
}

// Time will return the appointment time session as integer.
func (d *ApptDetails) Time() int {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in check changes method in appointment:", err)
	// 	}
	// }()

	return d.ApptTime
}

// CheckChanges will check if appointment information has been change.
// It will takes in time session, month, day as integer input, dentist name as string input and return true if there is any changes to the appointment, else return false.
func (d *ApptDetails) CheckChanges(timeSession int, month int, date int, name string) (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in check changes method in appointment:", err)
		}
	}()

	if d.Name == name && d.ApptTime == timeSession && d.ApptMonth == month && d.ApptDate == date {
		return false, errors.New("nothing has been changed to the appointment")
	}
	return true, nil
}

// this function will update the appointment list if there is any changes
func (l *AppointmentList) updateAppt(name string, drId int, timeSession int, month int, date int, apptId int, userName string) error {

	if l.Head == nil {
		return fmt.Errorf("appointment id %d not found", apptId)
	} else {
		currentNode := l.Head
		for currentNode != nil {
			if currentNode.ApptId == apptId {
				currentNode.Details.Name = name
				currentNode.Details.DrId = drId
				currentNode.Details.ApptMonth = month
				currentNode.Details.ApptDate = date
				currentNode.Details.ApptTime = timeSession
				currentNode.Details.PatientUser = userName
				return nil
			}
			currentNode = currentNode.Next
		}
	}
	return fmt.Errorf("appointment id %d not found", apptId)
}

// Update will update the appointment details if same dentist is chosen.
// It will take in dr Id, time session, month, dat and appointment Id as integer input, dentist name and patient's username as string input.
// If isAppt is true means to update appointment using appointment id else means update appointment using drId.
// Pointer to waitgroup will need to be passed in if a goroutines is executed on this function.
func (a *ApptHashTable) Update(name string, drId int, timeSession int, month int, date int, apptId int, isAppt bool, wg *sync.WaitGroup, userName string) {

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

	a.Arr[index].updateAppt(name, drId, timeSession, month, date, apptId, userName)
	// a.size++
	// return apptNum
}

// this function remove san appointment from the appointment list
func (l *AppointmentList) removeAppt(name string, drId int, timeSession int, month int, date int, apptId int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in remove appointment method in appointment:", err)
		}
	}()

	if l.Head == nil {
		return fmt.Errorf("no appointment for Dr. %s to remove", name)
	} else {
		currentNode := l.Head
		if currentNode.ApptId == apptId {
			l.Head = currentNode.Next
			return nil
		}
		// prevNode := l.head
		for currentNode != nil {
			prevNode := currentNode
			currentNode = currentNode.Next
			if currentNode.ApptId == apptId {
				prevNode.Next = currentNode.Next
				return nil
			}

		}
	}
	return fmt.Errorf("appointment id %d not found", apptId)
}

// UpdateDiffDr will update the appointment details if different dentist is chosen.
// It will take in dr Id, time session, month, dat and appointment Id as integer input, dentist name and patient's username as string input.
// If isAppt is true means to update appointment using appointment id else means update appointment using drId.
// Pointer to waitgroup will need to be passed in if a goroutines is executed on this function.
func (a *ApptHashTable) UpdateDiffDr(name string, newDrId int, oldDrId int, timeSession int, month int, date int, apptId int, isAppt bool, wg *sync.WaitGroup, userName string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in update different dentist method in appointment:", err)
		}
	}()

	defer wg.Done()

	var index int

	index = hashFunction(oldDrId)

	a.Arr[index].removeAppt(name, oldDrId, timeSession, month, date, apptId)

	index = hashFunction(newDrId)
	a.Arr[index].addAppt(name, newDrId, timeSession, month, date, apptId, userName)

}

// this function adds an appointment to the appintment list
func (l *AppointmentList) addAppt(name string, drId int, timeSession int, month int, date int, apptId int, userName string) {

	details := &ApptDetails{name, drId, month, date, timeSession, userName}
	newAppt := &AppointmentNode{apptId, details, nil}
	if l.Head == nil {
		l.Head = newAppt
	} else {
		newAppt.Next = l.Head
		l.Head = newAppt
	}

}

// Add adds an appoitnment to the appointment linkedlist.
// It will take in dr Id, time session, month, dat and appointment Id as integer input, dentist name and patient's username as string input.
// If isAppt is true means to update appointment using appointment id else means update appointment using drId.
// Pointer to waitgroup will need to be passed in if a goroutines is executed on this function.
func (a *ApptHashTable) Add(name string, drId int, timeSession int, month int, date int, apptIdChn chan int, isAppt bool, wg *sync.WaitGroup, userName string) {

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

	a.Arr[index].addAppt(name, drId, timeSession, month, date, apptNum, userName)
	a.Size++

}

func (l *AppointmentList) searchAppt(apptId int) (bool, *ApptDetails, string, error) {

	// var output string

	if l.Head == nil {
		return false, nil, "", fmt.Errorf("appointment id %d not found", apptId)
	} else {
		currentNode := l.Head
		for currentNode != nil {
			if currentNode.ApptId == apptId {
				output := fmt.Sprintf("Appointment Id %d is on %d-%02d-%02d at %s with Dr. %s for %s.\n",
					apptId, ApptYear, currentNode.Details.ApptMonth, currentNode.Details.ApptDate, TimeArr[currentNode.Details.ApptTime-1], currentNode.Details.Name, currentNode.Details.PatientUser)
				fmt.Printf("Appointment Id %d is on %d-%02d-%02d at %s with Dr. %s for %s.\n",
					apptId, ApptYear, currentNode.Details.ApptMonth, currentNode.Details.ApptDate, TimeArr[currentNode.Details.ApptTime-1], currentNode.Details.Name, currentNode.Details.PatientUser)
				return true, currentNode.Details, output, nil
			}
			currentNode = currentNode.Next
		}
	}
	return false, nil, "", fmt.Errorf("appointment id %d not found", apptId)

}

// Search check if an appointment is available using the appointment Id
// It takes in appointment id as integer input and return true if the appointment Id exists, else return false.
func (a *ApptHashTable) Search(apptId int) (bool, *ApptDetails, string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search method in appointment:", err)
		}
	}()

	index := hashFunction(apptId)

	if _, details, output, err := a.Arr[index].searchAppt(apptId); err != nil {
		// fmt.Println("cant find name")
		return false, nil, output, err
	} else {
		return true, details, output, nil
	}
}

func (l *AppointmentList) browseDrAppt() ([]string, error) {

	var output []string

	if l.Head == nil {
		return []string{}, errors.New("no appointment has been made for this dentist")
	} else {
		currentNode := l.Head
		for currentNode != nil {
			s := fmt.Sprintf("Appointment Id %d is on %d-%02d-%02d %s for %s.\n",
				currentNode.ApptId, ApptYear, currentNode.Details.ApptMonth, currentNode.Details.ApptDate, TimeArr[currentNode.Details.ApptTime-1], currentNode.Details.PatientUser)
			output = append(output, s)
			fmt.Printf("Appointment Id %d is on %d-%02d-%02d %s for %s.\n",
				currentNode.ApptId, ApptYear, currentNode.Details.ApptMonth, currentNode.Details.ApptDate, TimeArr[currentNode.Details.ApptTime-1], currentNode.Details.PatientUser)
			currentNode = currentNode.Next
		}
	}
	return output, nil

}

// Browse wlll browse all the appointments by a particular dentist.
// It takes in dr id as integer input and return a list of appointments by the chosen dentist.
func (a *ApptHashTable) Browse(drId int) ([]string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in browse method in appointment:", err)
		}
	}()

	index := hashFunction(drId)
	return a.Arr[index].browseDrAppt()
}

// this function will do hashing on the input provided
// and return remainder when divided by the size of appointment array
// the remainder is used as the index in an array which is a pointer that points to a linkedlist that stores the details of appointment
func hashFunction(id int) int {
	return id % ApptArrSize
}

// InitApptHashTable initialized 2 hash table i) appointment hash table and  ii) dentist appointment table.
// Appoinment hash table is hashed using the appointment id.
// Dentist Appoinment hash table is hashed using the dr id.
func InitApptHashTable() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in intialize appointment hashtable:", err)
		}
	}()

	ApptHash = &ApptHashTable{}
	for i := 0; i < ApptArrSize; i++ {
		ApptHash.Arr[i] = &AppointmentList{}
	}
	DrApptHash = &ApptHashTable{}
	for i := 0; i < ApptArrSize; i++ {
		DrApptHash.Arr[i] = &AppointmentList{}
	}

}

// InitAppt initialized some appointments to the appoinment linkedlist.
// These appointments are stored in both the Appointment Hash Table as well as Dentist Appointment Hash Table.
func InitAppt() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize appointment:", err)
		}
	}()

	// var monthArr = []int{5, 6, 7}
	var dayArr = []int{11, 12, 13}
	var timeArr = []int{1, 2, 3}
	var userName = []string{"jpschew", "ken", "ivan"}

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
				DentistHash.UpdateTimeSlot(dentistArr.DentistList[i], session, 5, day, chn, &wg)
				// close(chn)
				ApptHash.Add(dentistArr.DentistList[i], DentistHash.DrId(dentistArr.DentistList[i]), session, 5, day, chn, true, &wg, userName)
				DrApptHash.Add(dentistArr.DentistList[i], DentistHash.DrId(dentistArr.DentistList[i]), session, 5, day, chn, false, &wg, userName)
				close(chn)
			}
		}
	}

}
