package models

import "errors"

var (
	UserExistsError        = errors.New("user exists")
	NewUserDataError       = errors.New("user exists")
	SlugAlreadyExistsError = errors.New("slug already exists")
)
