package tests

import "fmt"

type DummyScope uint8

func (scope DummyScope) Key() string {
	return fmt.Sprintf("scope%d", scope)
}

func (scope DummyScope) Name() string {
	return fmt.Sprintf("Scope %d", scope)
}

func (scope DummyScope) Description() string {
	return fmt.Sprintf("Sample scope %d", scope)
}
