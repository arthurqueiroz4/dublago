package user

import "errors"

type User struct {
	ID    uint
	Email string
	Name  string
	SsoId string
}

var ErrUserAlreadyExists = errors.New("User already exists")
