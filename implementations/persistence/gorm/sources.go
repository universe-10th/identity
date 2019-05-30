package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/universe-10th/identity/stub"
	"fmt"
)


type gormSource struct {
	db *gorm.DB
}


// Returns a GORM-compatible lookup: it will perform a lookup against
// the underlying model's table, using model's identification field and
// case sensitivity inside a GORM connection.
func (gormSource *gormSource) ByIdentification(resultHolder stub.Credential, identification interface{}) error {
	caseSensitive := resultHolder.IdentificationIsCaseSensitive()
	query := ""
	if caseSensitive {
		query = fmt.Sprintf("%s = ?", resultHolder.IdentificationField())
	} else {
		query = fmt.Sprintf("UPPER(%s) = UPPER(?)", resultHolder.IdentificationField())
	}
	return gormSource.db.Where(query, identification).First(resultHolder).Error
}


// Returns a GORM-compatible lookup: it will perform a lookup against
// the underlying model's table, using model's primary key field inside
// a GORM connection.
func (gormSource *gormSource) ByPrimaryKey(resultHolder stub.Credential, pk interface{}) error {
	return gormSource.db.Where(resultHolder.PrimaryKeyField() + " = ?", pk).First(resultHolder).Error
}


// Saves a credential (either updating or creating) - returns its associated error
func (gormSource *gormSource) Save(credential stub.Credential) error {
	return gormSource.db.Save(credential).Error
}


// Deletes a credential (either hard or soft) - returns its associated error
func (gormSource *gormSource) Delete(credential stub.Credential) error {
	return gormSource.db.Delete(credential).Error
}


// Instantiates a GORM-compatible lookup source for a particular db
// connection given as argument.
func NewSource(db *gorm.DB) stub.Source {
	return &gormSource{db}
}
