package xorm

import (
	"github.com/go-xorm/xorm"
	"github.com/luismasuelli/go-identity/interfaces"
	"fmt"
	"errors"
)


type XORMSource struct {
	db *xorm.Engine
}


var NotFound = errors.New("record not found")


func (xormSource XORMSource) Lookup(resultHolder interfaces.Credential, identification string) error {
	caseSensitive := resultHolder.IdentificationIsCaseSensitive()
	query := ""
	if caseSensitive {
		query = fmt.Sprintf("%s = ?", resultHolder.IdentificationField())
	} else {
		query = fmt.Sprintf("UPPER(%s) = UPPER(?)", resultHolder.IdentificationField())
	}
	if got, err := xormSource.db.Where(query, identification).Get(resultHolder); err != nil {
		return err
	} else if !got {
		return NotFound
	} else {
		return nil
	}
}


func NewSource(db *xorm.Engine) interfaces.Source {
	return &XORMSource{db}
}