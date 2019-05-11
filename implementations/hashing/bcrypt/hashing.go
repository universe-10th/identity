package bcrypt

import "golang.org/x/crypto/bcrypt"

type BCryptHashingEngine struct {
	cost int
}

func (bcryptHashingEngine *BCryptHashingEngine) Hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcryptHashingEngine.cost)
	if err != nil {
		return "", err
	} else {
		return string(result), err
	}
}

func (bcryptHashingEngine *BCryptHashingEngine) Validate(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (bcryptHashingEngine *BCryptHashingEngine) Name() string {
	return "bcrypt"
}