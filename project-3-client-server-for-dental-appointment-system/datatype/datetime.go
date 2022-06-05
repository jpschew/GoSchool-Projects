package datatype

import "fmt"

const NumOfMonth = 12

const NumOfDay = 31

// added for go in action 1 to add available time slot to slice
var timeSlotAvail []string

// construct a pointer array represent Jan to Dec (index 0 to 11)
// each index will be a pointer to a dayNode which ranges from 1 to 31
type DateHashTable struct {
	month [NumOfMonth]*dayHashTable
}

// construct a pointer array represent day 1 to 31 (index 0 to 30)
// each index will be a pointer to a timeBST
type dayHashTable struct {
	day [NumOfDay]*timeBST
}

type timeBST struct {
	root *timeNode
	size int
}

type timeNode struct {
	left  *timeNode
	right *timeNode
	data  int
}

func InitDateTable() *DateHashTable {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in initialize date hashtable:", err)
		}
	}()

	// set up the array with each month node access another array with each day node
	result := &DateHashTable{}
	for i := range result.month { // 12
		result.month[i] = &dayHashTable{}
		for j := range result.month[i].day { // 31

			// at each day node, attach a BST root for the timing
			result.month[i].day[j] = &timeBST{nil, 0}

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

func (t *timeBST) addNode(r **timeNode, number int) error {

	if *r == nil {
		newNode := &timeNode{nil, nil, number}
		*r = newNode
		t.size++
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

func (t *timeBST) Add(number int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in add method in datetime:", err)
		}
	}()

	return t.addNode(&t.root, number)
}

func (t *timeBST) findSuccessor(r *timeNode) *timeNode {

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

func swap(target *timeNode, root *timeNode) {
	target.data, root.data = root.data, target.data
}
func (t *timeBST) deleteNode(r **timeNode, number int) (*timeNode, error) {

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
				t.size--
				return (*r).right, nil
			} else if (*r).right == nil {
				t.size--
				return (*r).left, nil
			} else { // both right/left child available

				// find successor to replace node to delete
				// choose right subtree to replace
				targetNode = t.findSuccessor((*r).right)
				// fmt.Println(*targetNode, "target")
				swap(targetNode, *r)

				(*r).right, err = t.deleteNode(&((*r).right), number)
				// return (*r).right

				t.size-- // only when we get the data then we deduct the size

			}

		}
	}
	// fmt.Println("before size is", t.size)
	// t.size--
	// fmt.Println("after size is", t.size)
	return *r, err
}

func (t *timeBST) Delete(number int) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in delete method in datetime:", err)
		}
	}()

	var err error
	t.root, err = t.deleteNode(&t.root, number)
	return err
}

func (t *timeBST) searchNode(r *timeNode, number int) (bool, error) {

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

func (t *timeBST) Search(number int) (bool, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic occurred in search method in datetime:", err)
		}
	}()
	return t.searchNode(t.root, number)
}

// modified for go in action 1 to return
// []string for timeslot
func (t *timeBST) inOrderTraversal(r *timeNode) []string {

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

// for printing of timeslot as it will print from smallest to largest (earliest to latest timeslot)
// modified for go in action 1 to return
// []string for timeslot
func (t *timeBST) InOrder() []string {
	// added for go in action 1 to clear time slot
	timeSlotAvail = []string{}
	return t.inOrderTraversal(t.root)
}

func (t *timeBST) postOrderTraversal(r *timeNode) {

	if r != nil {
		t.postOrderTraversal(r.left)
		t.postOrderTraversal(r.right)
		fmt.Println(r.data)
	}
}

func (t *timeBST) PostOrder() {
	t.postOrderTraversal(t.root)
}
