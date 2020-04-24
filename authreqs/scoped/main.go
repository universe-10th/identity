package scoped

import (
	"github.com/universe-10th/identity/traits/credential/scoped"
	"errors"
)

// This requirement contract type checks that
// a particular scoped [credential] satisfies
// it, or not.
type ScopeTreeNode interface {
	SatisfiedBy(scoped scoped.Scoped) bool
}

// This requirement checks satisfaction from
// a scoped one by checking if inner scope's
// key matches any of the keys among the ones
// in the scoped [credential].
type ScopeTreeLeaf struct {
	scope scoped.Scope
}

// Expects a single requirement.
func (leaf *ScopeTreeLeaf) SatisfiedBy(scoped scoped.Scoped) bool {
	_, ok := scoped.Scopes()[leaf.scope.Key()]
	return ok
}

// This scope is, actually, a list of scopes.
// It will involve requiring ALL its children
// requirements.
type AllScopesTreeNode []ScopeTreeNode

// Expects all of the requirements in the list.
func (node AllScopesTreeNode) SatisfiedBy(scoped scoped.Scoped) bool {
	for _, child := range node {
		if !child.SatisfiedBy(scoped) {
			return false
		}
	}
	return true
}

// This scope is, actually, a list of scopes.
// It will involve requiring ANY of its children
// requirements.
type AnyScopeTreeNode []ScopeTreeNode

// Expects all of the requirements in the list.
func (node AnyScopeTreeNode) SatisfiedBy(scoped scoped.Scoped) bool {
	for _, requirement := range node {
		if requirement.SatisfiedBy(scoped) {
			return true
		}
	}
	return false
}

type ScopeSpec interface{}

var ErrInvalidSpec = errors.New("input must be either a traits/credential/scoped.Scope or a ScopeTreeNode type")

// Convenience function to make an AllScopesTreeNode
// from Scope instances or other ScopeTreeNode instances.
func All(scopes ...ScopeSpec) AllScopesTreeNode {
	result := make(AllScopesTreeNode, len(scopes))
	for index, spec := range scopes {
		switch v := spec.(type) {
		case scoped.Scope:
			result[index] = &ScopeTreeLeaf{v}
		case ScopeTreeNode:
			result[index] = v
		default:
			panic(ErrInvalidSpec)
		}
	}
	return result
}

// Convenience function to make an AnyScopesTreeNode
// from Scope instances or other ScopeTreeNode instances.
func Any(scopes ...ScopeSpec) AnyScopeTreeNode {
	result := make(AnyScopeTreeNode, len(scopes))
	for index, spec := range scopes {
		switch v := spec.(type) {
		case scoped.Scope:
			result[index] = &ScopeTreeLeaf{v}
		case ScopeTreeNode:
			result[index] = v
		default:
			panic(ErrInvalidSpec)
		}
	}
	return result
}
