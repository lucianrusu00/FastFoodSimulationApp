package BusinessLogicUtilities

import (
	"DatabaseModels"
	"fmt"

	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Simulate struct {
	Db     *gorm.DB
	Orders chan DatabaseModels.Order
}

// type PreparationMachine struct {
// 	Name             string
// 	Capacity         int
// 	NumberOfMachines int
// }

// type Ingredient struct {
// 	Name               string
// 	PreparationMachine PreparationMachine
// 	PreparationTime    int // number of seconds nedeed for preparation
// }

// type FoodItem struct {
// 	Name               string
// 	Ingredients        []Ingredient
// 	FoodItems          []FoodItem
// 	PreparationMachine PreparationMachine
// 	PreparationTime    int
// }

// type Order struct {
// 	Name      string
// 	FoodItems []FoodItem
// }

// type Employee struct {
// 	Name string
// }

// This will the buffer of the machine
func createPreparationMachineIngredientBuffer(preparationMachine DatabaseModels.PreparationMachine, buffers map[string]([]chan DatabaseModels.Ingredient), capacityControl map[string]([]chan bool), preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient))) {
	preparationMachinesIngredients[preparationMachine.Name] = make(map[string]chan DatabaseModels.Ingredient)
	capacityControl[preparationMachine.Name] = make([]chan bool, preparationMachine.NumberOfMachines)
	buffers[preparationMachine.Name] = make([]chan DatabaseModels.Ingredient, preparationMachine.NumberOfMachines)
	for i := 0; i < preparationMachine.NumberOfMachines; i++ {
		capacityControl[preparationMachine.Name][i] = make(chan bool, preparationMachine.Capacity)
		buffers[preparationMachine.Name][i] = make(chan DatabaseModels.Ingredient, preparationMachine.Capacity)
	}

}

func createPreparationMachineFoodItemBuffer(preparationMachine DatabaseModels.PreparationMachineFoodItem, buffers map[string]([]chan DatabaseModels.FoodItem), capacityControl map[string]([]chan bool), preparationMachinesFoodItems map[string](map[string](chan DatabaseModels.FoodItem))) {
	preparationMachinesFoodItems[preparationMachine.Name] = make(map[string]chan DatabaseModels.FoodItem)
	capacityControl[preparationMachine.Name] = make([]chan bool, preparationMachine.NumberOfMachines)
	buffers[preparationMachine.Name] = make([]chan DatabaseModels.FoodItem, preparationMachine.NumberOfMachines)
	for i := 0; i < preparationMachine.NumberOfMachines; i++ {
		capacityControl[preparationMachine.Name][i] = make(chan bool, preparationMachine.Capacity)
		buffers[preparationMachine.Name][i] = make(chan DatabaseModels.FoodItem, preparationMachine.Capacity)
	}

}

// Send the ingredient to be prepared
func addIngredientToPrepMachine(ingredient DatabaseModels.Ingredient, channels map[string](map[string](chan DatabaseModels.Ingredient))) {
	//channels[ingredient.PreparationMachine.Name] = make(map[string]chan Ingredient)
	channels[ingredient.PreparationMachine.Name][ingredient.Name] = make(chan DatabaseModels.Ingredient, ingredient.PreparationMachine.Capacity)
}

func addFoodItemToPrepMachine(foodItem DatabaseModels.FoodItem, channels map[string](map[string](chan DatabaseModels.FoodItem))) {
	//channels[foodItem.PreparationMachine.Name] = make(map[string]chan FoodItem)
	channels[foodItem.PreparationMachineFoodItem.Name][foodItem.Name] = make(chan DatabaseModels.FoodItem, foodItem.PreparationMachineFoodItem.Capacity)
}

