package scoped

import (
	"github.com/universe-10th/identity/stub"
)

/**
 * Scope implementation for GORM engine.
 */
type ModelBackedScope struct {
	Id uint `gorm:"primary_key"`
	Code string `gorm:"size:30;not null;unique"`
	Label string `gorm:"size:50;not null"`
	HelpText string `gorm:"type:text;not null"`
}

/**
 * A scope holder is an interface (intended for the
 *   credentials) returning associated instances of
 *   ModelBackedScope for it.
 */
type ModelBackedScopeHolder interface {
	GetScopes(forceRefresh bool) map[string]*ModelBackedScope
}

func (scope *ModelBackedScope) Key() string {
	return scope.Code
}

func (scope *ModelBackedScope) Name() string {
	return scope.Label
}

func (scope *ModelBackedScope) Description() string {
	return scope.HelpText
}

func (scope *ModelBackedScope) SatisfiedBy(credential stub.Credential) bool {
	if holder, isHolder := credential.(ModelBackedScopeHolder); isHolder {
		_, ok := holder.GetScopes(false)[scope.Key()]
		return ok
	} else {
		return false
	}
}
