package bare

import (
	"github.com/jinzhu/gorm"
	"github.com/luismasuelli/go-identity/stub"
	"github.com/luismasuelli/go-identity/support/types"
	"errors"
)


var UserIsDeleted = errors.New("user is deleted")


/**
 * User bare implementation for GORM engine.
 */
type User struct {
	gorm.Model
	Username string `gorm:"size:30;not null;unique"`
	Password string `gorm:"size:100;not null"`
}

func (*User) IdentificationField() string {
	return "username"
}

func (*User) IdentificationIsCaseSensitive() bool {
	return false
}

func (user *User) SetHashedPassword(password string) {
	user.Password = password
}

func (user *User) ClearPassword() {
	user.SetHashedPassword("")
}

func (user *User) HashedPassword() string {
	return user.Password
}

func (user *User) HashingEngine() stub.PasswordHashingEngine {
	panic("not implemented")
}

func (user *User) CheckLogin(stage types.LoginStage) error {
	if stage == types.BeforePasswordCheck {
		if user.DeletedAt != nil {
			return UserIsDeleted
		}
	}
	return nil
}

