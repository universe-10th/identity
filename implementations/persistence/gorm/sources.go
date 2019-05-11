package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/luismasuelli/go-identity/interfaces"
	"fmt"
)


/**
 * Source implementation for GORM engine.
 */
type GORMSource struct {
	db *gorm.DB
}


func (gormSource GORMSource) Lookup(resultHolder interfaces.Credential, identification string) error {
	caseSensitive := resultHolder.IdentificationIsCaseSensitive()
	query := ""
	if caseSensitive {
		query = fmt.Sprintf("%s = ?", resultHolder.IdentificationField())
	} else {
		query = fmt.Sprintf("UPPER(%s) = UPPER(?)", resultHolder.IdentificationField())
	}
	return gormSource.db.Where(query, identification).Find(resultHolder).Error
}


func NewSource(db *gorm.DB) interfaces.Source {
	return &GORMSource{db}
}
