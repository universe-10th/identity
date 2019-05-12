package xorm

import "github.com/luismasuelli/go-identity/interfaces"

/**
 * Scope implementation for GORM engine.
 */
type Scope struct {
	Id uint `xorm:"pk"`
	Code string `xorm:"varchar(30) not null unique"`
	Label string `xorm:"varchar(50) not null"`
	HelpText string `xorm:"longtext not null"`
}

/**
 * A scope holder is an interface (intended for the
 *   credentials) returning associated instances of
 *   Scope for it.
 */
type ScopeHolder interface {
	Scopes() map[string]Scope
}

func (scope *Scope) Key() string {
	return scope.Code
}

func (scope *Scope) Name() string {
	return scope.Label
}

func (scope *Scope) Description() string {
	return scope.HelpText
}

func (scope *Scope) SatisfiedBy(credential interfaces.Credential) bool {
	if holder, isHolder := credential.(ScopeHolder); isHolder {
		_, ok := holder.Scopes()[scope.Key()]
		return ok
	} else {
		return false
	}
}