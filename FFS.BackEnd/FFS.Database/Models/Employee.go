package DatabaseModels

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name              string               `json:"name"`
	UsableMachineList []PreparationMachine `json:"preparationMachines" gorm:"foreignKey:ID"`
}
