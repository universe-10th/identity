package tests

import "testing"

func TestAdminIsSuperuser(t *testing.T) {
	requirements, realms := MakeUserExampleInstances()
	adminRealm := realms[0]
	adminReq23 := requirements[0]

	if credential, err := adminRealm.Login("SU", "admin-su$123"); err != nil {
		t.Errorf("This login must not fail! Error received: %s\n", err)
	} else if !adminReq23.SatisfiedBy(credential) {
		t.Error("SU account must satisfy any Admin requirement")
	}
}

func TestAdminSatisfiesScopes(t *testing.T) {
	requirements, realms := MakeUserExampleInstances()
	adminRealm := realms[0]
	adminReq23 := requirements[0]
	adminReq57 := requirements[1]

	if credential, err := adminRealm.Login("S1", "admin-s1$123"); err != nil {
		t.Errorf("This login must not fail! Error received: %s\n", err)
	} else if !adminReq23.SatisfiedBy(credential) {
		t.Error("S1 has scopes 2 and 3. It must satisfy the first admin requirement")
	} else if adminReq57.SatisfiedBy(credential) {
		t.Error("S1 has scopes 2 and 3. It must NOT satisfy the second admin requirement")
	}

	if credential, err := adminRealm.Login("S3", "admin-s3$123"); err != nil {
		t.Errorf("This login must not fail! Error received: %s\n", err)
	} else if !adminReq23.SatisfiedBy(credential) {
		t.Error("S1 has scopes 5 and 3. It must satisfy the first admin requirement")
	} else if !adminReq57.SatisfiedBy(credential) {
		t.Error("S1 has scopes 5 and 3. It must satisfy the second admin requirement")
	}
}

func TestAdminNothing(t *testing.T) {
	requirements, realms := MakeUserExampleInstances()
	userRealm := realms[1]
	adminReq23 := requirements[0]
	adminReq57 := requirements[1]

	if credential, err := userRealm.Login("U1", "user1$123"); err != nil {
		t.Errorf("This login must not fail! Error received: %s\n", err)
	} else if adminReq23.SatisfiedBy(credential) || adminReq57.SatisfiedBy(credential) {
		t.Error("U1 is a non-staff, non-superuser user. It must not satisfy any admin requirement")
	}
}

func TestTryAllIsSuperuser(t *testing.T) {

}

func TestTryAllSatisfiesScopes(t *testing.T) {

}

func TestTryAllDoesNotSatisfyScopes(t *testing.T) {

}

func TestTryAllNothing(t *testing.T) {

}