func runPreparationMachineIngredient(preparationMachine DatabaseModels.PreparationMachine, capacityControl map[string]([]chan bool), preparationMachinesBuffer map[string]([]chan DatabaseModels.Ingredient), preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient)), machineNumber int, employeeList *EmployeeList) {
	fmt.Printf("-- Setup --: Run %s number %d\n", preparationMachine.Name, machineNumber)
	for {
		// <- preparationMachinesBuffer[preparationMachine.Name][machineNumber]
		ingredient := <-preparationMachinesBuffer[preparationMachine.Name][machineNumber]
		fmt.Printf("-- PrepMahineIngredient Thread --: Received ingredient %s in ingredient machine %s number %d\n", ingredient.Name, preparationMachine.Name, machineNumber)

		go func(i DatabaseModels.Ingredient, employeeList *EmployeeList) {

			fmt.Printf("-- PrepMahineIngredient Thread --: Assign Employee \n")
			employee := DatabaseModels.Employee{}
			employee = assignEmployee(employeeList, employee)
			fmt.Printf("-- PrepMahineIngredient Thread --: Assigned %s\n", employee.Name)

			fmt.Printf("-- PrepMahineIngredient Thread --: Waiting for ingredient %s to be ready.\n", i.Name)
			time.Sleep(time.Duration(i.PreparationTime) * time.Second)

			preparationMachinesIngredients[preparationMachine.Name][i.Name] <- i // sent the prepared ingredient
			fmt.Printf("-- PrepMahineIngredientThread --: Ingredient %s has been sent back to the foodItem machine thread\n", i.Name)
			<-capacityControl[preparationMachine.Name][machineNumber] // frees one capacity cell

			freeEmployee(employeeList, employee)
			fmt.Printf("-- PrepMahineIngredient Thread --: Freed employee %s\n", employee.Name)
		}(ingredient, employeeList)

	}
}

