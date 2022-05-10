package BusinessLogicModels

type OrderDTO struct {
	Name          string   `json:"name" binding:"required"`
	FoodItemsName []string `json:"food_items_name"`
}
