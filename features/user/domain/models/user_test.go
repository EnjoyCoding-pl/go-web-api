package models

import "testing"

func TestNewUser_Success(t *testing.T) {
	name := "testUser"
	password := "test_password"
	repeatedPassword := "test_password"

	_, err := NewUser(name, password, repeatedPassword)
	if err != nil {
		t.Fatalf("Create new user failed, msg: %v", err)
	}
}

func TestNewUser_FailedDifferentPassword(t *testing.T) {
	name := "testUser"
	password := "test_password"
	repeatedPassword := "test_password2"

	_, err := NewUser(name, password, repeatedPassword)
	if err == nil {
		t.Fatalf("User created with different passwords, msg: %v", err)
	}
}

func TestCheckPassword_Success(t *testing.T) {
	name := "testUser"
	password := "test_password"
	testPassword := "test_password"

	user, err := NewUser(name, password, password)
	if err != nil {
		t.Fatalf("Create new user failed, msg: %v", err)
	}

	err = user.CheckPassword(testPassword)

	if err != nil {
		t.Fatalf("Check user password failed, msg: %v", err)
	}
}

func TestCheckPassword_FailedDifferentPasswords(t *testing.T) {
	name := "testUser"
	password := "test_password"
	testPassword := "test_password1"

	user, err := NewUser(name, password, password)
	if err != nil {
		t.Fatalf("Create new user failed, msg: %v", err)
	}

	err = user.CheckPassword(testPassword)

	if err == nil {
		t.Fatalf("Different passwords pass verification, msg: %v", err)
	}
}
