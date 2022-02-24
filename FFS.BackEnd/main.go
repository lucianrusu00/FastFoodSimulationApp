package main

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"BusinessLogicModels"
	"BusinessLogicServices"

	"Controllers"

	"DatabaseModels"
)

// albums slice to seed record album data.
var albums = []BusinessLogicModels.AlbumDTO{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Where("1 = 1").Unscoped().Delete(&DatabaseModels.Album{}) // Deleting all the albums

	db.AutoMigrate(&DatabaseModels.Album{})                                                               // creating the migration for the album table
	db.Create(&DatabaseModels.Album{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99}) // adding an album to the database

	var album DatabaseModels.Album
	db.First(&album, "ID = ?", "1") // getting the album that has the id = 1
	fmt.Printf("The title of the album is: ")
	fmt.Println(album.Title)

	//db.Unscoped().Delete(&album, "ID = ?", "1") // deleting permanently the album with id = 1, (if you want to soft delete it don't use unscoped keyword)

	fmt.Printf("Hello mate\n")
	router := gin.Default()

	albumService := BusinessLogicServices.AlbumService{db}
	albumController := Controllers.AlbumController{albumService, router}
	albumController.Run()
	//router.GET("/albums", getAlbums)

	router.Run("localhost:8080")

}
