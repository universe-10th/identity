package bare

import "github.com/jinzhu/gorm"


func MigrateCredentials(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
