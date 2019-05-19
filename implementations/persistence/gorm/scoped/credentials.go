package scoped

import "github.com/universe-10th/identity/implementations/persistence/gorm/bare"


type User struct {
	bare.User
	SuperUser   bool                         `gorm:"not null"`
	Scopes      []*ModelBackedScope          `gorm:"many2many:user_scopes"`
	scopesMap   map[string]*ModelBackedScope `gorm:"-"`
}


func (user *User) IsSuperUser() bool {
	return user.SuperUser
}


func (user *User) GetScopes(forceRefresh bool) map[string]*ModelBackedScope {
	if user.scopesMap == nil || forceRefresh {
		user.scopesMap = map[string]*ModelBackedScope{}
		for _, scope := range user.Scopes {
			user.scopesMap[scope.Key()] = scope
		}
	}
	return user.scopesMap
}