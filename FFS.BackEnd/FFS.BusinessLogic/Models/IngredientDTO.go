package BusinessLogicModels

type IngredientDTO struct {
	Name                   string `json:"Name"`
	PreparationTime        int    `json:"preparation_time"`
	PreparationMachineName string `json:"preparation_machine_name"`
}
