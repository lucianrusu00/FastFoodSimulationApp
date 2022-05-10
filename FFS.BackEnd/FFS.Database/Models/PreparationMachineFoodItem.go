package DatabaseModels

import (
	"gorm.io/gorm"
)

type PreparationMachineFoodItem struct {
	gorm.Model
	Name             string `json:"name"`
	Capacity         int    `json:"capacity"`
	NumberOfMachines int    `json:"number_of_machines"`
}
