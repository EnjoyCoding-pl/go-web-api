package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Hash  string
	Login string
}

func NewUser(login string, password string, repeatedPassword string) (*User, error) {

	if password != repeatedPassword {
		return nil, errors.New("password not matched")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{Login: login, Hash: string(hash)}, nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
