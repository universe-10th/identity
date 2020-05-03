package tests

import "errors"

type DummyHasher uint

func (DummyHasher) Name() string { return "dummy" }
func (DummyHasher) Hash(password string) (string, error) {
	return "Hashed[" + password + "]", nil
}
func (h DummyHasher) Validate(password string, hash string) error {
	if h, _ := h.Hash(password); h == hash {
		return nil
	} else {
		return errors.New("bad password")
	}
}
