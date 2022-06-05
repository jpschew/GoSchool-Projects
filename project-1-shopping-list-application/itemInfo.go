package main

type itemInfo struct {
	category int
	quantity int
	unitCost float64
}

func (i itemInfo) getCategory() int {
	return i.category
}

func (i itemInfo) getQuantity() int {
	return i.quantity
}

func (i itemInfo) getUnitCost() float64 {
	return i.unitCost
}

//need use * ptr so that we get the address of i
//to update the value for setter
//however map cannot call pointer
func (i *itemInfo) setCategory(cat int) {
	i.category = cat
}

// func (i *itemInfo) setQuantity(q int) {
// 	i.quantity = q
// }

// func (i *itemInfo) setUnitCost(u float64) {
// 	i.unitCost = u
// }

func (i itemInfo) totalCost() float64 {
	return float64(i.quantity) * i.unitCost
}
