package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
	"time"
)

// global variables declared
var userChoice int

// use pointer for item struct, so we can set the value by referencing
var itemMap map[string]*itemInfo
var category []string
var list []shoppingList
var (
	// returnMainMenu = "You will be returning to Main Menu."
	emptyMessage = "There is no item in the shopping list."
	delay        = time.Second * 1
	listCount    = 0
)

func init() {
	// first args = true to initialize shopping list according to assignment
	// first args = false to initialize empty shopping list
	// second args = true to initialize shopping list for testing modifying/deleting items
	// cannot be both true
	initializeShoppingList(true, false)
}
func main() {

	for {
		mainMenu()

		switch userChoice {
		case 1:
			clearScreen()
			// args = true to show items sorted by category
			// args = false show shopping list contents
			displayShoppingList(itemMap, category, false)
			// args = false for returning to main menu
			// args = true for exiting
			returnToMainMenu(false)
		case 2:
			clearScreen()
			generateReport()
			returnToMainMenu(false)
		case 3:
			clearScreen()
			addItem()
			returnToMainMenu(false)
		case 4:
			clearScreen()
			modifyItem()
			returnToMainMenu(false)
		case 5:
			clearScreen()
			deleteItem()
			returnToMainMenu(false)
		case 6:
			clearScreen()
			printData()
			returnToMainMenu(false)
		case 7:
			clearScreen()
			addNewCategory()
			returnToMainMenu(false)
		case 8:
			clearScreen()
			deleteCategory()
			returnToMainMenu(false)
		case 9:
			clearScreen()
			saveShoppingList()
			returnToMainMenu(false)
		case 10:
			clearScreen()
			retrieveShoppingList()
			returnToMainMenu(false)
		default:
			fmt.Println("Please key in the correct input! The program is exiting now!")
			// args = false for returning to main menu
			// args = true for exiting
			returnToMainMenu(true)
			os.Exit(3)
		}
	}

}

func returnToMainMenu(exit bool) {
	if exit {
		fmt.Print("Exiting")
		delayThreeSeconds()
		return
	} else {
		fmt.Print("Returning to Main Menu")
	}
	delayThreeSeconds()
	clearScreen()
	main()
}

func delayThreeSeconds() {
	time.Sleep(delay)
	fmt.Print(".")
	time.Sleep(delay)
	fmt.Print(".")
	time.Sleep(delay)
	fmt.Print(".")
	time.Sleep(delay)
}

func initializeShoppingList(isInit bool, isTest bool) {

	// use pointer for item struct, so we can set the value by referencing
	itemMap = make(map[string]*itemInfo)

	if isInit && !isTest {
		items := []string{"Fork", "Plates", "Cups", "Bread", "Cake", "Coke", "Sprite"}
		itemList := []itemInfo{}

		item1 := createItems(0, 4, 3)
		item2 := createItems(0, 4, 3)
		item3 := createItems(0, 5, 3)
		item4 := createItems(1, 2, 2)
		item5 := createItems(1, 3, 1)
		item6 := createItems(2, 5, 2)
		item7 := createItems(2, 5, 2)

		itemList = append(itemList, item1, item2, item3, item4, item5, item6, item7)

		for index, item := range items {
			// use pointer for item struct, so we can set the value by referencing
			itemMap[item] = &itemList[index]
		}

		category = append(category, []string{"Household", "Food", "Drinks"}...)
	}

	if isTest && !isInit {
		// testing for modifying/deleting items
		items := []string{"Fork", "Plates", "Sprite"}
		itemList := []itemInfo{}

		item1 := createItems(0, 4, 3)
		item2 := createItems(1, 4, 3)
		item3 := createItems(2, 5, 3)
		itemList = append(itemList, item1, item2, item3)

		for index, item := range items {
			// use pointer for item struct, so we can set the value by referencing
			itemMap[item] = &itemList[index]
		}

		category = append(category, []string{"Household", "Food", "Drinks"}...)
	}

}

func createItems(category int, quantity int, unitCost float64) itemInfo {

	newItem := itemInfo{
		category: category,
		quantity: quantity,
		unitCost: unitCost,
	}

	return newItem
}

func mainMenu() {

	// can break line after comma, opening parenthesis,
	// after dot notation and after binary operators
	start := []string{"Shopping List Application", "=========================", "1. View entire shopping list.", "2. Generate Shopping List Report", "3. Add items.",
		"4. Modify items.", "5. Delete items.", "6. Print Current Data", "7. Add New Category Name.", "8. Delete a Category", "9. Save Shopping List",
		"10. Retrieve Shopping List", "Press other keys to exit.", "Select you choice: "}

	for _, message := range start {
		fmt.Println(message)
	}
	fmt.Scan(&userChoice)
}

