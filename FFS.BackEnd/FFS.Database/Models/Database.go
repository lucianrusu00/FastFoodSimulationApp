package DatabaseModels

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func (x *Database) Run() {
	var err error
	x.Db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//x.Db.AutoMigrate(&Album{}) // creating the migration for the album table

	x.Db.AutoMigrate(&PreparationMachine{})
	x.Db.AutoMigrate(&PreparationMachineFoodItem{})
	x.Db.AutoMigrate(&Ingredient{})
	x.Db.AutoMigrate(&FoodItem{})
	x.Db.AutoMigrate(&Employee{})
	x.Db.AutoMigrate(&Order{})

	// x.Db.Unscoped().Where("1 = 1").Delete(&PreparationMachine{})
	// x.Db.Unscoped().Where("1 = 1").Delete(&PreparationMachineFoodItem{})
	// x.Db.Unscoped().Where("1 = 1").Delete(&Ingredient{})
	// x.Db.Unscoped().Where("1 = 1").Delete(&FoodItem{})
	// x.Db.Unscoped().Where("1 = 1").Delete(&Employee{})
	x.Db.Unscoped().Where("1 = 1").Delete(&Order{})

	// fmt.Println("creating a preparation machine")
	// x.Db.Create(&PreparationMachine{Name: "nume", Capacity: uint(12)})
	// fmt.Println("created the preparation machine")

}
