package tests

import (
	"github.com/universe-10th/identity/hashing"
	"testing"
)

func TestMultiHasherParsingH0Passwords(t *testing.T) {
	multi, h0, _ := MakeMultiHasherExampleInstances()
	h, _ := h0.Hash("foo$123")
	hashed := h0.Name() + ":" + h

	if err := multi.Validate("foo$123", hashed); err != nil {
		t.Errorf("Password-check using hasher 0 in multi hasher must succeed. Error received: %s\n", err)
	}
}

func TestMultiHasherParsingH1Passwords(t *testing.T) {
	multi, _, h1 := MakeMultiHasherExampleInstances()
	h, _ := h1.Hash("foo$123")
	hashed := h1.Name() + ":" + h

	if err := multi.Validate("foo$123", hashed); err != nil {
		t.Errorf("Password-check using hasher 1 in multi hasher must succeed. Error received: %s\n", err)
	}
}

func TestMultiHasherParsingUnregisteredEnginePassword(t *testing.T) {
	multi, _, _ := MakeMultiHasherExampleInstances()
	hu := DummyHasher(2)
	h, _ := hu.Hash("foo$123")
	hashed := hu.Name() + ":" + h

	if err := multi.Validate("foo$123", hashed); err != hashing.ErrUnregisteredEngine {
		t.Errorf("Password-check using unregistered hasher in multi hasher must fail with hashing.ErrUnregisteredEngine. Error returned instead: %s\n", err)
	}
}

func TestMultiHasherHashingByDefault(t *testing.T) {
	multi, h0, _ := MakeMultiHasherExampleInstances()
	h, _ := h0.Hash("foo$123")
	hashed := h0.Name() + ":" + h
	hashedDefault, _ := multi.Hash("foo$123")
	if hashed != hashedDefault {
		t.Errorf("Hashing with multi hasher should hash as if the h0 was selected, plus a prefix. Hashed instead: %s\n", hashedDefault)
	}
}
