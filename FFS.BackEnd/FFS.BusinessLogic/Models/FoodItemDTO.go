package BusinessLogicModels

type FoodItemDTO struct {
	Name                   string   `json:"name"`
	IngredientNameList     []string `json:"ingredient_name_list"`
	PreparationMachineName string   `json:"preparation_machine_name"`
	FoodItemNameList       []string `json:"food_item_name_list"`
	PreparationTime        int      `json:"preparation_time"`
}
