package main

import (
	"fmt"
)

type shoppingList struct {
	id           int
	itemList     map[string]itemInfo
	categoryList []string
}

func (s shoppingList) getId() int {
	return s.id
}

func (s shoppingList) getCategoryList() []string {
	return s.categoryList
}

func (s shoppingList) getItemList() map[string]itemInfo {
	return s.itemList
}

func (s shoppingList) printShoppingList() {

	fmt.Println("Shopping list", s.getId(), "consists:")
	for key, value := range s.getItemList() {

		// fmt.Println(key, value)
		fmt.Println(key, "- Category is", s.getCategoryList()[value.getCategory()], "-", "Quantity is", value.getQuantity(), "-",
			" Unit Cost is", value.getUnitCost())

	}
}
