package Controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"BusinessLogicModels"
	"BusinessLogicServices"
)

type AlbumController struct {
	AlbumService BusinessLogicServices.AlbumService
	Router       *gin.Engine
}

func (x AlbumController) getAlbumWithID(c *gin.Context) {
	id := c.Param("id")

	var albumDTO BusinessLogicModels.AlbumDTO
	albumDTO = x.AlbumService.GetAlbumWithID(id)

	c.IndentedJSON(http.StatusOK, albumDTO)

}

func (x AlbumController) Run() {

	x.Router.GET("/albumByID/:id", x.getAlbumWithID)

}
