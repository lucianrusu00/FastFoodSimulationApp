package DatabaseModels

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Name      string     `json:"name" binding:"required"`
	FoodItems []FoodItem `json:"food_items" gorm:"many2many:order_foodItems"`
}
