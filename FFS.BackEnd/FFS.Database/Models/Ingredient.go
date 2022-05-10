package DatabaseModels

import (
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name               string             `json:"name"`
	PreparationTime    int                `json:"preparation_time"`
	PreparationMachine PreparationMachine `json:"preparation_machine" gorm:"foreignKey:ID"`
}
