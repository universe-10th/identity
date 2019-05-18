package scoped

import "github.com/jinzhu/gorm"

func MigrateCredentials(db *gorm.DB) {
	db.AutoMigrate(&User{})
}


func MigrateScopes(db *gorm.DB) {
	db.AutoMigrate(&ModelBackedScope{})
}