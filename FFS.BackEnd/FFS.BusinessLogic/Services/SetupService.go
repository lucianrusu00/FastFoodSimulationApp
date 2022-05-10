package BusinessLogicServices

import (
	"BusinessLogicModels"
	"BusinessLogicUtilities"
	"DatabaseModels"

	"fmt"

	"gorm.io/gorm"
)

type SetupService struct {
	Db       *gorm.DB
	Simulate BusinessLogicUtilities.Simulate
}

func (x SetupService) AddPrepMachine(prepMachineDTO BusinessLogicModels.PreparationMachineDTO) {
	if prepMachineDTO.Type == "Ingredient" {
		prepMachine := DatabaseModels.PreparationMachine{Name: prepMachineDTO.Name, NumberOfMachines: prepMachineDTO.NumberOfMachines, Capacity: prepMachineDTO.Capacity}
		x.Db.Create(&prepMachine)
	} else {
		prepMachine := DatabaseModels.PreparationMachineFoodItem{Name: prepMachineDTO.Name, NumberOfMachines: prepMachineDTO.NumberOfMachines, Capacity: prepMachineDTO.Capacity}
		x.Db.Create(&prepMachine)
	}

	// Add channel for preparation machine

	fmt.Println("preparation machine added")

	// Whenever a preparation machine gets added, create a buffered channel with name prepMachine.Name and capcity equal to prepMachine.Capacity
	// queue := make(chan DatabaseModels.Ingredient, prepMachineDTO.Capacity)
}

func (x SetupService) AddIngredient(ingredientDTO BusinessLogicModels.IngredientDTO) {
	var preparationMachine DatabaseModels.PreparationMachine

	x.Db.First(&preparationMachine, "Name = ?", ingredientDTO.PreparationMachineName)

	ingredient := DatabaseModels.Ingredient{Name: ingredientDTO.Name, PreparationTime: ingredientDTO.PreparationTime, PreparationMachine: preparationMachine}

	x.Db.Create(&ingredient)

	fmt.Printf("Added ingredient with name %s and preparation machine %s with capacity %d \n", ingredient.Name, ingredient.PreparationMachine.Name, ingredient.PreparationMachine.Capacity)
}

func (x SetupService) AddFoodItem(foodItemDTO BusinessLogicModels.FoodItemDTO) {
	var preparationMachineFoodItem DatabaseModels.PreparationMachineFoodItem

	x.Db.First(&preparationMachineFoodItem, "Name = ?", foodItemDTO.PreparationMachineName)

	var ingredientList []DatabaseModels.Ingredient

	x.Db.Preload("PreparationMachine").Where("name in ?", foodItemDTO.IngredientNameList).Find(&ingredientList) // this is how you include diferent nested structs in the select statement

	for _, ingredient := range ingredientList {
		fmt.Printf("Ingredient %s with prep machine %s and capacity %d \n", ingredient.Name, ingredient.PreparationMachine.Name, ingredient.PreparationMachine.Capacity)
	}

	var foodItemList []DatabaseModels.FoodItem

	x.Db.Preload("PreparationMachineFoodItem").Preload("IngredientList").Preload("FoodItems").Where("name in ?", foodItemDTO.FoodItemNameList).Find(&foodItemList)

	foodItem := DatabaseModels.FoodItem{Name: foodItemDTO.Name, IngredientList: ingredientList, PreparationMachineFoodItem: preparationMachineFoodItem, PreparationTime: foodItemDTO.PreparationTime, FoodItems: foodItemList}

	x.Db.Create(&foodItem)

	//fmt.Printf("Added foodItem with name %s and ingredients: %s with prep machine %s and %s with prep machine %s \n", foodItem.Name, foodItem.IngredientList[0].Name, foodItem.IngredientList[0].PreparationMachine.Name, foodItem.IngredientList[1].Name, foodItem.IngredientList[1].PreparationMachine.Name)
}

func (x SetupService) AddEmployee(employeeDTO BusinessLogicModels.EmployeeDTO) {
	var preparationMachineList []DatabaseModels.PreparationMachine

	x.Db.Where("name in ?", employeeDTO.UsableMachineNameList).Find(&preparationMachineList)

	x.Db.Create(&DatabaseModels.Employee{Name: employeeDTO.Name, UsableMachineList: preparationMachineList})
}

func (x SetupService) StartSimulation() {
	go x.Simulate.StartSimulation()
}

func (x SetupService) AddOrder(orderDTO BusinessLogicModels.OrderDTO) {
	var foodItemList []DatabaseModels.FoodItem

	x.Db.Preload("PreparationMachineFoodItem").Preload("IngredientList").Preload("IngredientList.PreparationMachine").Preload("FoodItems").Where("name in ?", orderDTO.FoodItemsName).Find(&foodItemList)

	order := DatabaseModels.Order{Name: orderDTO.Name, FoodItems: foodItemList}

	x.Db.Create(&order)

	go x.Simulate.AddOrderToQueue(order)
}
