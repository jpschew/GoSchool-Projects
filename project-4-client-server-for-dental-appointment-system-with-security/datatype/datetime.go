// Package dataype includes the data structure to be used for the implementing the appointment list, dentist list as well as the available time slot tagged to the dentist.
package datatype

import "fmt"

// declared the number of month in a year and number of days in a month as constant
const (
	numOfMonth = 12
	numOfDay   = 31
)

// added for go in action 1 to add available time slot to slice
var timeSlotAvail []string

// dateHashTable is a structure that will contains all the time slots from January to December.
// it will construct a pointer array represent Jan to Dec (index 0 to 11).
// each index will be a pointer to a dayNode which ranges from 1 to 31 (index 0 to 30).
type dateHashTable struct {
	month [numOfMonth]*dayHashTable
}

// construct a pointer array represent day 1 to 31 (index 0 to 30)
// each index will be a pointer to a timeBST
type dayHashTable struct {
	day [numOfDay]*TimeBST
}

// TimeBST is a BST structure that stores the time session for an appointment.
// Time slot starts from 9AM to 5PM represented by 1 to 15 (each slot has 30minutes duration).
type TimeBST struct {
	Root *timeNode // pointer to the root node
	Size int       // number of time slots available (max. 15)
}

// time node of BST
type timeNode struct {
	left  *timeNode
	right *timeNode
	data  int
}

// InitDateTable initialized all the available time slots (session 1 to 15) to all the dates from January to December.
func InitDateTable() *dateHashTable {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize date hashtable:", err)
		}
	}()

	// set up the array with each month node access another array with each day node
	result := &dateHashTable{}
	for i := range result.month { // 12
		result.month[i] = &dayHashTable{}
		for j := range result.month[i].day { // 31

			// at each day node, attach a BST root for the timing
			result.month[i].day[j] = &TimeBST{nil, 0}

			// add timing slot to BST in such a way the middle is added first
			// so that we can have a balanced tree
			result.month[i].day[j].Add(8)
			result.month[i].day[j].Add(4)
			result.month[i].day[j].Add(12)
			result.month[i].day[j].Add(2)
			result.month[i].day[j].Add(14)
			result.month[i].day[j].Add(6)
			result.month[i].day[j].Add(10)
			result.month[i].day[j].Add(1)
			result.month[i].day[j].Add(3)
			result.month[i].day[j].Add(5)
			result.month[i].day[j].Add(7)
			result.month[i].day[j].Add(9)
			result.month[i].day[j].Add(11)
			result.month[i].day[j].Add(13)
			result.month[i].day[j].Add(15)

		}

	}

	return result
}

func (t *TimeBST) addNode(r **timeNode, number int) error {

	if *r == nil {
		newNode := &timeNode{nil, nil, number}
		*r = newNode
		t.Size++
		return nil
	} else {
		if (*r).data < number {
			t.addNode(&((*r).right), number)
		} else if (*r).data > number {
			t.addNode(&((*r).left), number)
		} else {
			return fmt.Errorf("cannot add same timeslot on the same day")
		}
	}
	// t.size++
	return nil
}

// Add adds the time node to the BST.
// It takes in integer (session slot) as input and return error if there is any.
func (t *TimeBST) Add(number int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add method in datetime:", err)
		}
	}()

	return t.addNode(&t.Root, number)
}

func (t *TimeBST) findSuccessor(r *timeNode) *timeNode {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in find successor method in datetime:", err)
		}
	}()

	// choose right subtree to replace
	// get the left most of right subtree
	for r.left != nil {
		r = r.left
	}
	// fmt.Println(r, "target root", r.left)
	return r

}

// this function swaps the value of the node to be removed and the value of the successor node
func swap(target *timeNode, root *timeNode) {
	target.data, root.data = root.data, target.data
}

func (t *TimeBST) deleteNode(r **timeNode, number int) (*timeNode, error) {

	var targetNode *timeNode
	var err error

	if *r == nil {
		// fmt.Println("no nodes in tree")
		return nil, fmt.Errorf("the item is not in the tree")
	} else {
		if (*r).data < number {
			(*r).right, err = t.deleteNode(&((*r).right), number)
			// return (*r).right
		} else if (*r).data > number {
			(*r).left, err = t.deleteNode(&((*r).left), number)
			// return (*r).left
		} else { // number at root node

			// only right/left node
			if (*r).left == nil {
				t.Size--
				return (*r).right, nil
			} else if (*r).right == nil {
				t.Size--
				return (*r).left, nil
			} else { // both right/left child available

				// find successor to replace node to delete
				// choose right subtree to replace
				targetNode = t.findSuccessor((*r).right)
				// fmt.Println(*targetNode, "target")
				swap(targetNode, *r)

				(*r).right, err = t.deleteNode(&((*r).right), number)
				// return (*r).right

				t.Size-- // only when we get the data then we deduct the size

			}

		}
	}
	// fmt.Println("before size is", t.size)
	// t.size--
	// fmt.Println("after size is", t.size)
	return *r, err
}

// Delete removes the time node from the BST.
// It takes in integer (session slot) as input and return error if there is any.
func (t *TimeBST) Delete(number int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in delete method in datetime:", err)
		}
	}()

	var err error
	t.Root, err = t.deleteNode(&t.Root, number)
	return err
}

func (t *TimeBST) searchNode(r *timeNode, number int) (bool, error) {

	if r == nil {
		return false, fmt.Errorf("the item is not in the tree")
	} else {
		if r.data == number {
			return true, nil
		} else {
			if r.data > number {
				return t.searchNode(r.left, number)
			} else {
				return t.searchNode(r.right, number)
			}
		}
	}

}

// Search check if the time node exists in the BST.
// It takes in integer (session slot) as input and return true if it exists, else return false.
func (t *TimeBST) Search(number int) (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search method in datetime:", err)
		}
	}()
	return t.searchNode(t.Root, number)
}

// modified for go in action 1 to return
// []string for timeslot
func (t *TimeBST) inOrderTraversal(r *timeNode) []string {

	// added for go into action 1
	// var timeSlot []string

	if r != nil {
		t.inOrderTraversal(r.left)
		fmt.Println(TimeArr[r.data-1])
		timeSlotAvail = append(timeSlotAvail, TimeArr[r.data-1])
		t.inOrderTraversal(r.right)
	}

	return timeSlotAvail
}

// InOrder will print the time slot from smallest to largest (earliest to latest timeslot).
// This function will also return a slice of time slot from smallest to largest (earliest to latest timeslot).
// for printing of timeslot as it will print from smallest to largest (earliest to latest timeslot)
// modified for go in action 1 to return
// []string for timeslot
func (t *TimeBST) InOrder() []string {
	// added for go in action 1 to clear time slot
	timeSlotAvail = []string{}
	return t.inOrderTraversal(t.Root)
}

// func (t *timeBST) postOrderTraversal(r *timeNode) {

// 	if r != nil {
// 		t.postOrderTraversal(r.left)
// 		t.postOrderTraversal(r.right)
// 		fmt.Println(r.data)
// 	}
// }

// func (t *timeBST) PostOrder() {
// 	t.postOrderTraversal(t.root)
// }
