package BusinessLogicModels

type PreparationMachineDTO struct {
	Name             string `json:"name"`
	NumberOfMachines int    `json:"number_of_machines"`
	Capacity         int    `json:"capacity"`
	Type             string `json:"type"` // Type of the machine. Either Ingredient or FoodItem
}