func displayShoppingList(itemInfoList map[string]*itemInfo, itemCat []string, sorted bool) {

	// sorted = true to show items sorted by category
	// sorted = false show shopping list contents
	if !sorted {
		fmt.Println("Shopping List Contents:")
	}

	if len(itemCat) != 0 {
		for key, item := range itemInfoList {
			fmt.Println("Category:", itemCat[item.getCategory()], "-", "Item: ", key, " Quantity:", item.getQuantity(), "Unit Cost:", item.getUnitCost())
		}
	} else {
		fmt.Println(emptyMessage)
	}
}

func generateReport() {
	var choice int
	reportContents := []string{"Generate Report", "1. Total Cost of each category", "2. List of item by category", "3. Main Menu"}

	for _, message := range reportContents {
		fmt.Println(message)
	}
	fmt.Print("\n\n")
	fmt.Println("Choose your report: ")
	fmt.Scan(&choice)

	if choice == 1 {
		clearScreen()
		fmt.Println("Total Cost by Category.")
		if len(category) != 0 {
			for index, cat := range category {
				var cost float64 = 0
				for _, value := range itemMap {
					if value.getCategory() == 0+index && cat == category[index] {
						cost += totalCost(value)
					}
				}
				fmt.Println(cat, "cost:", cost)
			}
		} else {
			fmt.Println(emptyMessage)
		}

	} else if choice == 2 {
		clearScreen()
		if len(category) != 0 {
			for index, cat := range category {
				filteredItems := make(map[string]*itemInfo)
				for key, item := range itemMap {
					if item.getCategory() == 0+index && cat == category[index] {
						filteredItems[key] = item
					}
				}
				// args = true to show items sorted by category
				// sorted = false show shopping list contents
				displayShoppingList(filteredItems, category, true)
			}
		} else {
			fmt.Println(emptyMessage)
		}

	} else if choice == 3 {
		// args = false for returning to main menu
		// args = true for exiting
		returnToMainMenu(false)
	} else {
		fmt.Println("Please enter the choice between 1 and 3")
		fmt.Print("Returning to Generate Report")
		delayThreeSeconds()
		clearScreen()
		generateReport()

	}

}

func clearScreen() {
	// below code is to clear screen
	fmt.Print("\033c")
}

func addItem() {
	var itemName, itemCategory string
	var itemUnit int
	var itemCost float64

	fmt.Println("What is the name of you item?")
	fmt.Scan(&itemName)
	fmt.Println("What category does it belong to?")
	fmt.Scan(&itemCategory)
	fmt.Println("How many units are there?")
	fmt.Scan(&itemUnit)
	fmt.Println("How much does it cost per unit?")
	fmt.Scan(&itemCost)

	//// convert to Title format for strings input
	//// deprecated in 1.18
	//itemName = strings.Title(strings.ToLower(itemName))
	//itemCategory = strings.Title(strings.ToLower(itemCategory))
	// fmt.Println(itemName, itemCategory)

	// use cases.Title() instead
	itemName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(itemName))
	itemCategory = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(itemCategory))

	// check if key in map
	if _, ok := itemMap[itemName]; !ok {

		addToContent(itemName, itemCategory, itemUnit, itemCost)

		fmt.Println(itemName, "is added to the shopping list.")

	} else {
		fmt.Println(itemName, "is already in the shopping list.")
	}
}

func addToContent(item string, cat string, unit int, cost float64) {

	var updatedCat int

	if len(category) == 0 {
		category = append(category, cat)
		updatedCat = 0
	} else {
		index, found := checkCategory(cat, category)
		if found {
			updatedCat = index
		} else {
			updatedCat = index
			category = append(category, cat)
		}

	}

	itemMap[item] = &itemInfo{
		category: updatedCat,
		quantity: unit,
		unitCost: cost,
	}
}

func checkCategory(cat string, catSlice []string) (int, bool) {
	for i, value := range catSlice {
		if cat == value {
			return i, true
		}
	}
	return len(category), false
}

