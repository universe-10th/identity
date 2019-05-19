package gorm

import (
	"github.com/universe-10th/identity/stub"
	"github.com/universe-10th/identity/support/types"
	"github.com/jinzhu/gorm"
	"errors"
)


// Tells whether a user is deleted (when it tries to login).
var UserIsDeleted = errors.New("user is deleted")


// A default implementation of users, which will also hold
// scopes and a superuser flag. It implements Credential
// (to be able to log in), WithSuperUserFlag (to be able
// to become superuser), and ModelBackedScopeHolder (so
// it can be recognized by the default implementation of
// scope: ModelBackedScope).
//
// This class is still abstract. You have, at least, to
// compose it inside another class and implement this
// method: HashingEngine() - please return a fixed
// instance.
type User struct {
	gorm.Model
	Username  string                       `gorm:"size:30;not null;unique"`
	Password  string                       `gorm:"size:100;not null"`
	Superuser bool                         `gorm:"not null"`
	Scopes    []*ModelBackedScope          `gorm:"many2many:user_scopes"`
	scopesMap map[string]*ModelBackedScope `gorm:"-"`
}


// Gets its primary key from the ID field.
func(user *User) PrimaryKey() interface{} {
	return user.ID
}


// Its serialized primary key will be its ID field.
func (*User) PrimaryKeyField() string {
	return "id"
}


// Its identification field is "username".
func (*User) IdentificationField() string {
	return "username"
}


// Its identification lookup is case insensitive.
func (*User) IdentificationIsCaseSensitive() bool {
	return false
}


// The hashed password will be set in the "password" field.
func (user *User) SetHashedPassword(password string) {
	user.Password = password
}


// Clearing the password will put an empty string in it.
func (user *User) ClearPassword() {
	user.SetHashedPassword("")
}


// The current hashed password will be retrieved from the
// "password" field.
func (user *User) HashedPassword() string {
	return user.Password
}


// This method is not implemented. You have to compose this
// type into a new one and implement this method: ensure a
// fixed instance is returned.
func (user *User) HashingEngine() stub.PasswordHashingEngine {
	panic("not implemented")
}


// The only login check this type performs, is the check of
// deletion (having a deletion date).
func (user *User) CheckLogin(stage types.LoginStage) error {
	if stage == types.BeforePasswordCheck {
		if user.DeletedAt != nil {
			return UserIsDeleted
		}
	}
	return nil
}


// Checking if it is superuser is done by a check on the
// "superuser" field.
func (user *User) IsSuperUser() bool {
	return user.Superuser
}


// Getting the (model backed) scopes is done by caching a
// map. Please ensure that while you load a user from the
// database you also preload the scopes and/or call the
// db.Model(obj).Related("Scopes") method.
func (user *User) GetScopes(forceRefresh bool) map[string]*ModelBackedScope {
	if user.scopesMap == nil || forceRefresh {
		user.scopesMap = map[string]*ModelBackedScope{}
		for _, scope := range user.Scopes {
			user.scopesMap[scope.Key()] = scope
		}
	}
	return user.scopesMap
}