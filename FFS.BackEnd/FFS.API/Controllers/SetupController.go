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

/*

You have a list of N+1 integers between 1 and N.
You know there's at least one duplicate, but there might be more.
Print out a number that appears in the list more than once.

For example:
N=3        list: 3, 1, 1, 3 -> result 1 or 3
N=5        list: 5 2 3 4 5 1 -> result 5


  sumMax = 1+..+N + N

  sumList = arr[0] + ... + arr[arr.size()-1]

  sumMax < sumList

  N/2 -> N

  sumMaxLeft = N/2 + N/2+1 + ... + N

  sumList = arr[0] + ... + arrr[N-1] // only if arr[i] >= N/2

int getDuplicateRec(int arr[], sumMax, sumList, NMin, NMax){
  int sumListLeft

  for(int i = 0; i < N+1; i++){
    if(arr[i] >= NMin && arr[i] <= Nmax)
  }
}


int getDuplicate(int arr[]){

  int sumMax = N*(N+1)/2 + N
  int sumList = 0

  for(int i = 0; i < N+1; i++){
    sumList += arr[i]
  }

  if(sumMax < sumList){
    int sumListLeft
    int sumMaxLeft = sumMax - ((N\2-1)(N\2)/2)
    for(int i = 0; i < N+1; i++)
      if(arr[i] >= N/2)
        sumListLeft += arr[i]
  }
}


//without added memory
int getDuplicate(int arr[]){
  for(int i = 0; i < arr.size()-1; i++){
    for(int j = i+1; j<arr.size(); j++){
      if(arr[i] == arr[j])
        return arr[i]
    }
  }

  return -1
}

int getDuplicate(int arr[]){
  arr.sort()

	for(int i = 0; i < arr.size()-1; i++){
    if(arr[i] == arr[i+1])
      return arr[i]
  }

  return -1
}






*/