func modifyItem() {
	var itemName, itemCategory string
	var itemUnit int
	var itemCost float64
	var indexToRm int
	var oldCat string // added

	fmt.Println("Which item would you wish to modify")
	fmt.Scan(&itemName)

	//// convert to Title format for strings input
	//// deprecated in 1.18
	//itemName = strings.Title(strings.ToLower(itemName))

	// use cases.Title() instead
	itemName = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(itemName))

	// check if key in map
	if _, ok := itemMap[itemName]; !ok {
		fmt.Println("There is no", itemName, "in the shopping list.")
	} else {

		oldName := itemName
		oldCat = category[itemMap[itemName].getCategory()] // added
		fmt.Println("Current item name is", itemName, "- Category is", category[itemMap[itemName].getCategory()], "-", "Quantity is", itemMap[itemName].getQuantity(), "-",
			" Unit Cost is", itemMap[itemName].getUnitCost())

		fmt.Println("Enter new name. Enter for no change.")
		if _, err := fmt.Scanln(&itemName); err != nil {
			defer fmt.Println("No changes to item name made")
		}
		fmt.Println("Enter new Category. Enter for no change.")
		if _, err := fmt.Scanln(&itemCategory); err != nil {
			defer fmt.Println("No changes to category made")
		}
		fmt.Println("Enter new Quantity. Enter for no change.")
		if _, err := fmt.Scanln(&itemUnit); err != nil {
			defer fmt.Println("No changes to quantity made")
		}
		fmt.Println("Enter new Unit cost. Enter for no change.")
		if _, err := fmt.Scanln(&itemCost); err != nil {
			defer fmt.Println("No changes to unit cost made")
		}

		//// convert to Title format for strings input
		//// deprecated in 1.18
		//itemCategory = strings.Title(strings.ToLower(itemCategory))

		// use cases.Title() instead
		itemCategory = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(itemCategory))

		appendCat := modifyKeyValue(itemName, oldName, itemCategory, oldCat, itemCost, itemUnit) //modify

		if appendCat { // added condition
			indexToRm = modifyCategoryContents()
			modifyItemContents(indexToRm)
		}
	}

}

func modifyKeyValue(newItem string, oldItem string, cat string, oldCat string, cost float64, quantity int) bool {

	var index int
	var found bool
	// true means will append on to the category slice
	// false means will replace
	// initialise to true unless check that category is replaced
	appendCat := true
	if cat == "" {
		index = itemMap[oldItem].getCategory()
	} else {
		// check if category is in slice
		index, found = checkCategory(cat, category)
		if !found {
			appendCat = checkReplaceOrAppend(oldCat, cat)
			// check category again after replacing/appending category
			// and get the latest index
			index, _ = checkCategory(cat, category)
		}
	}
	if cost == 0 {
		cost = itemMap[oldItem].getUnitCost()
		// itemMap[newItem].setUnitCost(cost)
	}
	if quantity == 0 {
		quantity = itemMap[oldItem].getQuantity()
		// itemMap[newItem].setQuantity(quantity)
	}
	itemMap[newItem] = &itemInfo{
		category: index,
		quantity: quantity,
		unitCost: cost,
	}

	if newItem != oldItem {
		delete(itemMap, oldItem)
	}

	return appendCat

}

// added function
func checkReplaceOrAppend(oldCat string, newCat string) bool {
	var index int
	count := 0
	for _, value := range itemMap {
		if category[value.getCategory()] == oldCat {
			count += 1
		}
	}
	// fmt.Println(count, "count")
	if count == 1 {
		index, _ = checkCategory(oldCat, category)
		// fmt.Println(index, "index")
		category[index] = newCat
		// fmt.Println(category)
		return false
	}
	return true
}

func modifyCategoryContents() int {

	var existItemCat []string
	var indexToRm int

	if len(itemMap) == 0 {
		category = []string{}
	} else {
		existItemCat = getExistCategory()
		indexToRm = getIndexToRm(existItemCat)

	}
	category = existItemCat

	return indexToRm
}

func getExistCategory() []string {

	var existItemCat []string
	inSlice := make(map[string]bool)
	for _, catName := range category { // added condition
		for _, val := range itemMap {
			// this if statement help to remove duplicates
			if catName == category[val.getCategory()] { // added condition
				if _, ok := inSlice[category[val.getCategory()]]; !ok {

					inSlice[category[val.getCategory()]] = true
					existItemCat = append(existItemCat, category[val.getCategory()])
				}
			}
		}
	}
	return existItemCat
}

func getIndexToRm(itemCat []string) int {
	var indexToRm int

	for index, cat := range category {
		// this var check if category in shopping list
		check := []bool{}
		for _, value := range itemCat {
			if value == cat {
				check = append(check, true)
				break
			}
		}
		// if length is 0 means not in shopping list
		if len(check) == 0 {
			// indexToRm = append(indexToRm, index)
			indexToRm = index
		}
	}
	return indexToRm
}

func modifyItemContents(indexToRm int) {

	for _, value := range itemMap {
		if value.getCategory() > indexToRm {

			value.setCategory(value.getCategory() - 1)
		}
	}
}

