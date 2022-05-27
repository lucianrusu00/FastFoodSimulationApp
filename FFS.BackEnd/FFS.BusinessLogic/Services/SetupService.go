package BusinessLogicServices

import (
	"BusinessLogicModels"
	"BusinessLogicUtilities"
	"DatabaseModels"

	"fmt"

	"gorm.io/gorm"

	"github.com/lib/pq"
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
		prepMachine := DatabaseModels.PMFoodItem{Name: prepMachineDTO.Name, NumberOfMachines: prepMachineDTO.NumberOfMachines, Capacity: prepMachineDTO.Capacity}
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

	var ingredientJustAdded DatabaseModels.Ingredient

	x.Db.Preload("PreparationMachine").First(&ingredientJustAdded, "Name = ?", ingredientDTO.Name)

	fmt.Printf("Added ingredient with name %s and preparation machine %s with capacity %d \n", ingredientJustAdded.Name, ingredientJustAdded.PreparationMachine.Name, ingredientJustAdded.PreparationMachine.Capacity)
}

func (x SetupService) AddFoodItem(foodItemDTO BusinessLogicModels.FoodItemDTO) {

	var preparationMachineFoodItem DatabaseModels.PMFoodItem

	fmt.Printf("----------------------------------------------The preparation machine for the food item %s is: ", foodItemDTO.Name)

	x.Db.First(&preparationMachineFoodItem, "Name = ?", foodItemDTO.PreparationMachineName)
	fmt.Println(preparationMachineFoodItem.Name)

	var ingredientList []DatabaseModels.Ingredient

	x.Db.Preload("PreparationMachine").Where("name in ?", foodItemDTO.IngredientNameList).Find(&ingredientList) // this is how you include diferent nested structs in the select statement

	for _, ingredient := range ingredientList {
		fmt.Printf("Ingredient %s with prep machine %s and capacity %d \n", ingredient.Name, ingredient.PreparationMachine.Name, ingredient.PreparationMachine.Capacity)
	}

	for _, foodItemString := range foodItemDTO.FoodItemNameList {
		var foodItemFromDB DatabaseModels.FoodItem

		x.Db.Preload("PMFoodItem").First(&foodItemFromDB, "Name = ?", foodItemString)

		fmt.Printf("FoodItem %s with preparation machine %s\n", foodItemFromDB.Name, foodItemFromDB.PMFoodItem.Name)
	}

	foodItem := DatabaseModels.FoodItem{Name: foodItemDTO.Name, IngredientList: ingredientList, PMFoodItem: preparationMachineFoodItem, PreparationTime: foodItemDTO.PreparationTime, FoodItemsString: pq.StringArray(foodItemDTO.FoodItemNameList)}

	fmt.Printf("======================The food item %s with prep machine %s was created\n", foodItem.Name, foodItem.PMFoodItem.Name)
	x.Db.Create(&foodItem)

	var foodItemJustCreated DatabaseModels.FoodItem
	//x.Db.Preload("PMFoodItem").Preload("FoodItemsString").First(&foodItemJustCreated, "Name = ?", foodItemDTO.Name)

	x.Db.Preload("PMFoodItem").First(&foodItemJustCreated, "Name = ?", foodItemDTO.Name)

	fmt.Printf("Added food item %s with prep machine %s and the food items: ", foodItemJustCreated.Name, foodItemJustCreated.PMFoodItem.Name)
	for _, foodItemStingAddedNow := range foodItemJustCreated.FoodItemsString {
		fmt.Printf("%s ", foodItemStingAddedNow)
	}
	fmt.Println()
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

func (x SetupService) getAllRecipe(foodItem *DatabaseModels.FoodItem) {
	for i := range foodItem.FoodItemsString {
		var foodItemNew DatabaseModels.FoodItem
		x.Db.Preload("PMFoodItem").Preload("IngredientList").Preload("IngredientList.PreparationMachine").First(&foodItemNew, "Name = ?", foodItem.FoodItemsString[i])
		x.getAllRecipe(&foodItemNew)

		foodItem.FoodItems = append(foodItem.FoodItems, foodItemNew)
	}
}

func (x SetupService) AddOrder(orderDTO BusinessLogicModels.OrderDTO) {
	var foodItemList []DatabaseModels.FoodItem

	x.Db.Preload("PMFoodItem").Preload("IngredientList").Preload("IngredientList.PreparationMachine").Where("name in ?", orderDTO.FoodItemsName).Find(&foodItemList)

	//x.Db.Preload("PMFoodItem").Preload("IngredientList").Preload("IngredientList.PreparrationMachine").Preload("FoodItems").Preload("FoodItems.PMFoodItems").Find(&foodItemsDB)

	for i := range foodItemList {
		x.getAllRecipe(&foodItemList[i])
	}

	order := DatabaseModels.Order{Name: orderDTO.Name, FoodItems: foodItemList}

	x.Db.Create(&order)

	go x.Simulate.AddOrderToQueue(order)
}
