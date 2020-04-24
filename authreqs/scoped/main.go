package scoped

import (
	"github.com/universe-10th/identity"
	scoped2 "github.com/universe-10th/identity/traits/credential/scoped"
)

// This requirement checks whether a
// credential is a scoped one satisfying
// a "scope tree" requirements.
type Scoped struct {
	tree ScopeTreeNode
}

// Tests whether the credential is superuser or not.
func (scoped Scoped) SatisfiedBy(credential identity.Credential) bool {
	if scopedCredential, ok := credential.(scoped2.Scoped); ok && scoped.tree.SatisfiedBy(scopedCredential) {
		return true
	} else {
		return false
	}
}

