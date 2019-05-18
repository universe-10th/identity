package scoped

import "github.com/jinzhu/gorm"


func MigrateAll(db *gorm.DB) {
	MigrateScopes(db)
	db.AutoMigrate(&User{})
}


func MigrateScopes(db *gorm.DB) {
	db.AutoMigrate(&ModelBackedScope{})
}