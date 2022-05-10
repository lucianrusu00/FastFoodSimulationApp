package Controllers

import (
	"BusinessLogicModels"
	"net/http"

	"github.com/gin-gonic/gin"

	"BusinessLogicServices"
)

type SetupController struct {
	SetupService BusinessLogicServices.SetupService
	Router       *gin.Engine
}

func (x SetupController) AddPrepMachine(c *gin.Context) {

	var prepMachine BusinessLogicModels.PreparationMachineDTO
	c.BindJSON(&prepMachine)

	x.SetupService.AddPrepMachine(prepMachine)

	// fmt.Printf("The name of the prep machine is: %s \n", c.Param("name"))

	// name := c.PostForm("name")
	// fmt.Printf("The name is: %s \n", name)

	// capacity, _ := strconv.Atoi(c.PostForm("capacity"))
	// fmt.Printf("The capacity is: %d \n", capacity)

	// x.SetupService.AddPrepMachine(c.PostForm("name"), uint(capacity))

	c.JSON(http.StatusOK, gin.H{"response": "prep machine added"})

}

func (x SetupController) AddIngredients(c *gin.Context) {

	var ingredient BusinessLogicModels.IngredientDTO
	c.BindJSON(&ingredient)

	x.SetupService.AddIngredient(ingredient)

	// fmt.Printf("The name of the ingredient is: ")
	// fmt.Println(ingredient.Name)

	// fmt.Printf("The preparation time is: ")
	// fmt.Println(ingredient.PreparationTime)

	// fmt.Printf("It has the prepMachine: %s\n", ingredient.PreparationMachineName)

}

func (x SetupController) AddFoodItem(c *gin.Context) {
	var foodItemDTO BusinessLogicModels.FoodItemDTO
	c.BindJSON(&foodItemDTO)

	x.SetupService.AddFoodItem(foodItemDTO)
}

func (x SetupController) AddEmployee(c *gin.Context) {
	var employeeDTO BusinessLogicModels.EmployeeDTO
	c.BindJSON(&employeeDTO)

	x.SetupService.AddEmployee(employeeDTO)

}

func (x SetupController) AddOrder(c *gin.Context) {
	var orderDTO BusinessLogicModels.OrderDTO
	c.BindJSON(&orderDTO)

	x.SetupService.AddOrder(orderDTO)

}

func (x SetupController) StartSimulation(c *gin.Context) {
	x.SetupService.StartSimulation()
}

func (x SetupController) Run() {

	//x.Router.GET("/albumByID/:id")
	x.Router.POST("/addPreparationMachine", x.AddPrepMachine)
	x.Router.POST("/addIngredient", x.AddIngredients)
	x.Router.POST("/addFoodItem", x.AddFoodItem)
	x.Router.POST("/addEmployee", x.AddEmployee)
	x.Router.POST("/addOrder", x.AddOrder)
	x.Router.GET("/startSimulation", x.StartSimulation)

}
