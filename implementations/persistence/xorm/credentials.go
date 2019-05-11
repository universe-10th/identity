package xorm

import (
	"github.com/luismasuelli/go-identity/interfaces"
	"github.com/luismasuelli/go-identity/types"
	"time"
)

type User struct {
	Id uint `xorm:"pk"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Username string `xorm:"varchar(30) not null unique"`
	Password string `xorm:"varchar(100) not null"`
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