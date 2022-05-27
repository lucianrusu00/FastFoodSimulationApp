package DatabaseModels

import (
	"gorm.io/gorm"
)

type PMFoodItem struct {
	gorm.Model
	Name             string `json:"name" binding:"required"`
	Capacity         int    `json:"capacity" binding:"required"`
	NumberOfMachines int    `json:"number_of_machines"`
}
