package gorm

import (
	"github.com/universe-10th/identity"
)

// This is an implementation of the Scope interface
// consisting on database-stored instances. There
// will no more satisfaction checks than the scope
// belonging to a scope map in the related credential,
// who must also implement ModelBackedScopeHolder
// interface and contain this scope in the map for
// the satisfaction question to be answered as true.
type ModelBackedScope struct {
	Id       uint   `gorm:"primary_key"`
	Code     string `gorm:"size:30;not null;unique"`
	Label    string `gorm:"size:50;not null"`
	HelpText string `gorm:"type:text;not null"`
}

// Credential implementations must implement this
// interface if they expect ModelBackedScope instances
// return true in the satisfaction question.
type ModelBackedScopeHolder interface {
	GetScopes(forceRefresh bool) map[string]*ModelBackedScope
}

// The scope's key is stored in its "key" field.
func (scope *ModelBackedScope) Key() string {
	return scope.Code
}

// The scope's name is stored in its "name" field.
func (scope *ModelBackedScope) Name() string {
	return scope.Label
}

// The scope's description is stored in its "description"
// field.
func (scope *ModelBackedScope) Description() string {
	return scope.HelpText
}

// Satisfaction check against a credential implies the
// credential implements ModelBackedScopeHolder and also
// contains the current scope among their scopes.
func (scope *ModelBackedScope) SatisfiedBy(credential identity.Credential) bool {
	if holder, isHolder := credential.(ModelBackedScopeHolder); isHolder {
		_, ok := holder.GetScopes(false)[scope.Key()]
		return ok
	} else {
		return false
	}
}
