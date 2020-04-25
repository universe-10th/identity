package scoped

import (
	"github.com/universe-10th/identity/credentials"
	scoped2 "github.com/universe-10th/identity/credentials/traits/scoped"
)

// This requirement checks whether a
// credential is a scoped one satisfying
// a "scope tree" requirements.
type Scoped struct {
	tree ScopeTreeNode
}

// Tests whether the credential is superuser or not.
func (scoped *Scoped) SatisfiedBy(credential credentials.Credential) bool {
	if scopedCredential, ok := credential.(scoped2.Scoped); ok && scoped.tree.SatisfiedBy(scopedCredential) {
		return true
	} else {
		return false
	}
}

// Creates a new Scoped requirement checking all the
// given Scopes or tree nodes (i.e. Any() or All())
// in an "any" way. This means: alternative scopes
// or nodes will be tried and one of them must be
// entirely satisfied. Intended for main user/player
// interaction with the server.
func RequireScopesAmong(scopes ...ScopeSpec) *Scoped {
	return &Scoped{Any(scopes...)}
}
