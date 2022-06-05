package datatype

import (
	"errors"
	"fmt"
	"sync"
)

// Size is set to 10 so each array will has 10 indexes.
const ArrSize = 10

var (
	DentistHash *DentistHashTable                         // DentishHash is a pointer that points to dentistHashTable structure that will contains all available dentists.
	NewArr      = &DentistNameList{make([]string, 10), 0} // NewArr is a pointer that points to a structure thats contains updated list of available dentists.

	dentistArr = DentistNameList{
		DentistList: []string{
			"James",
			"Michelle",
			"Ronald",
		},
		Size: 3,
	}
	idNum    = 0
	calendar *dateHashTable
)

// DentistNameList is a structure that contains the updated list of dentists.
type DentistNameList struct {
	DentistList []string // list of updated list of available dentists
	Size        int      // number of available dentists
}

// DentistInfo is a structure that contains the dentist information.
type DentistInfo struct {
	DrId      int            // dr Id
	TimeAvail *dateHashTable // pointer to a structure that will contains all the time slots from January to December.
}

// DentistList is a linkedlink structure where its node will contain dentist information.
type DentistList struct { //linkedlist
	Head *DentistNode // points to the head of dentist linkedlist
}

// DentistNode is a linkedlink node where dentist information is stored.
type DentistNode struct { //node
	Name string       // name of dentist
	Info *DentistInfo // info will contains drId and time available, refer to DentistInfo
	Next *DentistNode // points to the next node
}

// DentistHashTable is a structure that will contains all the available dentists.
// It has an array of size 10 and each index will be a pointer to a linkedlist.
type DentistHashTable struct { //hash map
	Dentist [ArrSize]*DentistList // contain a list of pointer pointing to linked-link of size 10
	Size    int                   // number of available dentists
}

// Insert inserts the dentist name to the linkedlist in the dentistHashTable.
// It takes in dentist name as string input and update the dentistHashTable.
func (t *DentistHashTable) Insert(name string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in insert method in dentist:", err)
		}
	}()

	_, index := hash(name)

	t.Dentist[index].insertDentist(name)
	t.Size++
	addToDentistList(name, t.Size)
}

func (d *DentistList) insertDentist(name string) {

	idNum++
	calendar = InitDateTable()
	newDentistInfo := &DentistInfo{idNum, calendar}
	newDentist := &DentistNode{name, newDentistInfo, nil}

	if d.Head == nil {
		d.Head = newDentist
	} else {
		newDentist.Next = d.Head
		d.Head = newDentist
	}
}

// this function add a new dentist to NewArr which is a pointer to a slice of available dentists
// and update the number of dentists in the clinic
func addToDentistList(name string, size int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add to dentist list:", err)
		}
	}()

	(*NewArr).DentistList[size-1] = name
	(*NewArr).Size++

}

// Search check if the date and timeslot of a particular dentist is available.
// It takes in timeSession, month and day as integer input, dentist name as string input and return true if the timeslot of the chosen for the dentist exists, else return false.
func (t *DentistHashTable) Search(name string, timeSession int, month int, date int) (bool, string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in search method in dentist:", err)
	// 	}
	// }()

	_, index := hash(name)
	if _, info, err := t.Dentist[index].searchDentist(name); err != nil {
		return false, "", err
	} else {
		if _, err := (*info).TimeAvail.month[month-1].day[date-1].Search(timeSession); err != nil {
			return false, "", err
		}
	}
	return true, name, nil
}

func (d *DentistList) searchDentist(name string) (bool, *DentistInfo, error) {

	if d.Head != nil {
		currentNode := d.Head
		for currentNode != nil {
			if currentNode.Name == name {
				return true, currentNode.Info, nil
			}
			currentNode = currentNode.Next
		}
	}
	return false, nil, fmt.Errorf("%s not found", name)
}

// UpdateTimeSlot deletes a time slot from a dentist available time.
// It takes in timeSession, month and day as integer input, dentist name as string input and return true if the timeslot of the chosen for the dentist exists, else return false.
// A buffered channel of 1 is also need as the input so that only one update will happen at anytime to avoid racing condition.
// The waitgroup pointer is passed in if a goroutines is used to execute this function.
func (t *DentistHashTable) UpdateTimeSlot(name string, timeSession int, month int, date int, data chan int, wg *sync.WaitGroup) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in update time slot in dentist:", err)
		}
	}()

	defer wg.Done()

	// need handle race condition issue so wont delete 2 times at once
	// here we have a buffered channel of length 1
	_, open := <-data

	if open {
		_, index := hash(name)
		_, info, err := t.Dentist[index].searchDentist(name)

		if err != nil {
			panic(err)
		}

		// need handle race condition issue so wont delete 2 times at once
		(*info).TimeAvail.month[month-1].day[date-1].Delete(timeSession)
	}

}

