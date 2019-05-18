package multiple

import (
	"strings"
	"errors"
	"github.com/luismasuelli/go-identity/stub"
)

/**
 * Wraps many password hashing engines in one and allows
 *   using all of them (with a default for one to write).
 */
type MultiplePasswordHashingEngine struct {
	defaultEngine string
	registeredEngines map[string]stub.PasswordHashingEngine
}

var InvalidHash = errors.New("invalid hash string")
var UnregisteredEngine = errors.New("unregistered engine")

func (multiplePasswordHashingEngine *MultiplePasswordHashingEngine) Hash(password string) (string, error) {
	engine := multiplePasswordHashingEngine.registeredEngines[multiplePasswordHashingEngine.defaultEngine]
	if hashed, err := engine.Hash(password); err != nil {
		return "", err
	} else {
		return engine.Name() + ":" + hashed, nil
	}
}

func (multiplePasswordHashingEngine *MultiplePasswordHashingEngine) Validate(password string, hash string) error {
	parts := strings.SplitN(hash, ":", 2)
	// By default, if the password is not <key>:<hash>, we
	//   just take the default engine. If the prefix was
	//   specified, then a specific engine will be used to
	//   validate that password against the hash.
	engineKey := multiplePasswordHashingEngine.defaultEngine
	switch len(parts) {
	case 1:
		engineKey = multiplePasswordHashingEngine.defaultEngine
		hash = parts[0]
	case 2:
		engineKey = parts[0]
		hash = parts[1]
	default:
		return InvalidHash
	}
	// Given the engine key, password, and hash, calculate the validation.
	if engine, ok := multiplePasswordHashingEngine.registeredEngines[engineKey]; !ok {
		return UnregisteredEngine
	} else {
		return engine.Validate(password, hash)
	}
}

func (multiplePasswordHashingEngine *MultiplePasswordHashingEngine) Name() string {
	panic("multiple password hashing engines have no name")
}


func NewWithDefault(defaultEngine stub.PasswordHashingEngine, engines ...stub.PasswordHashingEngine) stub.PasswordHashingEngine {
	if len(engines) == 0 {
		panic("no engines were specified")
	}

	mphe := &MultiplePasswordHashingEngine{}

	for _, engine := range engines {
		if engine == nil {
			panic("nil engine instance")
		}
		name := engine.Name()
		if name == "" {
			panic("empty engine name (this is an engine implementation issue)")
		}
		if _, ok := mphe.registeredEngines[name]; ok {
			panic("duplicate engine name (two hashing engines of the same type were registered)")
		} else {
			mphe.registeredEngines[name] = engine
		}
	}

	if defaultEngine == nil {
		mphe.defaultEngine = engines[0].Name()
		return mphe
	} else if engine, _ := mphe.registeredEngines[defaultEngine.Name()]; engine != defaultEngine {
		panic("default engine is specified but not present among hashing engines")
	} else {
		mphe.defaultEngine = defaultEngine.Name()
		return mphe
	}
}


func New(engines ...stub.PasswordHashingEngine) stub.PasswordHashingEngine {
	return NewWithDefault(nil, engines...)
}