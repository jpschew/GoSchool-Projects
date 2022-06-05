package main

type cost interface {
	totalCost() float64
}

func totalCost(c cost) float64 {
	return c.totalCost()
}
