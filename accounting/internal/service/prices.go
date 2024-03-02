package service

import "math/rand/v2"

const (
	assignTaskMinPrice   = 10
	assignTaskMaxPrice   = 20
	completeTaskMinPrice = 20
	completeTaskMaxPrice = 40
)

func priceAssignTask() int {
	return assignTaskMinPrice + rand.IntN(assignTaskMaxPrice-assignTaskMinPrice+1)
}

func priceCompleteTask() int {
	return completeTaskMinPrice + rand.IntN(completeTaskMaxPrice-completeTaskMinPrice+1)
}
