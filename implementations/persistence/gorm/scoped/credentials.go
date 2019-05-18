package scoped

import "github.com/luismasuelli/go-identity/implementations/persistence/gorm/bare"


type User struct {
	bare.User
	IsSuperuser bool                         `gorm:"not null"`
	Scopes      []*ModelBackedScope          `gorm:"many2many:user_scopes"`
	scopesMap   map[string]*ModelBackedScope `gorm:"-"`
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