func deleteItem() {

	var itemToDelete string
	var indexToRm int

	if len(itemMap) == 0 || len(category) == 0 {
		fmt.Println(emptyMessage)
	} else {
		fmt.Println("Delete Item.")
		fmt.Println("Enter item name to delete:")
		fmt.Scanln(&itemToDelete)

		//// convert to Title format for strings input
		//// deprecated in 1.18
		//itemToDelete = strings.Title(strings.ToLower(itemToDelete))

		// use cases.Title() instead
		itemToDelete = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(itemToDelete))

		found := findItem(itemToDelete)
		if found {
			delete(itemMap, itemToDelete)
			fmt.Printf("Deleted %s\n", itemToDelete)
			indexToRm = modifyCategoryContents()
			modifyItemContents(indexToRm)
		} else {
			fmt.Println("Item not found. Nothing to delete.")
			time.Sleep(delay)
			// args = false for returning to main menu
			// args = true for exiting
			returnToMainMenu(false)
		}
	}

}

func findItem(item string) bool {
	if _, ok := itemMap[item]; !ok {
		fmt.Println("Item not found")
		return false
	}
	return true
}

func printData() {
	fmt.Println("Print Current Data")
	if len(itemMap) == 0 {
		fmt.Println("No Data found!")
	} else {
		for key, value := range itemMap {
			fmt.Println(key, "-", *value)
		}
	}
}

func addNewCategory() {
	var itemCat string
	fmt.Println("Add New Category Name")
	fmt.Println("What is the New Category Name to add?")

	if _, err := fmt.Scanln(&itemCat); err != nil {
		defer fmt.Println("No Input Found!")
	} else {
		if index, found := checkCategory(itemCat, category); found {
			fmt.Println("Category:", itemCat, "already in exist at index", index)
		} else {
			fmt.Println("New category:", itemCat, "is added at index", index)
		}
	}

}

func deleteCategory() {

	var userCategory string
	// var indexToRm []int6

	fmt.Println("Please enter the category you want to delete: ")
	fmt.Scan(&userCategory)

	//// change to Title format for strings input
	//// deprecated in 1.18
	//userCategory = strings.Title(strings.ToLower(userCategory))

	// use cases.Title() instead
	userCategory = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(userCategory))

	if index, found := checkCategory(userCategory, category); !found {
		fmt.Println(userCategory + " category does not exists")
	} else {

		removeItemFromMap(index)
		modifyItemContents(index)
		fmt.Println("Deleting", userCategory, "from the shopping list")
		// delayThreeSeconds()
	}
}

func removeItemFromMap(index int) {

	//update map
	for key, value := range itemMap {
		if value.getCategory() == index {
			delete(itemMap, key)
		}
	}

	//get existing category
	category = getExistCategory()
	time.Sleep(delay)
}

func saveShoppingList() {

	var choice string

	// declare a new var without referencing itemInfo struct
	shopList := make(map[string]itemInfo)

	fmt.Println("Do you want to save your shopping list? (Y/N)")
	fmt.Scan(&choice)

	// copy the data to new itemInfo struct without referencing to assign to itemList field of shoppingList struct
	for key, value := range itemMap {
		shopList[key] = *value
	}

	// // convert to uppercase format for strings input
	choice = strings.ToUpper(choice)
	switch choice {
	case "Y":
		shoppingList := shoppingList{
			id:           listCount + 1,
			itemList:     shopList,
			categoryList: category,
		}
		listCount += 1

		// everytime user save will append shoppingList to a slice
		list = append(list, shoppingList)
		fmt.Println(listCount, "shopping list(s) saved!")
	case "N":
		// args = false for returning to main menu
		// args = true for exiting
		returnToMainMenu(false)
	default:
		fmt.Print("Wrong input choice! Please try again")
		delayThreeSeconds()
		fmt.Println("")
		clearScreen()
		saveShoppingList()
	}

}

func retrieveShoppingList() {

	var userIndex int
	newItemMap := make(map[string]*itemInfo)
	var itemList []itemInfo
	var keyList []string

	if len(list) == 0 {
		fmt.Println("No shopping list is saved!")
	} else {
		fmt.Printf("Pleae enter the index (start from 0) of the shopping list you want to retrieve (currently have %v):\n", len(list))
		fmt.Scanln(&userIndex)

		// if user key in index that is out of range, error message
		if userIndex >= len(list) {
			fmt.Println("You do not have so many shopping list saved!")
		} else {
			// print retrieve shopping list using shoppingList struct printShoppingList() method
			list[userIndex].printShoppingList()

			// update category
			category = list[userIndex].getCategoryList()

			// 2 var to store the key(item) and value(itemInfo struct) of the item
			for key, value := range list[userIndex].getItemList() {
				itemList = append(itemList, value)
				keyList = append(keyList, key)
			}
			// re-assign back newly declared itemInfo var with reference,
			// so we are able to update the value for other features to capture
			for index, key := range keyList {
				newItemMap[key] = &itemList[index]
			}

			itemMap = newItemMap
		}
	}
}
