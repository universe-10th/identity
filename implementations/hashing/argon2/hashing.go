package argon2

import (
	"errors"
	"github.com/universe-10th/identity/support/utils"
	"golang.org/x/crypto/argon2"
	"encoding/base64"
	"strings"
	"github.com/universe-10th/identity/stub"
	"runtime"
)


var InvalidHash = errors.New("invalid hash string")
var PasswordMismatch = errors.New("password mismatch")


// Argon2 hashing facade.
type argon2HashingEngine struct {
	time, memory, keyLen uint32
	threads, saltLength uint8
}


func (argon2HashingEngine *argon2HashingEngine) Hash(password string) (string, error) {
	salt := utils.Salt(int(argon2HashingEngine.saltLength))
	result := argon2.IDKey(
		[]byte(password), []byte(salt),
		argon2HashingEngine.time,
		argon2HashingEngine.memory,
		argon2HashingEngine.threads,
		argon2HashingEngine.keyLen,
    )
	return salt + "$" + base64.StdEncoding.EncodeToString(result), nil
}


func (argon2HashingEngine *argon2HashingEngine) Validate(password string, hash string) error {
	parts := strings.SplitN(hash, "$", 2)
	if len(parts) != 2 {
		return InvalidHash
	}

	salt := parts[0]
	result := argon2.IDKey(
		[]byte(password), []byte(salt),
		argon2HashingEngine.time,
		argon2HashingEngine.memory,
		argon2HashingEngine.threads,
		argon2HashingEngine.keyLen,
	)
	if base64.StdEncoding.EncodeToString(result) != parts[1] {
		return PasswordMismatch
	} else {
		return nil
	}
}


func (argon2HashingEngine *argon2HashingEngine) Name() string {
	return "argon2"
}


func New(saltLength uint8, time, memory uint32, threads uint8, keyLen uint32) stub.PasswordHashingEngine {
	return &argon2HashingEngine{time, memory, keyLen, threads, saltLength}
}


var Default = New(8,1, 2 << 16, uint8(runtime.NumCPU()), 32)