package tests

import (
	"github.com/universe-10th/identity/authreqs"
	"github.com/universe-10th/identity/authreqs/compound"
	scoped2 "github.com/universe-10th/identity/authreqs/scoped"
	"github.com/universe-10th/identity/authreqs/superuser"
	"github.com/universe-10th/identity/credentials"
	"github.com/universe-10th/identity/credentials/traits/scoped"
	"github.com/universe-10th/identity/hashing"
	"github.com/universe-10th/identity/realms"
	"github.com/universe-10th/identity/realms/login/activity"
	"github.com/universe-10th/identity/realms/login/password"
	"github.com/universe-10th/identity/realms/login/punish"
	"reflect"
	"time"
)

func MakeUserExampleInstances() ([]authreqs.AuthorizationRequirement, []*realms.Realm) {
	hasher := (&BaseUser{}).Hasher()
	hash := func(input string) string {
		hashed, _ := hasher.Hash(input)
		return hashed
	}
	ago := func(duration time.Duration) *time.Time {
		result := new(time.Time)
		*result = time.Now().Add(-duration)
		return result
	}
	ptr := func(duration time.Duration) *time.Duration {
		return &duration
	}

	scope2 := DummyScope(2)
	scope3 := DummyScope(3)
	scope5 := DummyScope(5)
	scope7 := DummyScope(7)

	adminSU := &Admin{BaseUser: BaseUser{active: true, hashedPassword: hash("admin-su$123")}, superuser: true, scopes: nil}
	adminS1 := &Admin{BaseUser: BaseUser{active: true, hashedPassword: hash("admin-s1$123")}, superuser: false, scopes: map[string]scoped.Scope{
		scope2.Key(): scope2,
		scope3.Key(): scope3,
	}}
	adminS2 := &Admin{BaseUser: BaseUser{active: true, hashedPassword: hash("admin-s2$123")}, superuser: false, scopes: map[string]scoped.Scope{
		scope5.Key(): scope5,
		scope7.Key(): scope7,
	}}
	adminS3 := &Admin{BaseUser: BaseUser{active: true, hashedPassword: hash("admin-s3$123")}, superuser: false, scopes: map[string]scoped.Scope{
		scope5.Key(): scope5,
		scope3.Key(): scope3,
	}}
	user1 := &User{BaseUser: BaseUser{active: true, hashedPassword: hash("user1$123")}}
	user2 := &User{BaseUser: BaseUser{active: false, hashedPassword: hash("user2$123")}}
	user3 := &User{
		BaseUser: BaseUser{active: true, hashedPassword: hash("user3$123")},
		// Punishment: expired
		punishedOn:  ago(time.Hour * 24 * 7),
		punishedFor: ptr(time.Hour * 24 * 3),
		punishment:  "Sample Punishment  (expired)",
		punisher:    adminS1,
	}
	user4 := &User{
		BaseUser: BaseUser{active: true, hashedPassword: hash("user4$123")},
		// Punishment: current
		punishedOn:  ago(time.Hour * 24 * 7),
		punishedFor: ptr(time.Hour * 24 * 8),
		punishment:  "Sample Punishment (active)",
		punisher:    adminS1,
	}
	user5 := &User{
		BaseUser: BaseUser{active: true, hashedPassword: hash("user5$123")},
		// Punishment: eternal
		punishedOn:  ago(time.Hour * 24 * 7),
		punishedFor: nil,
		punishment:  "Sample Punishment (eternal)",
		punisher:    adminS1,
	}
	broker := &DummyBroker{
		dataByIndex: map[reflect.Type]map[int]credentials.Credential{
			reflect.TypeOf(&Admin{}): {
				0: adminSU,
				1: adminS1,
				2: adminS2,
				3: adminS3,
			},
			reflect.TypeOf(&User{}): {
				1: user1,
				2: user2,
				3: user3,
				4: user4,
				5: user5,
			},
		},
		dataByIdentifier: map[reflect.Type]map[string]credentials.Credential{
			reflect.TypeOf(&Admin{}): {
				"SU": adminSU,
				"S1": adminS1,
				"S2": adminS2,
				"S3": adminS3,
			},
			reflect.TypeOf(&User{}): {
				"U1": user1,
				"U2": user2,
				"U3": user3,
				"U4": user4,
				"U5": user5,
			},
		},
	}

	users := credentials.NewSource(broker, &User{})
	admins := credentials.NewSource(broker, &Admin{})

	adminReq1 := compound.Admin(scope2, scope3)
	adminReq2 := compound.Admin(scope5, scope7)
	tryAll := compound.TryAll{superuser.RequireSuperuser, scoped2.RequireScopesAmong(scope2, scope7)}

	adminRealm := realms.NewRealm(admins, activity.ActivityStep(0), password.PasswordCheckingStep(0))
	userRealm := realms.NewRealm(users, activity.ActivityStep(0), password.PasswordCheckingStep(0), &punish.PunishmentCheckStep{TimeFormat: "2006-01-02T15:04:05"})

	// We have the realms, the users, and the admin requirements.
	requirements := []authreqs.AuthorizationRequirement{
		// Remember to use these in order.
		adminReq1, adminReq2, tryAll,
	}
	realms := []*realms.Realm{
		// Remember to use these in order.
		adminRealm, userRealm,
	}
	return requirements, realms
}

func MakeMultiHasherExampleInstances() (hashing.HashingEngine, DummyHasher, DummyHasher) {
	h0 := DummyHasher(0)
	h1 := DummyHasher(1)
	multi := hashing.NewMultipleHashingEngine(h0, h1)
	return multi, h0, h1
}