func runPreparationMachineFoodItem(preparationMachine DatabaseModels.PreparationMachineFoodItem, capacityControl map[string]([]chan bool), preparationMachinesFoodItemsBuffer map[string]([]chan DatabaseModels.FoodItem), preparationMachinesIngredientsBuffer map[string]([]chan DatabaseModels.Ingredient), preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient)), preparationMachineFoodItems map[string](map[string](chan DatabaseModels.FoodItem)), machineNumber int, employeeList *EmployeeList) {
	fmt.Printf("-- Setup --: Run %s number %d\n", preparationMachine.Name, machineNumber)
	for {
		fmt.Printf("-- PrepMahineFoodItem Thread --: Waiting for a food Item to prepare\n\n")
		foodItem := <-preparationMachinesFoodItemsBuffer[preparationMachine.Name][machineNumber] // takes the foodItem out of the buffer
		fmt.Printf("-- PrepMahineFoodItem Thread --: Food Item %s received in preparation machine %s number %d\n", foodItem.Name, preparationMachine.Name, machineNumber)
		//foodItem := <-preparationMachineFoodItems[preparationMachine.Name]
		// send ingredient to it's prepMachine
		for _, foodI := range foodItem.FoodItems {
			fmt.Printf("-- PrepMahineFoodItem Thread --: Sending foodItem %s to one of it's prepMachines\n", foodI.Name)
			for k := 0; k < foodI.PreparationMachineFoodItem.NumberOfMachines; k++ {
				select {
				case capacityControl[foodI.PreparationMachineFoodItem.Name][k] <- true:
					preparationMachinesFoodItemsBuffer[foodI.PreparationMachineFoodItem.Name][k] <- foodI
					// sent the ingredient to be prepared ^^^
					fmt.Printf("-- PrepMahineFoodItem Thread --: Selected prep machine: %s, number %d for foodItem %s\n", foodI.PreparationMachineFoodItem.Name, k, foodI.Name)
					k = foodI.PreparationMachineFoodItem.NumberOfMachines
				default:
					fmt.Printf("-- PrepMachineFoodItem Thread --: Machine %s number %d is full, going to try machine %d\n", foodI.PreparationMachineFoodItem.Name, k, k+1)
				}
				if k == foodI.PreparationMachineFoodItem.NumberOfMachines-1 {
					fmt.Printf("-- ERROR -- PrepMahineFoodItem Thread --: Empty machine not found for FoodItem %s\n", foodI.Name)
				}
			}
		}

		for _, i := range foodItem.IngredientList {
			fmt.Printf("-- PrepMahineFoodItem Thread --: Sending ingredient %s to one of it's prepMachines\n", i.Name)
			for k := 0; k < i.PreparationMachine.NumberOfMachines; k++ {
				select {
				case capacityControl[i.PreparationMachine.Name][k] <- true:
					preparationMachinesIngredientsBuffer[i.PreparationMachine.Name][k] <- i
					// sent the ingredient to be prepared ^^^
					fmt.Printf("-- PrepMahineFoodItem Thread --: Selected prep machine: %s, number %d for ingredient %s\n", i.PreparationMachine.Name, k, i.Name)
					k = i.PreparationMachine.NumberOfMachines
				default:
					fmt.Printf("-- PrepMachineFoodItem Thread --: Machine %s number %d is full, going to try machine %d\n", i.PreparationMachine.Name, k, k+1)
				}
				if k == i.PreparationMachine.NumberOfMachines-1 {
					fmt.Printf("-- ERROR -- PrepMahineFoodItem Thread --: Empty machine not found for Ingredient %s\n", i.Name)
				}
			}

		}

		go func(fI DatabaseModels.FoodItem, employeeList *EmployeeList) { // wait for the ingredients to be ready
			for _, foodI := range fI.FoodItems {
				foodItemMade := <-preparationMachineFoodItems[foodI.PreparationMachineFoodItem.Name][foodI.Name]
				fmt.Printf("-- PrepMahineFoodItem Thread --: FoodItem %s has been prepared and received for assembly\n", foodItemMade.Name)
			}

			for _, i := range fI.IngredientList {
				ingredientMade := <-preparationMachinesIngredients[i.PreparationMachine.Name][i.Name]
				fmt.Printf("-- PrepMahineFoodItem Thread --: Ingredient %s has been prepared and received for assembly\n", ingredientMade.Name)
			}

			fmt.Printf("-- PrepMahineFoodItem Thread --: Assigning employee \n")
			employee := DatabaseModels.Employee{}
			employee = assignEmployee(employeeList, employee)
			fmt.Printf("-- PrepMahineFoodItem Thread --: Assigned %s\n", employee.Name)

			fmt.Printf("-- PrepMahineFoodItem Thread --: Waiting %d seconds for assemblying %s\n", fI.PreparationTime, fI.Name)
			time.Sleep(time.Duration(fI.PreparationTime) * time.Second)
			fmt.Printf("-- PrepMahineFoodItem Thread --: Food Item %s is finished and is being sent to the order thread from %s\n", fI.Name, preparationMachine.Name)
			fmt.Println(preparationMachineFoodItems[preparationMachine.Name][fI.Name])
			preparationMachineFoodItems[preparationMachine.Name][fI.Name] <- fI // send the made food item back to the order
			fmt.Printf("-- PrepMahineFoodItem Thread --: Food Item %s has been sent to the order thread\n", fI.Name)
			<-capacityControl[preparationMachine.Name][machineNumber] // gives one more capacity

			freeEmployee(employeeList, employee)
			fmt.Printf("-- PrepMahineFoodItem Thread --: Freed Employee %s\n", employee.Name)
		}(foodItem, employeeList)

	}
}

