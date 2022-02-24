package BusinessLogicServices

import (
	"BusinessLogicModels"

	"DatabaseModels"

	"gorm.io/gorm"
)

type AlbumService struct {
	Db *gorm.DB
}

func (x AlbumService) GetAlbumWithID(id string) BusinessLogicModels.AlbumDTO {
	var album DatabaseModels.Album
	x.Db.First(&album, "ID = ?", id)

	albumDTO := BusinessLogicModels.AlbumDTO{ID: album.ID, Title: album.Title, Artist: album.Artist, Price: album.Price}

	return albumDTO
}
