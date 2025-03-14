package user

import "errors"

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	SsoId string `json:"ssoId"`
}

var ErrUserAlreadyExists = errors.New("User already exists")
