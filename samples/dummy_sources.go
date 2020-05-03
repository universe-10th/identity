package samples

import (
	"errors"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/scoped"
	"github.com/universe-10th/identity/hashing"
	"reflect"
	"time"
)

type BaseUser struct {
	active         bool
	hashedPassword string
}

func (user *BaseUser) Active() bool {
	return user.active
}

func (user *BaseUser) SetActive(active bool) {
	user.active = active
}

func (user *BaseUser) HashedPassword() string {
	return user.hashedPassword
}

func (user *BaseUser) SetHashedPassword(password string) {
	user.hashedPassword = password
}

func (user *BaseUser) Hasher() hashing.HashingEngine {
	return DummyHasher(0)
}

type User struct {
	BaseUser
	punishedOn  *time.Time
	punishedFor *time.Duration
	punishment  interface{}
	punisher    credentials.Credential
}

func (user *User) PunishedFor() (punishedOn *time.Time, forTime *time.Duration, reason interface{}, by credentials.Credential) {
	return user.punishedOn, user.punishedFor, user.punishment, user.punisher
}

func (user *User) Punish(forTime *time.Duration, reason interface{}, by credentials.Credential) {
	now := new(time.Time)
	*now = time.Now()
	punishedFor := (*time.Duration)(nil)
	if forTime != nil {
		punishedFor = new(time.Duration)
		*punishedFor = *forTime
	}
	user.punishedOn = now
	user.punishedFor = punishedFor
	user.punishment = reason
	user.punisher = by
}

func (user *User) Unpunish() {
	user.punishedOn = nil
	user.punishedFor = nil
	user.punishment = nil
	user.punisher = nil
}

type Admin struct {
	BaseUser
	superuser bool
	scopes    map[string]scoped.Scope
}

func (Admin) Staff() bool {
	return true
}

func (admin *Admin) Superuser() bool {
	return admin.superuser
}

func (admin *Admin) Scopes() map[string]scoped.Scope {
	return admin.scopes
}

type DummyBroker struct {
	dataByIdentifier map[reflect.Type]map[string]credentials.Credential
	dataByIndex      map[reflect.Type]map[int]credentials.Credential
}

func (broker *DummyBroker) Allows(template credentials.Credential) bool {
	_, ok := broker.dataByIndex[reflect.TypeOf(template)]
	return ok
}

var ErrNotFound = errors.New("identifier not found")

func (broker *DummyBroker) ByIdentifier(identifier interface{}, template credentials.Credential) (credentials.Credential, error) {
	if result, ok := broker.dataByIdentifier[reflect.TypeOf(template)][identifier.(string)]; !ok {
		return nil, ErrNotFound
	} else {
		return result, nil
	}
}

func (broker *DummyBroker) ByIndex(index interface{}, template credentials.Credential) (credentials.Credential, error) {
	if result, ok := broker.dataByIndex[reflect.TypeOf(template)][index.(int)]; !ok {
		return nil, ErrNotFound
	} else {
		return result, nil
	}
}

func (broker *DummyBroker) Save(credential credentials.Credential) error {
	return nil
}
