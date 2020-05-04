package tests

import (
	"github.com/universe-10th/identity/realm"
	"testing"
	"time"
)

// These tests will make use of user 1, with password: user1$123.
// Login must always succeed in every case to retrieve the credential
// for the first time. Then, it will fail after password changes.

func TestSetPassword(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	_ = userRealm.SetPassword(credential, "user1$456")

	if _, err := userRealm.Login("U1", "user1$123"); err != realm.ErrLoginFailed {
		t.Errorf("After password change, the old password attempt must return realm.ErrLoginFailed. Error returned instead: %s\n", err)
	}

	if _, err := userRealm.Login("U1", "user1$456"); err != nil {
		t.Errorf("After password change, the new password attempt must return no error. Error returned instead: %s\n", err)
	}
}

func TestUnsetPassword(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	_ = userRealm.UnsetPassword(credential)

	if _, err := userRealm.Login("U1", "user1$123"); err != realm.ErrLoginFailed {
		t.Errorf("After password change, the old password attempt must return realm.ErrLoginFailed. Error returned instead: %s\n", err)
	}
}

func TestPasswordResetWithBadToken(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	if err := userRealm.ConfirmPasswordReset(credential, "ab214109sdfb", "new-password"); err != realm.ErrBadToken {
		t.Errorf("Trying to reset the password using a bad token (different token, or attempting a token when none is present) must return realm.ErrBadToken. Error returned instead: %s\n", err)
	}
}

func TestPasswordResetWithExpiredToken(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	_ = userRealm.PreparePasswordReset(credential, "abc123", time.Second)
	time.Sleep(2 * time.Second)
	if err := userRealm.ConfirmPasswordReset(credential, "abc123", "new-password"); err != realm.ErrBadToken {
		t.Errorf("Trying to reset the password using an expired token must return realm.ErrBadToken. Error returned instead: %s\n", err)
	}
}

func TestPasswordResetSuccess(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	_ = userRealm.PreparePasswordReset(credential, "abc123", time.Hour)
	if err := userRealm.ConfirmPasswordReset(credential, "abc123", "user1$456"); err != nil {
		t.Errorf("Trying to reset the password using a good token must return no error. Error returned instead: %s\n", err)
	}

	if _, err := userRealm.Login("U1", "user1$123"); err != realm.ErrLoginFailed {
		t.Errorf("After password reset, the old password attempt must return realm.ErrLoginFailed. Error returned instead: %s\n", err)
	}

	if _, err := userRealm.Login("U1", "user1$456"); err != nil {
		t.Errorf("After password reset, the new password attempt must return no error. Error returned instead: %s\n", err)
	}

}

func TestPasswordResetCancelWithNoToken(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")
	if err := userRealm.CancelPasswordReset(credential); err != nil {
		t.Errorf("Cancelling a password reset should not return an error on itself. Error returned instead: %s\n", err)
	}

	if _, err := userRealm.Login("U1", "user1$123"); err != nil {
		t.Errorf("After password reset cancel, the old password attempt must still work. Error returned instead: %s\n", err)
	}
}

func TestPasswordResetCancelWithToken(t *testing.T) {
	_, realms := MakeUserExampleInstances()
	userRealm := realms[1]

	credential, _ := userRealm.Login("U1", "user1$123")

	_ = userRealm.PreparePasswordReset(credential, "abc123", time.Hour)
	if err := userRealm.CancelPasswordReset(credential); err != nil {
		t.Errorf("Cancelling a password reset should not return an error on itself. Error returned instead: %s\n", err)
	}

	if _, err := userRealm.Login("U1", "user1$123"); err != nil {
		t.Errorf("After password reset cancel, the old password attempt must still work. Error returned instead: %s\n", err)
	}
}
