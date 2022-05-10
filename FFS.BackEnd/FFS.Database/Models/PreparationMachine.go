package DatabaseModels

import (
	"gorm.io/gorm"
)

type PreparationMachine struct {
	gorm.Model
	Name             string `json:"name" binding:"required"`
	NumberOfMachines int    `json:"number_of_machines"`
	Capacity         int    `json:"capacity" binding:"required"`
}
