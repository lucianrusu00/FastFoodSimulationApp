module main

go 1.15

require (
	BusinessLogicModels v1.0.0
	DatabaseModels v1.0.0
    Controllers v1.0.0
	BusinessLogicServices v1.0.0
	github.com/gin-gonic/gin v1.7.7
	github.com/mattn/go-sqlite3 v1.14.11 // indirect
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.1
)

replace Controllers v1.0.0 => ./FFS.API/Controllers

replace BusinessLogicModels v1.0.0 => ./FFS.BusinessLogic/Models

replace DatabaseModels v1.0.0 => ./FFS.Database/Models

replace BusinessLogicServices v1.0.0 => ./FFS.BusinessLogic/Services