func runOrderThread(preparationMachinesFoodItemsBuffer map[string]([]chan DatabaseModels.FoodItem), capacityControl map[string]([]chan bool), preparationMachineFoodItems map[string](map[string](chan DatabaseModels.FoodItem)), ordersSent chan DatabaseModels.Order, ordersReady chan DatabaseModels.Order) {
	fmt.Printf("-- Setup --: Run order thread\n")
	for {
		order := <-ordersSent
		fmt.Printf("-- Order Thread --: Received order %s\n", order.Name)

		for _, f := range order.FoodItems {
			fmt.Printf("-- Order Thread --: Sending food item %s to it's prep machine %s\n", f.Name, f.PreparationMachineFoodItem.Name)
			for k := 0; k < f.PreparationMachineFoodItem.NumberOfMachines; k++ {
				fmt.Printf("-- Order Thread --: Trying to send food item %s to prep machine number %d---------------------------------------------------------------------------\n", f.Name, k)
				select {
				case capacityControl[f.PreparationMachineFoodItem.Name][k] <- true:
					preparationMachinesFoodItemsBuffer[f.PreparationMachineFoodItem.Name][k] <- f
					fmt.Printf("-- Order Thread --: Selected prep machine: %s, number %d for FoodItem %s\n", f.PreparationMachineFoodItem.Name, k, f.Name)
					k = f.PreparationMachineFoodItem.NumberOfMachines
				default:
					fmt.Printf("-- Order Thread --: Machine %s number %d is full, going to try machine %d, total machine number = %d---------------------------------------------------------------------------\n", f.PreparationMachineFoodItem.Name, k, k+1, f.PreparationMachineFoodItem.NumberOfMachines)
				}
				if k == f.PreparationMachineFoodItem.NumberOfMachines-1 {
					fmt.Printf("-- ERROR -- Order Thread --: Empty machine not found for FoodItem %s\n", f.Name)
				}
			}
		}

		go func(o DatabaseModels.Order) { // wait for the ingredients to be ready
			fmt.Printf("-- Order Thread --: Starting to wait for order %s to be ready\n", o.Name)

			// This can receive food items in any order
			for _, f := range o.FoodItems {
				fmt.Printf("-- Order Thread --: Waiting for food item %s to be send back to the order thread from prep machine %s\n", f.Name, f.PreparationMachineFoodItem.Name)
				fmt.Println(preparationMachineFoodItems[f.PreparationMachineFoodItem.Name][f.Name])
				foodItemMade := <-preparationMachineFoodItems[f.PreparationMachineFoodItem.Name][f.Name]
				fmt.Printf("-- Order Thread --: Food item %s received back to the order thread\n", foodItemMade.Name)
			}
			fmt.Printf("-- Order Thread --: Order %s ready and being sent to ordersReady\n", o.Name)
			ordersReady <- o
			fmt.Printf("-- Order Thread --: Order %s has been sent to the oredersReady thread\n", o.Name)
			<-capacityControl["Order"][0]
		}(order)
	}
}

func sendOrder(order DatabaseModels.Order, capacityControl map[string]([]chan bool), ordersSent chan DatabaseModels.Order) {
	fmt.Printf("-- Basic THREAD --: Trying to send order %s\n", order.Name)

	capacityControl["Order"][0] <- true
	ordersSent <- order

	fmt.Printf("-- Basic THREAD --: Order %s sent\n", order.Name)

}

type EmployeeList struct {
	Lock      chan bool
	Employees []DatabaseModels.Employee
}

func createEmployeeList() EmployeeList {
	employees := []DatabaseModels.Employee{}
	lock := make(chan bool, 1)

	employeeList := EmployeeList{Lock: lock, Employees: employees}
	employeeList.Lock <- true

	return employeeList
}

func (x Simulate) addEmployee(list *EmployeeList, employee DatabaseModels.Employee) {
	<-list.Lock
	list.Employees = append(list.Employees, employee)
	list.Lock <- true
}

func extractEmployee(list *EmployeeList, index int) {
	list.Employees[index] = list.Employees[len(list.Employees)-1]
	list.Employees = list.Employees[:len(list.Employees)-1]
}

func freeEmployee(list *EmployeeList, usedEmployee DatabaseModels.Employee) {
	<-list.Lock

	list.Employees = append(list.Employees, usedEmployee)

	list.Lock <- true
}

func assignEmployee(list *EmployeeList, lastEmployee DatabaseModels.Employee) DatabaseModels.Employee {
	<-list.Lock

	//Busy wait untill an employee is free
	for len(list.Employees) <= 0 {
		list.Lock <- true
		time.Sleep(10)
		<-list.Lock
	}

	if len(list.Employees) == 1 {
		employeeSelected := list.Employees[0]
		extractEmployee(list, 0)

		list.Lock <- true
		return employeeSelected
	}

	//Try to assign the same employee that did the job last
	for index, employee := range list.Employees {
		if employee.Name == lastEmployee.Name {
			extractEmployee(list, index)
			list.Lock <- true
			return employee
		}

	}

	//Assign a random employee
	employeeIndex := rand.Intn(len(list.Employees) - 1)
	employeeSelected := list.Employees[employeeIndex]
	extractEmployee(list, employeeIndex)

	list.Lock <- true
	return employeeSelected

}

