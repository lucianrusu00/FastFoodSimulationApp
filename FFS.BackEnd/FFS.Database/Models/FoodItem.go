package DatabaseModels

import (
	"gorm.io/gorm"
)

type FoodItem struct {
	gorm.Model
	Name                       string                     `json:"name"`
	IngredientList             []Ingredient               `json:"ingredient_list" gorm:"foreignKey:ID"`
	FoodItems                  []FoodItem                 `json:"food_item_list" gorm:"foreignKey:ID"`
	PreparationMachineFoodItem PreparationMachineFoodItem `json:"assembly_point" gorm:"foreignKey:ID"`
	PreparationTime            int                        `json:"preparation_time"`
}
