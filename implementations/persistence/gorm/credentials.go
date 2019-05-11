package gorm

import (
	"github.com/luismasuelli/go-identity/interfaces"
	"github.com/luismasuelli/go-identity/types"
	"github.com/jinzhu/gorm"
)

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

func (user *User) HashingEngine() interfaces.PasswordHashingEngine {
	panic("not implemented")
}

func (user *User) CheckLogin(stage types.LoginStage) error {
	return nil
}