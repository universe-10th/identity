package identity


// Hashing engines are facades of regularly (already
// implemented) algorithms like bcrypt.
type HashingEngine interface {
	// Creates a hash, ready to be stored.
	Hash(password string)                  (string, error)
	// Validates a password against a hash.
	Validate(password string, hash string) error
}