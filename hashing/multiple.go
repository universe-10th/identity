package hashing

import (
	"errors"
	"strings"
)

// Wraps many password hashing engines in one and allows
// using all of them (with a default for one to write).
// These hashers are considered beforehand under the
// possibility of changing the default hasher some day.
type MultipleHashingEngine struct {
	defaultEngine     string
	registeredEngines map[string]HashingEngine
}

var InvalidHash = errors.New("invalid hash string")
var UnregisteredEngine = errors.New("unregistered engine")

// Creates a hash using the default hasher engine.
func (multipleHashingEngine *MultipleHashingEngine) Hash(password string) (string, error) {
	engine := multipleHashingEngine.registeredEngines[multipleHashingEngine.defaultEngine]
	if hashed, err := engine.Hash(password); err != nil {
		return "", err
	} else {
		return hashed, nil
	}
}

// Validates a hash using whatever hasher is matched.
func (multipleHashingEngine *MultipleHashingEngine) Validate(password string, hash string) error {
	parts := strings.SplitN(hash, ":", 2)
	// By default, if the password is not <key>:<hash>, we
	//   just take the default engine. If the prefix was
	//   specified, then a specific engine will be used to
	//   validate that password against the hash.
	engineKey := multipleHashingEngine.defaultEngine
	switch len(parts) {
	case 2:
		engineKey = parts[0]
		hash = parts[1]
	default:
		return InvalidHash
	}
	// Given the engine key, password, and hash, calculate the validation.
	if engine, ok := multipleHashingEngine.registeredEngines[engineKey]; !ok {
		return UnregisteredEngine
	} else {
		return engine.Validate(password, hash)
	}
}

// These implementations have no name.
func (multipleHashingEngine *MultipleHashingEngine) Name() string {
	return ""
}

// Panicked when no hashers were specified when instantiating a multiple hashing engine.
var ErrNoHashers = errors.New("no hashers were specified")

// Panicked when a nil hasher was specified when instantiating a multiple hashing engine.
var ErrNilHasher = errors.New("a nil hasher was specified")

// Panicked when a hasher in the list has empty name when instantiating a multiple hashing engine.
var ErrEmptyHasherName = errors.New("empty hasher name (this is an engine implementation issue)")

// Panicked when a hasher in the list has a duplicate name when instantiating a multiple hashing engine.
var ErrDuplicateHasherName = errors.New("duplicate hasher name (most likely, two hashing engines of the same " +
	"type were registered)")

// Panicked when an explicitly stated default hasher is not among the list when instantiating a multiple hashing engine.
var ErrMissingDefault = errors.New("explicitly specified default hasher is missing among list")

// Panicked when another multiple hashing engine is specified among the list when insantiating a multuple hashing engine.
var ErrNestedMultiHasher = errors.New("nesting multiple hashers is forbidden")

// Creates a new multiple hasher using an explicitly given engine as the default one.
// It panics of no engines, or duplicate-name engines, are given. It also panics if
// the given default engine is not present among the engines.
func NewMultipleHashingEngineWithDefault(defaultEngine HashingEngine, engines ...HashingEngine) HashingEngine {
	if len(engines) == 0 {
		panic(ErrNoHashers)
	}

	mphe := &MultipleHashingEngine{}

	for _, engine := range engines {
		if engine == nil {
			panic(ErrNilHasher)
		} else if _, ok := engine.(*MultipleHashingEngine); ok {
			panic(ErrNestedMultiHasher)
		}
		name := engine.Name()
		if name == "" {
			panic(ErrEmptyHasherName)
		}
		if _, ok := mphe.registeredEngines[name]; ok {
			panic(ErrDuplicateHasherName)
		} else {
			mphe.registeredEngines[name] = engine
		}
	}

	if defaultEngine == nil {
		mphe.defaultEngine = engines[0].Name()
		return mphe
	} else if engine, _ := mphe.registeredEngines[defaultEngine.Name()]; engine != defaultEngine {
		panic(ErrMissingDefault)
	} else {
		mphe.defaultEngine = defaultEngine.Name()
		return mphe
	}
}

// Creates a new multiple hasher using the first engine as the default one.
func NewMultipleHashingEngine(engines ...HashingEngine) HashingEngine {
	return NewMultipleHashingEngineWithDefault(nil, engines...)
}