func (x Simulate) getEmployeeList() EmployeeList {
	var employees []DatabaseModels.Employee

	x.Db.Find(&employees)

	lock := make(chan bool, 1)

	employeeList := EmployeeList{Lock: lock, Employees: employees}
	employeeList.Lock <- true

	return employeeList

}

func (x Simulate) initiatePreparationMachcinesIngredients(employeeList *EmployeeList, capacityControl map[string]([]chan bool), preparationMachinesIngredientsBuffer map[string][]chan DatabaseModels.Ingredient, preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient))) {
	var preparationMachinesIngredientsDB []DatabaseModels.PreparationMachine

	x.Db.Find(&preparationMachinesIngredientsDB)

	for _, prepMachine := range preparationMachinesIngredientsDB {
		createPreparationMachineIngredientBuffer(prepMachine, preparationMachinesIngredientsBuffer, capacityControl, preparationMachinesIngredients)
		for i := 0; i < prepMachine.NumberOfMachines; i++ {
			go runPreparationMachineIngredient(prepMachine, capacityControl, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, i, employeeList)
		}
	}
}

func (x Simulate) initiatePreparationMachinesFoodItem(employeeList *EmployeeList, capacityControl map[string]([]chan bool), preparationMachinesFoodItems map[string](map[string](chan DatabaseModels.FoodItem)), preparationMachinesFoodItemsBuffer map[string]([]chan DatabaseModels.FoodItem), preparationMachinesIngredientsBuffer map[string][]chan DatabaseModels.Ingredient, preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient))) {
	var preparationMachinesFoodItemDB []DatabaseModels.PreparationMachineFoodItem

	x.Db.Find(&preparationMachinesFoodItemDB)

	for _, prepMachine := range preparationMachinesFoodItemDB {
		createPreparationMachineFoodItemBuffer(prepMachine, preparationMachinesFoodItemsBuffer, capacityControl, preparationMachinesFoodItems)

		for i := 0; i < prepMachine.NumberOfMachines; i++ {
			go runPreparationMachineFoodItem(prepMachine, capacityControl, preparationMachinesFoodItemsBuffer, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, preparationMachinesFoodItems, i, employeeList)
		}
	}

}

func (x Simulate) initiatePreparationMachines(employeeList *EmployeeList, capacityControl map[string]([]chan bool), preparationMachinesIngredientsBuffer map[string][]chan DatabaseModels.Ingredient, preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient)), preparationMachinesFoodItems map[string](map[string](chan DatabaseModels.FoodItem)), preparationMachinesFoodItemsBuffer map[string]([]chan DatabaseModels.FoodItem)) {
	x.initiatePreparationMachcinesIngredients(employeeList, capacityControl, preparationMachinesIngredientsBuffer, preparationMachinesIngredients)
	x.initiatePreparationMachinesFoodItem(employeeList, capacityControl, preparationMachinesFoodItems, preparationMachinesFoodItemsBuffer, preparationMachinesIngredientsBuffer, preparationMachinesIngredients)
}

func (x Simulate) initiateIngredients(preparationMachinesIngredients map[string](map[string](chan DatabaseModels.Ingredient))) {
	var ingredientsDB []DatabaseModels.Ingredient

	x.Db.Preload("PreparationMachine").Find(&ingredientsDB)

	for _, ingredientDB := range ingredientsDB {
		preparationMachinesIngredients[ingredientDB.PreparationMachine.Name] = make(map[string]chan DatabaseModels.Ingredient)
		addIngredientToPrepMachine(ingredientDB, preparationMachinesIngredients)
	}

}

