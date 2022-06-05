package datatype

import (
	"errors"
	"fmt"
	"sync"
)

const size = 10

var idNum = 0 // need handle concuurency issues when updating appt num
var DentistHash *dentistHashTable

var calendar *DateHashTable

var NewArr = &dentistNameList{make([]string, 10), 0}

var DentistArr = dentistNameList{
	list: []string{
		"James",
		"Michelle",
		"Ronald",
	},
	size: 3,
}

type dentistNameList struct {
	list []string
	size int
}

type dentistInfo struct {
	drId int
	// name string
	timeAvail *DateHashTable
}

type dentistList struct { //linkedlist
	head *dentistNode
}

type dentistNode struct { //node
	name string       // name of dentist will be the key
	info *dentistInfo // info will contains drId and time available
	next *dentistNode
}

type dentistHashTable struct { //hash map
	dentist [size]*dentistList // will contain a list of pointer pointing to linked-link
	size    int
}

func (t *dentistHashTable) Insert(name string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in insert method in dentist:", err)
		}
	}()

	_, index := hash(name)

	t.dentist[index].insertDentist(name)
	t.size++
	addToDentistList(name, t.size)
}

func (d *dentistList) insertDentist(name string) {

	idNum++
	calendar = InitDateTable()
	newDentistInfo := &dentistInfo{idNum, calendar}
	newDentist := &dentistNode{name, newDentistInfo, nil}

	if d.head == nil {
		d.head = newDentist
	} else {
		newDentist.next = d.head
		d.head = newDentist
	}
}

func addToDentistList(name string, size int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add to dentist list:", err)
		}
	}()

	(*NewArr).list[size-1] = name
	(*NewArr).size++

}

func (t *dentistHashTable) Search(name string, timeSession int, month int, date int) (bool, string, error) {

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in search method in dentist:", err)
	// 	}
	// }()

	_, index := hash(name)
	if _, info, err := t.dentist[index].searchDentist(name); err != nil {
		return false, "", err
	} else {
		if _, err := (*info).timeAvail.month[month-1].day[date-1].Search(timeSession); err != nil {
			return false, "", err
		}
	}
	return true, name, nil
}

func (d *dentistList) searchDentist(name string) (bool, *dentistInfo, error) {

	if d.head != nil {
		currentNode := d.head
		for currentNode != nil {
			if currentNode.name == name {
				return true, currentNode.info, nil
			}
			currentNode = currentNode.next
		}
	}
	return false, nil, fmt.Errorf("%s not found", name)
}

func (t *dentistHashTable) UpdateTimeSlot(name string, timeSession int, month int, date int, data chan int, wg *sync.WaitGroup) {

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
		_, info, err := t.dentist[index].searchDentist(name)

		if err != nil {
			panic(err)
		}

		// need handle race condition issue so wont delete 2 times at once
		(*info).timeAvail.month[month-1].day[date-1].Delete(timeSession)
	}

}

func (t *dentistHashTable) AddTimeSlot(name string, timeSession int, month int, date int, wg *sync.WaitGroup) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add time slot in dentist:", err)
		}
	}()

	defer wg.Done()

	// _, open := <-data

	_, index := hash(name)
	_, info, err := t.dentist[index].searchDentist(name)
	if err != nil {
		panic(errors.New("enter wrong dentist name"))
	} else {
		return (*info).timeAvail.month[month-1].day[date-1].Add(timeSession)
	}
}

func (t *dentistHashTable) GetDrId(name string) int {

	_, index := hash(name)
	_, info, err := t.dentist[index].searchDentist(name)
	if err != nil {
		panic(errors.New("enter wrong dentist name"))
	}
	return info.drId
}

func (t *dentistHashTable) Delete(name string) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in delete method in dentist:", err)
		}
	}()

	_, index := hash(name)
	if t.dentist[index] != nil {
		if err := t.dentist[index].deleteDentist(name); err == nil {
			t.size--
			return nil
		}
	}
	return fmt.Errorf("%s not found", name)
}

func (d *dentistList) deleteDentist(name string) error {

	if d.head.name == name {
		d.head = d.head.next
		return nil
	}
	currentNode := d.head
	// prevNode := d.head
	for currentNode != nil {
		prevNode := currentNode
		currentNode = currentNode.next
		if currentNode.name == name {
			prevNode.next = currentNode.next
			return nil
		}
	}
	return fmt.Errorf("%s not found", name)
}

// this function is for initializing all the 15 timeslots all the month and day
func InitDentist() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize dentist:", err)
		}
	}()

	// initial with 3 dentists - ["James", "Michelle", "Ronald"]
	DentistHash = &dentistHashTable{}

	for i := range DentistHash.dentist {
		DentistHash.dentist[i] = &dentistList{nil}
	}

	for _, name := range DentistArr.list {
		DentistHash.Insert(name)
	}

}

func hash(name string) (int, int) {

	sum := 0
	for _, ch := range name {
		// casting change the alphabet into ascii number
		// A - Z from 65 - 90
		// a - z from 97 - 122
		sum += int(ch)
	}

	return sum, sum % size // get remainder of array size
}

// modified for go in action 1 to return
// []string of timeslot
func (t *dentistHashTable) ListAvailTime(name string, month int, date int) ([]string, error) {

	// commented out for go in action 1
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("Panic occurred in list available time method in dentist:", err)
	// 	}
	// }()

	_, index := hash(name)
	_, info, err := t.dentist[index].searchDentist(name)
	// fmt.Println(err)
	// fmt.Println(month-1, date-1)
	// fmt.Println(*info.timeAvail.month[4].day[11])
	if err != nil {
		// fmt.Println("dentist error")
		panic(errors.New("enter wrong dentist name"))
	}

	fmt.Println("The available time slot are:")
	timeSlot := info.timeAvail.month[month-1].day[date-1].InOrder()

	return timeSlot, nil
}

func (d *dentistNameList) List() []string {

	return d.list
}
