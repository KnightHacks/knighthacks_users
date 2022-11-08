package repository

import "errors"

var (
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user with id already exists")
)