// AddimeSlot adds a time slot from a dentist available time.
// It takes in timeSession, month and day as integer input, dentist name as string input.
// The waitgroup pointer is passed in if a goroutines is used to execute this function.
func (t *DentistHashTable) AddTimeSlot(name string, timeSession int, month int, date int, wg *sync.WaitGroup) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add time slot in dentist:", err)
		}
	}()

	defer wg.Done()

	// _, open := <-data

	_, index := hash(name)
	_, info, err := t.Dentist[index].searchDentist(name)
	if err != nil {
		panic(errors.New("enter wrong dentist name"))
	} else {
		return (*info).TimeAvail.month[month-1].day[date-1].Add(timeSession)
	}
}

// DrdId will return the Dr Id as integer given a dentist name.
func (t *DentistHashTable) DrId(name string) int {

	_, index := hash(name)
	_, info, err := t.Dentist[index].searchDentist(name)
	if err != nil {
		panic(errors.New("enter wrong dentist name"))
	}
	return info.DrId
}

// Delete deletes a dentist name from the linkedlist in the dentistHashTable.
// It takes in dentist name as string input and update the dentistHashTable.
func (t *DentistHashTable) Delete(name string) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in delete method in dentist:", err)
		}
	}()

	_, index := hash(name)
	if t.Dentist[index] != nil {
		if err := t.Dentist[index].deleteDentist(name); err == nil {
			t.Size--
			return nil
		}
	}
	return fmt.Errorf("%s not found", name)
}

func (d *DentistList) deleteDentist(name string) error {

	if d.Head.Name == name {
		d.Head = d.Head.Next
		return nil
	}
	currentNode := d.Head
	// prevNode := d.head
	for currentNode != nil {
		prevNode := currentNode
		currentNode = currentNode.Next
		if currentNode.Name == name {
			prevNode.Next = currentNode.Next
			return nil
		}
	}
	return fmt.Errorf("%s not found", name)
}

// InitDentist initialized the available dentist to the data structure.
func InitDentist() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize dentist:", err)
		}
	}()

	// initial with 3 dentists - ["James", "Michelle", "Ronald"]
	DentistHash = &DentistHashTable{}

	for i := range DentistHash.Dentist {
		DentistHash.Dentist[i] = &DentistList{nil}
	}

	for _, name := range dentistArr.DentistList {
		DentistHash.Insert(name)
	}

}

// this function will do hashing on the input provided
// and return the sum and remainder when divided by the size of array
// the remainder is used as the index in an array which is a pointer that points to a linkedlist that stores the dentist name and its information
func hash(name string) (int, int) {

	sum := 0
	for _, ch := range name {
		// casting change the alphabet into ascii number
		// A - Z from 65 - 90
		// a - z from 97 - 122
		sum += int(ch)
	}

	return sum, sum % ArrSize // get remainder of array size
}

// ListAvailTime lists the available time slot of a dentist given a month and date.
// It takes in dentist name as string input, month and day as integer input and return a slice of string containing the available time slots.
func (t *DentistHashTable) ListAvailTime(name string, month int, date int) ([]string, error) {

	// commented out for go in action 1
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in list available time method in dentist:", err)
	// 	}
	// }()

	_, index := hash(name)
	_, info, err := t.Dentist[index].searchDentist(name)
	// fmt.Println(err)
	// fmt.Println(month-1, date-1)
	// fmt.Println(*info.timeAvail.month[4].day[11])
	if err != nil {
		// fmt.Println("dentist error")
		panic(errors.New("enter wrong dentist name"))
	}

	fmt.Println("The available time slot are:")
	timeSlot := info.TimeAvail.month[month-1].day[date-1].InOrder()

	return timeSlot, nil
}

// List lists the update list of dentists available in the clinics.
// It will also return a slice of string containing the updated list of available dentists.
func (d *DentistNameList) List() []string {

	return d.DentistList
}