func (x Simulate) initiateFoodItems(preparationMachinesFoodItems map[string](map[string](chan DatabaseModels.FoodItem))) {

	var foodItemsDB []DatabaseModels.FoodItem

	x.Db.Preload("PreparationMachineFoodItem").Find(&foodItemsDB)

	for _, foodItemDB := range foodItemsDB {
		addFoodItemToPrepMachine(foodItemDB, preparationMachinesFoodItems)
	}

}

func (x Simulate) AddOrderToQueue(order DatabaseModels.Order) {

	x.Orders <- order
}

func (x Simulate) StartSimulation() {

	//x.orders = make(chan DatabaseModels.Order, 100) // We can add 100 orders to queue

	employeeList := x.getEmployeeList()

	//Only for the intial test
	// x.addEmployee(&employeeList, DatabaseModels.Employee{Name: "Boghy"})
	// x.addEmployee(&employeeList, DatabaseModels.Employee{Name: "Dan"})
	// x.addEmployee(&employeeList, DatabaseModels.Employee{Name: "Lucian"})
	// x.addEmployee(&employeeList, DatabaseModels.Employee{Name: "Sion"})

	//////////////////////////////////////////////////////////// NO NO
	// <-employeeList.Lock
	// for _, employee := range employeeList.Employees {
	// 	fmt.Println(employee.Name)
	// }
	// employeeList.Lock <- true

	// return

	capacityControl := make(map[string]([]chan bool))
	ordersReady := make(chan DatabaseModels.Order, 100) // 100 orders can wait to be received by the customer
	ordersSent := make(chan DatabaseModels.Order, 20)

	go func() {
		for {
			orderReady := <-ordersReady
			fmt.Printf("-- OrdersReady Thread --: Order %s is ready\n", orderReady.Name)
		}
	}()

	preparationMachinesIngredientsBuffer := make(map[string][]chan DatabaseModels.Ingredient)
	preparationMachinesIngredients := make(map[string](map[string](chan DatabaseModels.Ingredient)))
	preparationMachinesFoodItems := make(map[string](map[string](chan DatabaseModels.FoodItem)))
	preparationMachinesFoodItemsBuffer := make(map[string]([]chan DatabaseModels.FoodItem))

	x.initiatePreparationMachines(&employeeList, capacityControl, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, preparationMachinesFoodItems, preparationMachinesFoodItemsBuffer)

	//May use for test
	// prepMachineGrill := DatabaseModels.PreparationMachine{Name: "Grill", Capacity: 12, NumberOfMachines: 3} // We have 3 grills with capacity 12
	// createPreparationMachineIngredientBuffer(prepMachineGrill, preparationMachinesIngredientsBuffer, capacityControl, preparationMachinesIngredients)
	// for i := 0; i < prepMachineGrill.NumberOfMachines; i++ {
	// 	go runPreparationMachineIngredient(prepMachineGrill, capacityControl, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, i, &employeeList)
	// }

	// prepMachineCuttingBoard := DatabaseModels.PreparationMachine{Name: "CuttingBoard", Capacity: 2, NumberOfMachines: 2} // We have 2 cutting boards with capacity 2
	// for i := 0; i < prepMachineCuttingBoard.NumberOfMachines; i++ {
	// 	go runPreparationMachineIngredient(prepMachineCuttingBoard, capacityControl, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, i, &employeeList)
	// }
	// createPreparationMachineIngredientBuffer(prepMachineCuttingBoard, preparationMachinesIngredientsBuffer, capacityControl, preparationMachinesIngredients)

	//May use for test
	// prepMachineBurgerAssemblyPoint := PreparationMachine{Name: "Burger Assembly Point", Capacity: 1, NumberOfMachines: 4} // We have 4 assembly points each with capacity 1
	// createPreparationMachineFoodItemBuffer(prepMachineBurgerAssemblyPoint, preparationMachinesFoodItemsBuffer, capacityControl, preparationMachinesFoodItems)

	// for i := 0; i < prepMachineBurgerAssemblyPoint.NumberOfMachines; i++ {
	// 	go runPreparationMachineFoodItem(prepMachineBurgerAssemblyPoint, capacityControl, preparationMachinesFoodItemsBuffer, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, preparationMachinesFoodItems, i, &employeeList)
	// }

	// prepMachineSanvisusMaker := PreparationMachine{Name: "SanvisusMaker", Capacity: 3, NumberOfMachines: 1}
	// createPreparationMachineFoodItemBuffer(prepMachineSanvisusMaker, preparationMachinesFoodItemsBuffer, capacityControl, preparationMachinesFoodItems)

	// for i := 0; i < prepMachineSanvisusMaker.NumberOfMachines; i++ {
	// 	go runPreparationMachineFoodItem(prepMachineSanvisusMaker, capacityControl, preparationMachinesFoodItemsBuffer, preparationMachinesIngredientsBuffer, preparationMachinesIngredients, preparationMachinesFoodItems, i, &employeeList)
	// }

	x.initiateIngredients(preparationMachinesIngredients)

	//May use for testing
	// ingredientTomato := Ingredient{Name: "Tomato", PreparationMachine: prepMachineCuttingBoard, PreparationTime: 2} // 2 seconds to cut the tomato
	// preparationMachinesIngredients[ingredientTomato.PreparationMachine.Name] = make(map[string]chan Ingredient)
	// addIngredientToPrepMachine(ingredientTomato, preparationMachinesIngredients)
	// ingredientPatty := Ingredient{Name: "Patty", PreparationMachine: prepMachineGrill, PreparationTime: 3}
	// preparationMachinesIngredients[ingredientPatty.PreparationMachine.Name] = make(map[string]chan Ingredient)
	// addIngredientToPrepMachine(ingredientPatty, preparationMachinesIngredients)

	x.initiateFoodItems(preparationMachinesFoodItems)

	//May use for testing
	// ingredients := []Ingredient{ingredientPatty, ingredientTomato}

	// foodItemBurger2 := FoodItem{Name: "Burger2", Ingredients: ingredients, PreparationMachine: prepMachineBurgerAssemblyPoint, PreparationTime: 4}
	// addFoodItemToPrepMachine(foodItemBurger2, preparationMachinesFoodItems)

	// foodItemBurger := FoodItem{Name: "Burger1", Ingredients: ingredients, PreparationMachine: prepMachineBurgerAssemblyPoint, PreparationTime: 2}
	// addFoodItemToPrepMachine(foodItemBurger, preparationMachinesFoodItems)

	// foodItemGrilledBurger := FoodItem{Name: "GrilledBurger", Ingredients: []Ingredient{}, FoodItems: []FoodItem{foodItemBurger}, PreparationMachine: prepMachineSanvisusMaker, PreparationTime: 5}
	// addFoodItemToPrepMachine(foodItemGrilledBurger, preparationMachinesFoodItems)

	go runOrderThread(preparationMachinesFoodItemsBuffer, capacityControl, preparationMachinesFoodItems, ordersSent, ordersReady)
	capacityControl["Order"] = make([]chan bool, 20)
	capacityControl["Order"][0] = make(chan bool, 20) // Can only send 20 orders at a time

	//Waiting for orders
	for {
		order := <-x.Orders
		sendOrder(order, capacityControl, ordersSent)
	}

	//May use for testing
	// order := Order{Name: "Order 1", FoodItems: []FoodItem{foodItemBurger, foodItemBurger2}}
	// order2 := Order{Name: "Order 2", FoodItems: []FoodItem{foodItemBurger2, foodItemGrilledBurger}}

	// time.Sleep(5 * time.Second) // leaving time for the setup to be ready
	// sendOrder(order, capacityControl, ordersSent)
	// sendOrder(order2, capacityControl, ordersSent)

}

//Problems:
// 1. The preparation machine capacity is substracted before actually starting to work. The substraction of the capacity should only happen [...]
// when all of the ingredients for the food item have been gathered.

// TODO:
// 1. Make a diagram for the processes
// 2. Refactor code
// 3. Make the preparation food machine work with both food items and ingredients.
