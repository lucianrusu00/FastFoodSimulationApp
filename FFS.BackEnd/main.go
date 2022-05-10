package main

import (
	"BusinessLogicUtilities"
	"fmt"

	"github.com/gin-gonic/gin"

	"BusinessLogicServices"

	"Controllers"

	"DatabaseModels"
)

// introduction
// motivation
// contextual background
// overleaf.com - for latex

// rrelated work, methodology, fanancial figures

func main() {

	fmt.Printf("Starting the server\n")

	database := new(DatabaseModels.Database)

	database.Run()

	//database.Db.Create(&DatabaseModels.PreparationMachine{Name: "nume", Capacity: uint(28)})

	//database.Db.Create(&DatabaseModels.Album{ID: "3", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99}) // adding an album to the database

	// var album DatabaseModels.Album
	// database.Db.First(&album, "ID = ?", "1") // getting the album that has the id = 1
	// fmt.Printf("The title of the album is: ")
	// fmt.Println(album.Title)

	//db.Unscoped().Delete(&album, "ID = ?", "1") // deleting permanently the album with id = 1, (if you want to soft delete it don't use unscoped keyword)

	//database.Db.Unscoped().Where("1 = 1").Delete(&DatabaseModels.PreparationMachine{}) // delete all elements from the preparation_machine table

	fmt.Printf("Hello mate\n")
	router := gin.Default()

	albumService := BusinessLogicServices.AlbumService{database.Db}
	albumController := Controllers.AlbumController{albumService, router}
	albumController.Run()

	orders := make(chan DatabaseModels.Order, 100)
	setupSimulation := BusinessLogicUtilities.Simulate{Db: database.Db, Orders: orders}

	setupService := BusinessLogicServices.SetupService{Db: database.Db, Simulate: setupSimulation}
	setupController := Controllers.SetupController{SetupService: setupService, Router: router}
	setupController.Run()

	router.Run("localhost:8080")

}
