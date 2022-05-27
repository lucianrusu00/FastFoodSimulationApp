package DatabaseModels

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type FoodItem struct {
	gorm.Model
	Name         string `json:"name"`
	PMFoodItemID int
	PMFoodItem   PMFoodItem `json:"assembly_point"`

	IngredientList  []Ingredient   `json:"ingredient_list" gorm:"many2many:foodItem_ingredients"`
	FoodItemsString pq.StringArray `json:"food_items_strings" gorm:"type:text[]"`
	FoodItems       []FoodItem     `gorm:"-"`
	PreparationTime int            `json:"preparation_time"`
}
