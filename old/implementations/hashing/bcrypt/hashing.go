package bcrypt

import (
	"github.com/universe-10th/identity"
	"golang.org/x/crypto/bcrypt"
)

// BCrypt hashing facade.
type bcryptHashingEngine struct {
	cost int
}

func (bcryptHashingEngine *bcryptHashingEngine) Hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcryptHashingEngine.cost)
	if err != nil {
		return "", err
	} else {
		return string(result), err
	}
}

func (bcryptHashingEngine *bcryptHashingEngine) Validate(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (bcryptHashingEngine *bcryptHashingEngine) Name() string {
	return "bcrypt"
}

func New(cost int) identity.PasswordHashingEngine {
	return &bcryptHashingEngine{cost}
}

var Default = New(10)
