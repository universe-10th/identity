package tests

import (
	"github.com/universe-10th/identity/realm"
	"github.com/universe-10th/identity/realm/login/punish"
	"testing"
)

func TestAdminLoginBadPassword(t *testing.T) {
	_, realms := MakeExampleInstances()
	adminRealm := realms[0]

	// Trying to log the SU user with password admin-su$124, which is wrong.
	if _, err := adminRealm.Login("SU", "admin-su$124"); err != realm.ErrLoginFailed {
		t.Error("Login for user SU must fail with an invalid password")
	}
	// Trying to log the SU2 user, which does not exist.
	if _, err := adminRealm.Login("SU2", "admin-su$123"); err != realm.ErrLoginFailed {
		t.Error("Login for user SU must fail with an invalid user")
	}
}

func TestAdminLoginOK(t *testing.T) {
	_, realms := MakeExampleInstances()
	adminRealm := realms[0]

	// Trying to log the SU user with the right password.
	if _, err := adminRealm.Login("SU", "admin-su$124"); err != nil {
		t.Errorf("Login for user SU must succeed with the appropriate password. Error: %s\n", err)
	}
}

func TestUserLoginBadPassword(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U1 user with the password: user1$124, which is wrong.
	if _, err := userRealm.Login("U1", "user1$124"); err != realm.ErrLoginFailed {
		t.Error("Login for user U1 must fail with an invalid password")
	}
}

func TestUserLoginInactive(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U2 user with the password: user1$123, but being inactive.
	if _, err := userRealm.Login("U2", "user2$123"); err != realm.ErrLoginFailed {
		t.Error("Login for user U2 must fail since it is inactive")
	}
}

func TestUserLoginPunishmentStill(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U4 user with the password: user4$123, but being punished.
	if _, err := userRealm.Login("U4", "user4$123"); err == nil {
		t.Error("Login for user U4 must fail since it is with an active punishment")
	} else if pe, _ := err.(*punish.PunishedError); pe == nil {
		t.Error("Login for user U4 must fail with PunishedError since it is with an active punishment")
	}
}

func TestUserLoginPunishmentPermanent(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U5 user with the password: user5$123, but being punished.
	if _, err := userRealm.Login("U5", "user5$123"); err == nil {
		t.Error("Login for user U5 must fail since it is with an eternal punishment")
	} else if pe, _ := err.(*punish.PunishedError); pe == nil {
		t.Error("Login for user U5 must fail with PunishedError since it is with an eternal punishment")
	} else if pe.PunishedFor != nil {
		t.Error("Login for user U5 must fail with a nil-time PunishedError since it is with an eternal punishment")
	}
}

func TestUserLoginOKUnpunished(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U1 user with the password: user1$123, must succeed.
	if _, err := userRealm.Login("U1", "user1$123"); err != nil {
		t.Error("Login for user U1 must succeed")
	}
}

func TestUserLoginOKPunishmentExpired(t *testing.T) {
	_, realms := MakeExampleInstances()
	userRealm := realms[1]

	// Trying to log the U3 user with the password: user1$123 and punishment expired, must succeed.
	if _, err := userRealm.Login("U3", "user3$123"); err != nil {
		t.Error("Login for user U3 must succeed")
	}
}
