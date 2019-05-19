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
// the underlying model's table, using appropriate model's field and
// case sensitivity inside a GORM connection.
func (gormSource *gormSource) ByIdentification(resultHolder stub.Credential, identification string) error {
	caseSensitive := resultHolder.IdentificationIsCaseSensitive()
	query := ""
	if caseSensitive {
		query = fmt.Sprintf("%s = ?", resultHolder.IdentificationField())
	} else {
		query = fmt.Sprintf("UPPER(%s) = UPPER(?)", resultHolder.IdentificationField())
	}
	return gormSource.db.Where(query, identification).First(resultHolder).Error
}


// Instantiates a GORM-compatible lookup source for a particular db
// connection given as argument.
func NewSource(db *gorm.DB) stub.Source {
	return &gormSource{db}
}
